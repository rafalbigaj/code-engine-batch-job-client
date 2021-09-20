/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.ibm.com/coligo/batch-job-controller/pkg/ctxlog"
	"github.ibm.com/coligo/batch-job-controller/pkg/utils/errors"
)

type mergableJobDefContainer struct {
	corev1.Container
}

const (
	arraySpecDefault        = "0"
	retryLimitDefault       = 3
	maxExecutionTimeDefault = 7200
	// As of Kubernetes, the max limit on a name is 63 chars.
	// We reserve 10 for controller to add data{arraySpec, retryLimit}.
	// We reserve 5 for random prefix string. so final length is 63-10-5 = 48
	maxGenerateNameLength = 48
)

var (
	// see also: https://github.ibm.com/coligo/project-isolation-controller/blob/main/deployment/templates/limitrange/default.yaml
	staticResourceSizeCPU              = resource.MustParse("1")
	staticResourceSizeMemory           = resource.MustParse("4G")
	staticResourceSizeEphemeralStorage = resource.MustParse("4G")
)

func (jr *JobRun) MutateWithDefaults() {
	jr.Spec.JobDefinitionSpec.setArraySpecIfNil(arraySpecDefault)
	jr.Spec.JobDefinitionSpec.setRetryLimitIfNil(retryLimitDefault)
	jr.Spec.JobDefinitionSpec.setMaxExecutionTimeIfNil(maxExecutionTimeDefault)
	jr.trimGenerateNameIfTooLong()
}

func (jr *JobRun) MutateWithJobDefinition(jobDef *JobDefinition) (errs field.ErrorList) {
	errs = field.ErrorList{}
	jr.trimGenerateNameIfTooLong()
	if jr.Spec.JobDefinitionSpec.IsEmpty() {
		jobDef.Spec.DeepCopyInto(&jr.Spec.JobDefinitionSpec)
	} else {
		if errs = jr.Spec.JobDefinitionSpec.Template.checkDisallowedFieldsForMutation(jobDef.Spec.Template, field.NewPath("spec.jobDefinitionSpec.template")); len(errs) != 0 {
			return
		}

		jr.Spec.JobDefinitionSpec.setArraySpecIfNil(*jobDef.Spec.ArraySpec)
		jr.Spec.JobDefinitionSpec.setRetryLimitIfNil(*jobDef.Spec.RetryLimit)
		jr.Spec.JobDefinitionSpec.setMaxExecutionTimeIfNil(*jobDef.Spec.MaxExecutionTime)
		jr.Spec.JobDefinitionSpec.mergeTemplateWith(jobDef.Spec.Template)
	}

	jr.AddLabel(LabelJobDefName, jobDef.Name, false)
	jr.AddLabel(LabelJobDefUUID, string(jobDef.UID), false)
	jr.SetOwner(jobDef)

	return
}

func (jr *JobRun) MutateResourcesRequestsWithLimitRange(containerLimitRangeItem *corev1.LimitRangeItem) (errs field.ErrorList) {
	errs = field.ErrorList{}

	if len(jr.Spec.JobDefinitionSpec.Template.Containers) != 1 {
		msg := "there must be exactly one container"
		errs = append(errs, errors.ContainerCountUnsupported(field.NewPath("spec.jobDefinitionSpec.template.containers"), len(jr.Spec.JobDefinitionSpec.Template.Containers), msg))
		return
	}

	if containerLimitRangeItem == nil {
		return
	}

	if jr.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests == nil {
		jr.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests = corev1.ResourceList{}
	}

	jr.mutateRequestResource(corev1.ResourceCPU, containerLimitRangeItem)
	jr.mutateRequestResource(corev1.ResourceMemory, containerLimitRangeItem)
	jr.mutateRequestResource(corev1.ResourceEphemeralStorage, containerLimitRangeItem)

	return
}

func (jr *JobRun) mutateRequestResource(resourceName corev1.ResourceName, containerLimitRangeItem *corev1.LimitRangeItem) {
	if _, exists := jr.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests[resourceName]; !exists {
		defaultValue := containerLimitRangeItem.DefaultRequest[resourceName]
		jr.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests[resourceName] = defaultValue
	}
}

func (jr *JobRun) trimGenerateNameIfTooLong() {
	if len(jr.ObjectMeta.GenerateName) > maxGenerateNameLength {
		jr.ObjectMeta.GenerateName = jr.ObjectMeta.GenerateName[:maxGenerateNameLength]
	}
}

func (jr *JobRun) MutateResourcesRequestsWithStaticLimits(ctx context.Context) (errs field.ErrorList) {
	errs = field.ErrorList{}

	if len(jr.Spec.JobDefinitionSpec.Template.Containers) != 1 {
		msg := "there must be exactly one container"
		errs = append(errs, errors.ContainerCountUnsupported(field.NewPath("spec.jobDefinitionSpec.template.containers"), len(jr.Spec.JobDefinitionSpec.Template.Containers), msg))
		return
	}

	if jr.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests == nil {
		jr.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests = corev1.ResourceList{}
	}

	jr.Spec.JobDefinitionSpec.Template.setResourceRequestIfNil(ctx, corev1.ResourceCPU, staticResourceSizeCPU)
	jr.Spec.JobDefinitionSpec.Template.setResourceRequestIfNil(ctx, corev1.ResourceMemory, staticResourceSizeMemory)
	jr.Spec.JobDefinitionSpec.Template.setResourceRequestIfNil(ctx, corev1.ResourceEphemeralStorage, staticResourceSizeEphemeralStorage)

	return nil
}

func (pT *JobPodTemplate) setResourceRequestIfNil(ctx context.Context, resource corev1.ResourceName, quantity resource.Quantity) {
	if _, exists := pT.Containers[0].Resources.Requests[resource]; !exists {
		ctxlog.Warnf(ctx, "Falling back to static value %s for %s", quantity.String(), resource.String())
		pT.Containers[0].Resources.Requests[resource] = quantity
	}
}

func (jds *JobDefinitionSpec) setArraySpecIfNil(arrayspec string) {
	if jds.ArraySpec == nil {
		jds.ArraySpec = &arrayspec
	}
}

func (jds *JobDefinitionSpec) setRetryLimitIfNil(retryLimit int64) {
	if jds.RetryLimit == nil {
		jds.RetryLimit = &retryLimit
	}
}

func (jds *JobDefinitionSpec) setMaxExecutionTimeIfNil(maxExecTime int64) {
	if jds.MaxExecutionTime == nil {
		jds.MaxExecutionTime = &maxExecTime
	}
}

func (jds *JobDefinitionSpec) mergeTemplateWith(podTemplate JobPodTemplate) {
	jds.Template.mergeWith(podTemplate)
}

func (pT *JobPodTemplate) mergeWith(jobDefPodTemplate JobPodTemplate) {
	pT.ImagePullSecrets = jobDefPodTemplate.ImagePullSecrets
	pT.ServiceAccountName = jobDefPodTemplate.ServiceAccountName

	pT.mergeContainersWith(jobDefPodTemplate.Containers)
}

func (pT *JobPodTemplate) mergeContainersWith(jdContainers []corev1.Container) {
	if !pT.hasContainers() {
		pT.Containers = jdContainers
	} else {
		// Currently we only allow one container and therefore we always pick the first containers.
		// Grep for d31add5e to see other references to this constraint
		jdContainer := mergableJobDefContainer{jdContainers[0]}
		jdContainer.mergeWith(pT.Containers[0])

		pT.Containers = []corev1.Container{jdContainer.Container}
	}

}

func (pT *JobPodTemplate) hasContainers() bool {
	return len(pT.Containers) > 0
}

func (mC *mergableJobDefContainer) mergeWith(jobRunContainer corev1.Container) {
	mC.overwriteCommand(jobRunContainer.Command)
	mC.overwriteArgs(jobRunContainer.Args)
	mC.overwriteEnvFrom(jobRunContainer.EnvFrom)
	mC.mergeResourceRequests(jobRunContainer.Resources.Requests)

	mC.mergeEnvVars(jobRunContainer.Env)

}

func (mC *mergableJobDefContainer) overwriteCommand(command []string) {
	if command != nil {
		mC.Command = command
	}
}

func (mC *mergableJobDefContainer) overwriteArgs(args []string) {
	if args != nil {
		mC.Args = args
	}
}

func (mC *mergableJobDefContainer) overwriteEnvFrom(envFrom []corev1.EnvFromSource) {
	if envFrom != nil {
		mC.EnvFrom = envFrom
	}
}

func (mC *mergableJobDefContainer) mergeEnvVars(jobRunEnvVars []corev1.EnvVar) {
	if mC.Env == nil {
		mC.Env = jobRunEnvVars
	} else {
		envsToAppend := []corev1.EnvVar{}
		for _, jobRunEnvVar := range jobRunEnvVars {
			if ok, idx := mC.searchEnvVarWithName(jobRunEnvVar.Name); ok {
				if isEnvVarValueEmpty(jobRunEnvVar) {
					// delete
					mC.Env = append(mC.Env[:idx], mC.Env[idx+1:]...)
				} else {
					// replace
					mC.Env[idx] = jobRunEnvVar
				}
			} else {
				// append
				if !isEnvVarValueEmpty(jobRunEnvVar) {
					envsToAppend = append(envsToAppend, jobRunEnvVar)
				}
			}
		}
		mC.Env = append(mC.Env, envsToAppend...)
	}
}

func (mC *mergableJobDefContainer) mergeResourceRequests(jobRunContainerRequests corev1.ResourceList) {
	for resourceName, resourceQuantity := range jobRunContainerRequests {
		if mC.Container.Resources.Requests == nil {
			mC.Container.Resources.Requests = corev1.ResourceList{}
		}
		mC.Container.Resources.Requests[resourceName] = resourceQuantity
	}
}

func isEnvVarValueEmpty(envVar corev1.EnvVar) bool {
	return envVar.Value == "" && envVar.ValueFrom == nil
}

func (mC *mergableJobDefContainer) searchEnvVarWithName(name string) (bool, int) {
	for idx, env := range mC.Env {
		if env.Name == name {
			return true, idx
		}
	}
	return false, -1
}

func (pT *JobPodTemplate) checkDisallowedFieldsForMutation(jdPodTemplate JobPodTemplate, fld *field.Path) (errs field.ErrorList) {
	errs = field.ErrorList{}

	if pT.ServiceAccountName != "" {
		errs = append(errs, errors.Unsupported(fld.Child("serviceAccountName"), "disallow to set with referenced jobDefinition"))
	}

	if pT.ImagePullSecrets != nil {
		errs = append(errs, errors.Unsupported(fld.Child("imagePullSecrets"), "disallow to set with referenced jobDefinition"))
	}

	if pT.hasContainers() {
		// Currently we only allow one container and therefore we always pick the first containers.
		// Grep for d31add5e to see other references to this constraint
		jdContainer := jdPodTemplate.Containers[0]
		jrContainer := pT.Containers[0]
		errs = append(errs, checkDisallowedFieldsForMerging(jdContainer, jrContainer, fld.Child("containers").Index(0))...)
	}

	return
}

func checkDisallowedFieldsForMerging(jdContainer, jrContainer corev1.Container, fld *field.Path) (errs field.ErrorList) {
	errs = field.ErrorList{}

	if jrContainer.Name != "" && jdContainer.Name != jrContainer.Name {
		errs = append(errs, errors.Invalid(fld.Child("name"), fmt.Sprintf("%q", jdContainer.Name), "must be the same as referenced jobDefinition.spec.template.containers[0].name"))
	}

	if jrContainer.Image != "" && jdContainer.Image != jrContainer.Image {
		errs = append(errs, errors.Invalid(fld.Child("image"), fmt.Sprintf("%q", jdContainer.Image), "must be the same as referenced jobDefinition.spec.template.containers[0].image"))
	}

	if jrContainer.TerminationMessagePath != "" {
		errs = append(errs, errors.Unsupported(fld.Child("terminationMessagePath"), "disallow to set with referenced jobDefinition"))
	}

	if jrContainer.TerminationMessagePolicy != "" {
		errs = append(errs, errors.Unsupported(fld.Child("terminationMessagePolicy"), "disallow to set with referenced jobDefinition"))
	}

	if jrContainer.WorkingDir != "" {
		errs = append(errs, errors.Unsupported(fld.Child("workingDir"), "disallow to set with referenced jobDefinition"))
	}

	return
}
