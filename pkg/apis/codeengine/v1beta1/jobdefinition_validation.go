/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020, 2021
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1

import (
	"encoding/json"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.ibm.com/coligo/batch-job-controller/pkg/utils/errors"
)

type ValidationConfig struct {
	ContainerLimitRangeItem *corev1.LimitRangeItem
	ResourceRatios          ResourceRatios
}

type ResourceRatios []corev1.ResourceList

func (r ResourceRatios) String() string {
	var s []string
	for _, ratio := range r {
		s = append(s, fmt.Sprintf("%s / %s", ratio.Cpu().String(), ratio.Memory().String()))
	}
	return strings.Join(s, ", ")
}

// TemplateSizeLimit is the maximum allowed size of the "template" field.
const TemplateSizeLimit = 10 * 1024 // 10 KiB

// Validate ensures JobDefinition is properly configured.
func (jd *JobDefinition) Validate(validationConfig ValidationConfig) field.ErrorList {
	return ValidateJobDefinitionSpec(&jd.Spec, validationConfig, field.NewPath("spec"))
}

// ValidateJobDefinitionSpec ensures JobDefinitionSpec is properly configured.
func ValidateJobDefinitionSpec(jdSpec *JobDefinitionSpec, validationConfig ValidationConfig, fldPath *field.Path) field.ErrorList {
	errs := field.ErrorList{}

	if jdSpec.ArraySpec == nil {
		errs = append(errs, errors.Missing(fldPath.Child("arraySpec")))
	} else {
		_, err := jdSpec.CalculateArrayIndices()
		if err != nil {
			errs = append(errs, errors.Invalid(fldPath.Child("arraySpec"), fmt.Sprintf("%q", *jdSpec.ArraySpec), err.Error()))
		}
	}

	if jdSpec.RetryLimit == nil {
		errs = append(errs, errors.Missing(fldPath.Child("retryLimit")))
	}

	if jdSpec.Template.IsEmpty() {
		errs = append(errs, errors.Missing(fldPath.Child("template")))
	}

	return append(errs, ValidateJobPodTemplate(&jdSpec.Template, validationConfig, fldPath.Child("template"))...)
}

// ValidateJobPodTemplate ensures JobPodTemplate is properly configured.
func ValidateJobPodTemplate(template *JobPodTemplate, validationConfig ValidationConfig, fldPath *field.Path) field.ErrorList {
	errs := field.ErrorList{}

	if err := checkTemplateSize(*template, fldPath); err != nil {
		errs = append(errs, err)
	}

	if len(template.Containers) != 1 {
		msg := "there must be exactly one container"
		errs = append(errs, errors.ContainerCountUnsupported(fldPath.Child("containers"), len(template.Containers), msg))
	}

	return append(errs, ValidateContainers(template.Containers, validationConfig, fldPath.Child("containers"))...)
}

// ValidateContainers ensures Containers is properly configured.
func ValidateContainers(containers []corev1.Container, validationConfig ValidationConfig, fldPath *field.Path) field.ErrorList {
	errs := field.ErrorList{}

	for i, container := range containers {
		errs = checkDisallowedFields(container, fldPath.Index(i))

		errs = append(errs, validateEnvVarSource(container.Env, fldPath.Index(i).Child("env"))...)

		errs = append(errs, validateResources(container.Resources, validationConfig, fldPath.Index(i).Child("resources"))...)

		if container.Name == "" {
			errs = append(errs, errors.Missing(fldPath.Index(i).Child("name")))
		} else if e := validation.IsDNS1123Label(container.Name); len(e) != 0 {
			// IsDNS1123Label tests for the container name to follow the DNS label standard
			msg := fmt.Sprintf("invalid value: %s", strings.Join(e, ","))
			errs = append(errs, errors.Invalid(fldPath.Index(i).Child("name"), fmt.Sprintf("%q", container.Name), msg))
		}

		if container.Image == "" {
			errs = append(errs, errors.Missing(fldPath.Index(i).Child("image")))
		}
	}
	return errs
}

func validateEnvVarSource(envVars []corev1.EnvVar, fldPath *field.Path) field.ErrorList {
	errs := field.ErrorList{}

	for i, env := range envVars {
		if env.Name == JobIndex {
			msg := fmt.Sprintf("'%s' is reserved for batch API", env.Name)
			errs = append(errs, errors.Unsupported(fldPath.Index(i).Child("name"), msg))
		}

		valueFrom := env.ValueFrom
		if valueFrom != nil {
			fieldRef := valueFrom.FieldRef
			if fieldRef != nil {
				errs = append(errs, errors.Unsupported(fldPath.Index(i).Child("valueFrom.fieldRef"), ""))
			}
			resourceFieldRef := valueFrom.ResourceFieldRef
			if resourceFieldRef != nil {
				errs = append(errs, errors.Unsupported(fldPath.Index(i).Child("valueFrom.resourceFieldRef"), ""))
			}
		}
	}

	return errs
}

func validateResources(resources corev1.ResourceRequirements, validationConfig ValidationConfig, fldPath *field.Path) field.ErrorList {
	errs := field.ErrorList{}

	res := []corev1.ResourceName{corev1.ResourceCPU, corev1.ResourceMemory, corev1.ResourceEphemeralStorage}

	for _, r := range res {
		err := validateResourceQuantity(resources, r, validationConfig.ContainerLimitRangeItem, fldPath.Child("requests", r.String()))
		if err != nil {
			errs = append(errs, err)
		}
	}

	err := validateResourceRatios(resources, validationConfig.ResourceRatios, fldPath)
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}

func validateResourceRatios(resources corev1.ResourceRequirements, ratios ResourceRatios, fldPath *field.Path) *field.Error {
	if len(ratios) > 0 && len(resources.Requests) > 0 {
		for _, ratio := range ratios {
			if ratio.Cpu().Cmp(*resources.Requests.Cpu()) == 0 && ratio.Memory().Cmp(*resources.Requests.Memory()) == 0 {
				return nil
			}
		}
		return errors.UnsupportedResourceSpecification(fldPath.Child("requests"), resources.Requests, ratios.String())
	}
	return nil
}

func validateResourceQuantity(
	resources corev1.ResourceRequirements,
	resourceName corev1.ResourceName,
	containerLimitRangeItem *corev1.LimitRangeItem,
	fldPath *field.Path) (err *field.Error) {

	resourceQuantity, requestExists := resources.Requests[resourceName]
	if requestExists {
		if resourceQuantity.Cmp(resource.Quantity{}) <= 0 {
			msg := "must be greater than 0"
			err = errors.BeyondRange(fldPath, resourceQuantity.String(), msg)
		} else if containerLimitRangeItem != nil {
			err = validateResourceQuantityInLimitRange(resourceQuantity, *containerLimitRangeItem, resourceName, fldPath)
		}
	}

	return
}

func validateResourceQuantityInLimitRange(
	resourceQuantity resource.Quantity,
	limitRangeItem corev1.LimitRangeItem, resourceName corev1.ResourceName,
	fldPath *field.Path,
) (err *field.Error) {
	minQuantity, hasMin := limitRangeItem.Min[resourceName]
	maxQuantity, hasMax := limitRangeItem.Max[resourceName]

	// Return error if either resource < min or resource > max.
	if (hasMin && resourceQuantity.Cmp(minQuantity) < 0) || (hasMax && resourceQuantity.Cmp(maxQuantity) > 0) {
		minQuantityString := func() string {
			if hasMin {
				return minQuantity.String()
			}
			return "n/a"
		}()
		maxQuantityString := func() string {
			if hasMax {
				return maxQuantity.String()
			}
			return "n/a"
		}()
		msg := fmt.Sprintf("expected %v <= %v <= %v", minQuantityString, resourceQuantity.String(), maxQuantityString)
		err = errors.BeyondRange(fldPath, resourceQuantity.String(), msg)
	}

	return
}

func checkDisallowedFields(container corev1.Container, fldPath *field.Path) field.ErrorList {
	errs := field.ErrorList{}

	// Disallowed fields
	if container.Lifecycle != nil {
		errs = append(errs, errors.Unsupported(fldPath.Child("lifecycle"), ""))
	}

	if container.LivenessProbe != nil {
		errs = append(errs, errors.Unsupported(fldPath.Child("livenessProbe"), ""))
	}

	if container.ReadinessProbe != nil {
		errs = append(errs, errors.Unsupported(fldPath.Child("readinessProbe"), ""))
	}

	if container.StartupProbe != nil {
		errs = append(errs, errors.Unsupported(fldPath.Child("startupProbe"), ""))
	}

	if container.Ports != nil {
		errs = append(errs, errors.Unsupported(fldPath.Child("ports"), ""))
	}

	if container.SecurityContext != nil {
		errs = append(errs, errors.Unsupported(fldPath.Child("securityContext"), ""))
	}

	if container.Stdin {
		errs = append(errs, errors.Unsupported(fldPath.Child("stdin"), ""))
	}

	if container.StdinOnce {
		errs = append(errs, errors.Unsupported(fldPath.Child("stdinOnce"), ""))
	}

	if container.TTY {
		errs = append(errs, errors.Unsupported(fldPath.Child("tty"), ""))
	}

	if container.VolumeDevices != nil {
		errs = append(errs, errors.Unsupported(fldPath.Child("volumeDevices"), ""))
	}

	if container.VolumeMounts != nil {
		errs = append(errs, errors.Unsupported(fldPath.Child("volumeMounts"), ""))
	}

	return errs
}

func checkTemplateSize(template JobPodTemplate, fldPath *field.Path) *field.Error {
	templateBytes, _ := json.Marshal(template)
	templateSize := len(templateBytes)

	if templateSize > TemplateSizeLimit {
		msg := fmt.Sprintf("size limit is %d bytes, but actual size is %d bytes", TemplateSizeLimit, templateSize)
		return errors.SizeLimitExceeded(fldPath, msg)
	}

	return nil
}
