/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020, 2021
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1_test

import (
	"encoding/json"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/validation/field"

	cv1b1 "github.ibm.com/coligo/batch-job-controller/pkg/apis/codeengine/v1beta1"
	"github.ibm.com/coligo/batch-job-controller/pkg/utils/errors"
	"github.ibm.com/coligo/batch-job-controller/pkg/utils/pointers"
)

var _ = Describe("JobDefinitionValidation", func() {
	var (
		validationConfig cv1b1.ValidationConfig
	)

	BeforeEach(func() {
		validationConfig = cv1b1.ValidationConfig{}
	})

	When("validating a JobDefinition", func() {
		var (
			jobDef *cv1b1.JobDefinition
			err    field.ErrorList
		)
		JustBeforeEach(func() {
			err = jobDef.Validate(validationConfig)
		})

		Context("that is valid", func() {
			BeforeEach(func() {
				jobDef = &cv1b1.JobDefinition{
					Spec: cv1b1.JobDefinitionSpec{
						ArraySpec:  pointers.StringPtr("0"),
						RetryLimit: pointers.Int64Ptr(2),
						Template: cv1b1.JobPodTemplate{
							Containers: []corev1.Container{
								{
									Name:  "fake-container",
									Image: "fake-image",
									Resources: corev1.ResourceRequirements{
										Requests: corev1.ResourceList{
											corev1.ResourceCPU:    resource.MustParse("1"),
											corev1.ResourceMemory: resource.MustParse("64Mi"),
										},
									},
								},
							},
						},
					},
				}
			})

			It("should not throw an error", func() {
				Expect(len(err)).Should(Equal(0))
			})
		})
	})

	When("validating JobDefinitionSpec", func() {
		var (
			jobDefinitionSpec *cv1b1.JobDefinitionSpec
			err               field.ErrorList
		)

		BeforeEach(func() {
			jobDefinitionSpec = &cv1b1.JobDefinitionSpec{
				ArraySpec:  pointers.StringPtr("0"),
				RetryLimit: pointers.Int64Ptr(2),
				Template: cv1b1.JobPodTemplate{
					Containers: []corev1.Container{
						{
							Name:  "fake-container",
							Image: "fake-image",
						},
					},
				},
			}
		})

		JustBeforeEach(func() {
			err = cv1b1.ValidateJobDefinitionSpec(jobDefinitionSpec, validationConfig, field.NewPath("spec"))
		})

		Context("that is valid", func() {
			BeforeEach(func() {
				jobDefinitionSpec = &cv1b1.JobDefinitionSpec{
					ArraySpec:  pointers.StringPtr("0"),
					RetryLimit: pointers.Int64Ptr(2),
					Template: cv1b1.JobPodTemplate{
						Containers: []corev1.Container{
							{
								Name:  "fake-container",
								Image: "fake-image",
							},
						},
					},
				}
			})
			It("should not throw an error", func() {
				Expect(len(err)).Should(Equal(0))
			})
		})

		Context("that does not define template", func() {
			BeforeEach(func() {
				jobDefinitionSpec = &cv1b1.JobDefinitionSpec{
					ArraySpec:  pointers.StringPtr("0"),
					RetryLimit: pointers.Int64Ptr(2),
					Template:   cv1b1.JobPodTemplate{},
				}
			})
			It("returns a meaningful error message", func() {
				Expect(errors.AggregateErrorList(err)).Should(
					ContainSubstring("spec.template: Missing field: must be specified"),
				)
			})
		})

		Describe("ArraySpec", func() {
			Context("that does not define ArraySpec", func() {
				BeforeEach(func() {
					jobDefinitionSpec.ArraySpec = nil
				})

				It("invalid JobDefinitionSpec due to missing arraySpec", func() {
					Expect(errors.AggregateErrorList(err)).Should(
						ContainSubstring("spec.arraySpec: Missing field: must be specified"),
					)
				})

			})

			DescribeTable("that defines an invalid ArraySpec",
				func(arraySpec string, expectedErrorMessage string) {
					jobDefinitionSpec.ArraySpec = pointers.StringPtr(arraySpec)
					Expect(errors.AggregateErrorList(cv1b1.ValidateJobDefinitionSpec(jobDefinitionSpec, validationConfig, field.NewPath("spec")))).
						Should(ContainSubstring(expectedErrorMessage))
				},
				Entry("which is empty", "", "error getting start index of range '':"),
				Entry("which contains invalid range notation", "1-3-5, 7", "error parsing arraySpec range: '1-3-5'. Expect 2, got 3"),
				Entry("which contains alpha numerical value", "a-3", "error getting start index of range 'a-3':"),
				Entry("which contains alpha numerical value", "1-b", "error getting end index of range '1-b':"),
				Entry("which exceeds the limit range", "0-2147483648", "exceeded index range, must between 0 and 9999999, got: 0 and 2147483648"),
				Entry("which exceeds maximum arrray size", "0-500, 600-1200", "spec.arraySpec: Invalid value: \"0-500, 600-1200\": exceeded maximum array size: 1000"),
			)

		})

		Describe("RetryLimit", func() {
			Context("that does not define RetryLimit", func() {
				BeforeEach(func() {
					jobDefinitionSpec.RetryLimit = nil
				})
				It("returns a meaningful error message", func() {
					Expect(errors.AggregateErrorList(err)).Should(
						ContainSubstring("spec.retryLimit: Missing field: must be specified"),
					)
				})
			})
		})
	})

	When("validating JobPodTemplate", func() {
		var (
			jobPodTemplate *cv1b1.JobPodTemplate
		)

		BeforeEach(func() {
			jobPodTemplate = &cv1b1.JobPodTemplate{
				Containers: []corev1.Container{
					{
						Name:  "fake-container",
						Image: "fake-image",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("1"),
								corev1.ResourceMemory: resource.MustParse("64Mi"),
							},
						},
					},
				},
			}
		})

		It("valid JobDefinitionSpec", func() {
			err := cv1b1.ValidateJobPodTemplate(jobPodTemplate, validationConfig, field.NewPath("template"))
			Expect(len(err)).Should(Equal(0))
		})

		It("invalid JobDefinitionSpec due to missing containers with empty jobDefinitionRef", func() {
			jobPodTemplate = &cv1b1.JobPodTemplate{
				ServiceAccountName: "fake-serviceaccount",
			}
			err := cv1b1.ValidateJobPodTemplate(jobPodTemplate, validationConfig, field.NewPath("template"))
			Expect(len(err)).ShouldNot(Equal(0))
			Expect(errors.AggregateErrorList(err)).Should(ContainSubstring("template.containers: Unsupported container count: 0: there must be exactly one container"))
		})

		It("invalid JobDefinitionSpec due to missing container fields", func() {
			jobPodTemplate.Containers[0].Name = ""
			jobPodTemplate.Containers[0].Image = ""
			err := cv1b1.ValidateJobPodTemplate(jobPodTemplate, validationConfig, field.NewPath("template"))
			Expect(len(err)).ShouldNot(Equal(0))
			Expect(errors.AggregateErrorList(err)).Should(ContainSubstring("template.containers[0].name: Missing field: must be specified"))
			Expect(errors.AggregateErrorList(err)).Should(ContainSubstring("template.containers[0].image: Missing field: must be specified"))
		})

		It("invalid JobDefinitionSpec due to disallowed container fields", func() {
			jobPodTemplate.Containers[0].Stdin = true
			jobPodTemplate.Containers[0].StdinOnce = true
			err := cv1b1.ValidateJobPodTemplate(jobPodTemplate, validationConfig, field.NewPath("template"))
			Expect(len(err)).ShouldNot(Equal(0))
			Expect(errors.AggregateErrorList(err)).Should(ContainSubstring("template.containers[0].stdin: Unsupported field"))
			Expect(errors.AggregateErrorList(err)).Should(ContainSubstring("template.containers[0].stdinOnce: Unsupported field"))
		})

		It("invalid JobDefinitionSpec due to disallowed envVarSource fields", func() {
			jobPodTemplate.Containers[0].Env = []corev1.EnvVar{
				{
					Name: "MY_NODE_NAME",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "spec.nodeName",
						},
					},
				},
				{
					Name: "MY_CPU_REQUEST",
					ValueFrom: &corev1.EnvVarSource{
						ResourceFieldRef: &corev1.ResourceFieldSelector{
							ContainerName: "fake-container",
							Resource:      "requests.cpu",
						},
					},
				},
			}
			err := cv1b1.ValidateJobPodTemplate(jobPodTemplate, validationConfig, field.NewPath("template"))
			Expect(len(err)).ShouldNot(Equal(0))
			Expect(errors.AggregateErrorList(err)).Should(ContainSubstring("template.containers[0].env[0].valueFrom.fieldRef: Unsupported field"))
			Expect(errors.AggregateErrorList(err)).Should(ContainSubstring("template.containers[0].env[1].valueFrom.resourceFieldRef: Unsupported field"))
		})

		It("invalid JobDefinitionSpec due to invalid (negative) resource requirements", func() {
			jobPodTemplate.Containers[0].Resources = corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("-10"),
					corev1.ResourceMemory: resource.MustParse("-64"),
				},
			}
			err := cv1b1.ValidateJobPodTemplate(jobPodTemplate, validationConfig, field.NewPath("template"))
			Expect(len(err)).ShouldNot(Equal(0))
			Expect(errors.AggregateErrorList(err)).Should(ContainSubstring("template.containers[0].resources.requests.cpu: Out of range: -10: must be greater than 0"))
			Expect(errors.AggregateErrorList(err)).Should(ContainSubstring("template.containers[0].resources.requests.memory: Out of range: -64: must be greater than 0"))
		})

		It("invalid JobDefinitionSpec due to invalid (zero) resource requirements", func() {
			jobPodTemplate.Containers[0].Resources = corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("0"),
					corev1.ResourceMemory: resource.MustParse("0"),
				},
			}
			err := cv1b1.ValidateJobPodTemplate(jobPodTemplate, validationConfig, field.NewPath("template"))
			Expect(len(err)).ShouldNot(Equal(0))
			Expect(errors.AggregateErrorList(err)).Should(ContainSubstring("template.containers[0].resources.requests.cpu: Out of range: 0: must be greater than 0"))
			Expect(errors.AggregateErrorList(err)).Should(ContainSubstring("template.containers[0].resources.requests.memory: Out of range: 0: must be greater than 0"))
		})

		It("invalid JobDefinitionSpec due to two containers", func() {
			jobPodTemplate.Containers = []corev1.Container{
				{
					Name:  "fake-container1",
					Image: "fake-image",
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("1"),
							corev1.ResourceMemory: resource.MustParse("64Mi"),
						},
					},
				},
				{
					Name:  "fake-container2",
					Image: "fake-image",
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("1"),
							corev1.ResourceMemory: resource.MustParse("64Mi"),
						},
					},
				},
			}

			err := cv1b1.ValidateJobPodTemplate(jobPodTemplate, validationConfig, field.NewPath("template"))
			Expect(len(err)).ShouldNot(Equal(0))
			Expect(errors.AggregateErrorList(err)).Should(ContainSubstring("template.containers: Unsupported container count: 2: there must be exactly one container"))
		})

		It("invalid JobDefinitionSpec due to invalid env var", func() {
			jobPodTemplate.Containers[0].Env = []corev1.EnvVar{
				{
					Name:  cv1b1.JobIndex,
					Value: "fake-index",
				},
			}

			err := cv1b1.ValidateJobPodTemplate(jobPodTemplate, validationConfig, field.NewPath("template"))
			Expect(len(err)).ShouldNot(Equal(0))
			Expect(errors.AggregateErrorList(err)).Should(ContainSubstring(fmt.Sprintf("'%s' is reserved for batch API", cv1b1.JobIndex)))
		})

		It("invalid JobDefinitionSpec due to invalid container name", func() {
			jobPodTemplate.Containers[0].Name = "fake_name"

			err := cv1b1.ValidateJobPodTemplate(jobPodTemplate, validationConfig, field.NewPath("template"))
			Expect(len(err)).ShouldNot(Equal(0))
			Expect(errors.AggregateErrorList(err)).Should(SatisfyAll(
				ContainSubstring("invalid value: a lowercase RFC 1123 label must consist of lower case alphanumeric characters"),
				ContainSubstring("containers[0].name"),
			))
		})

		Context("JobPodTemplate size", func() {
			const sizeLimit = 10 * 1024
			var errs field.ErrorList

			BeforeEach(func() {
				// Get current size
				jobPodTemplate.ServiceAccountName = "x"
				bytes, _ := json.Marshal(jobPodTemplate)
				size := len(bytes)

				// Max the size
				jobPodTemplate.ServiceAccountName = jobPodTemplate.ServiceAccountName + strings.Repeat("x", sizeLimit-size)
			})

			JustBeforeEach(func() {
				errs = cv1b1.ValidateJobPodTemplate(jobPodTemplate, validationConfig, field.NewPath("template"))
			})

			When("does not exceed limit", func() {
				It("does not return errors", func() {
					Expect(errs).To(BeEmpty())
				})
			})

			When("exceeds limit", func() {
				BeforeEach(func() {
					jobPodTemplate.ServiceAccountName = jobPodTemplate.ServiceAccountName + "x"
				})

				It("returns errors", func() {
					Expect(errs).ToNot(BeEmpty())
				})

				It("returns ErrSizeLimitExceeded", func() {
					Expect(string(errs[0].Type)).To(Equal("ErrSizeLimitExceeded"))
				})

				It("should contain descriptive error message", func() {
					Expect(errors.AggregateErrorList(errs)).Should(SatisfyAll(
						ContainSubstring("template: Exceeding size limit: size limit is 10240 bytes"),
						ContainSubstring("but actual size is"),
					))
				})
			})
		})
	})

	When("validating a JobDefinition", func() {
		var (
			jobDef *cv1b1.JobDefinition
			err    field.ErrorList
		)
		BeforeEach(func() {})

		JustBeforeEach(func() {
			err = jobDef.Validate(validationConfig)
		})

		//Validation: if any of cpu, ephemeral-storage or memory is specified under JobDefinition.spec.template.containers.resources.requests:
		Context("that specifies containers.resources", func() {
			BeforeEach(func() {
				jobDef = &cv1b1.JobDefinition{
					Spec: cv1b1.JobDefinitionSpec{
						ArraySpec:  pointers.StringPtr("0"),
						RetryLimit: pointers.Int64Ptr(2),
						Template: cv1b1.JobPodTemplate{
							Containers: []corev1.Container{
								{
									Name:  "fake-container",
									Image: "fake-image",
									Resources: corev1.ResourceRequirements{
										Requests: corev1.ResourceList{
											corev1.ResourceCPU:    resource.MustParse("500m"),
											corev1.ResourceMemory: resource.MustParse("250M"),
										},
									},
								},
							},
						},
					},
				}
			})

			Context("and there is no LimitRange given", func() {
				BeforeEach(func() {
					validationConfig.ContainerLimitRangeItem = nil
				})

				It("does not error", func() {
					Expect(len(err)).Should(Equal(0))
				})
			})

			Context("and there is a LimitRange given", func() {
				Context("and JobDefinition specifies request.cpu", func() {
					BeforeEach(func() {
						validationConfig.ContainerLimitRangeItem = &corev1.LimitRangeItem{
							Type: corev1.LimitTypeContainer,
							Min: corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("10m"),
							},
							Max: corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("8"),
							},
						}
					})

					Context("which is within the bounds of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("1"),
							}
						})

						It("does not throw a validation error", func() {
							Expect(len(err)).Should(Equal(0))
						})
					})

					Context("which exceeds the upper bound of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("9"),
							}
						})

						It("returns a meaningful error message", func() {
							Expect(errors.AggregateErrorList(err)).To(
								SatisfyAll(
									ContainSubstring("expected 10m <= 9 <= 8"),
									ContainSubstring("spec.template.containers[0].resources.requests.cpu"),
								),
							)
						})
					})

					Context("which equals the upper bound of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("8"),
							}
						})

						It("does not throw a validation error", func() {
							Expect(len(err)).Should(Equal(0))
						})
					})

					Context("which equals the lower bound of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("10m"),
							}
						})

						It("does not throw a validation error", func() {
							Expect(len(err)).Should(Equal(0))
						})
					})

					Context("which falls below the lower bound of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("9m"),
							}
						})

						It("returns a meaningful error message", func() {
							Expect(errors.AggregateErrorList(err)).To(
								SatisfyAll(
									ContainSubstring("expected 10m <= 9m <= 8"),
									ContainSubstring("spec.template.containers[0].resources.requests.cpu"),
								),
							)
						})
					})
				})

				Context("and JobDefinition specifies request.ephemeral-storage", func() {
					BeforeEach(func() {
						validationConfig.ContainerLimitRangeItem = &corev1.LimitRangeItem{
							Type: corev1.LimitTypeContainer,
							Min: corev1.ResourceList{
								corev1.ResourceEphemeralStorage: resource.MustParse("100Mi"),
							},
							Max: corev1.ResourceList{
								corev1.ResourceEphemeralStorage: resource.MustParse("4Gi"),
							},
						}
					})

					Context("which is within the bounds of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceEphemeralStorage: resource.MustParse("2Gi"),
							}
						})

						It("does not throw a validation error", func() {
							Expect(len(err)).Should(Equal(0))
						})
					})

					Context("which equals the upper bound of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceEphemeralStorage: resource.MustParse("4Gi"),
							}
						})

						It("does not throw a validation error", func() {
							Expect(len(err)).Should(Equal(0))
						})
					})

					Context("which exceeds the upper bound of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceEphemeralStorage: resource.MustParse("4097Mi"),
							}
						})

						It("returns a meaningful error message", func() {
							Expect(errors.AggregateErrorList(err)).To(
								SatisfyAll(
									ContainSubstring("expected 100Mi <= 4097Mi <= 4Gi"),
									ContainSubstring("spec.template.containers[0].resources.requests.ephemeral-storage"),
								),
							)
						})
					})

					Context("which equals the lower bound of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceEphemeralStorage: resource.MustParse("100Mi"),
							}
						})

						It("does not throw a validation error", func() {
							Expect(len(err)).Should(Equal(0))
						})
					})

					Context("which falls below the lower bound of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceEphemeralStorage: resource.MustParse("102399Ki"),
							}
						})

						It("returns a meaningful error message", func() {
							Expect(errors.AggregateErrorList(err)).To(
								SatisfyAll(
									ContainSubstring("expected 100Mi <= 102399Ki <= 4Gi"),
									ContainSubstring("spec.template.containers[0].resources.requests.ephemeral-storage"),
								),
							)
						})
					})
				})

				Context("and JobDefinition specifies request.memory", func() {
					BeforeEach(func() {
						validationConfig.ContainerLimitRangeItem = &corev1.LimitRangeItem{
							Type: corev1.LimitTypeContainer,
							Min: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("128Mi"),
							},
							Max: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("32Gi"),
							},
						}
					})

					Context("which is within the bounds of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("512Mi"),
							}
						})

						It("does not throw a validation error", func() {
							Expect(len(err)).Should(Equal(0))
						})
					})

					Context("which equals the upper bound of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("32768Mi"),
							}
						})

						It("does not throw a validation error", func() {
							Expect(len(err)).Should(Equal(0))
						})
					})

					Context("which exceeds the upper bound of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("32769Mi"),
							}
						})

						It("returns a meaningful error message", func() {
							Expect(errors.AggregateErrorList(err)).To(
								SatisfyAll(
									ContainSubstring("expected 128Mi <= 32769Mi <= 32Gi"),
									ContainSubstring("spec.template.containers[0].resources.requests.memory"),
								),
							)
						})
					})

					Context("which equals the lower bound of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("131072Ki"),
							}
						})

						It("does not throw a validation error", func() {
							Expect(len(err)).Should(Equal(0))
						})
					})

					Context("which falls below the lower bound of LimitRange for the namespace", func() {
						BeforeEach(func() {
							jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("131071Ki"),
							}
						})

						It("returns a meaningful error message", func() {
							Expect(errors.AggregateErrorList(err)).To(
								SatisfyAll(
									ContainSubstring("expected 128Mi <= 131071Ki <= 32Gi"),
									ContainSubstring("spec.template.containers[0].resources.requests.memory"),
								),
							)
						})
					})
				})

				Context("with incomplete min/max values and a ResourceQuantity that's out of the limit", func() {
					BeforeEach(func() {
						validationConfig.ContainerLimitRangeItem = &corev1.LimitRangeItem{
							Type: corev1.LimitTypeContainer,
							Min: corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("10m"),
							},
						}
						jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
							corev1.ResourceCPU: resource.MustParse("9m"),
						}
					})

					It("returns a meaningful error message", func() {
						Expect(errors.AggregateErrorList(err)).To(
							SatisfyAll(
								ContainSubstring("expected 10m <= 9m <= n/a"),
								ContainSubstring("spec.template.containers[0].resources.requests.cpu"),
							))
					})
				})
			})

			Context("and there is no ResourceRatios given", func() {
				BeforeEach(func() {
					validationConfig.ResourceRatios = nil
				})

				It("does not error", func() {
					Expect(len(err)).Should(Equal(0))
				})
			})

			Context("and there is ResourceRatios given with zero elements", func() {
				BeforeEach(func() {
					validationConfig.ResourceRatios = make(cv1b1.ResourceRatios, 0)
				})

				It("does not error", func() {
					Expect(len(err)).Should(Equal(0))
				})
			})

			Context("and there is one ResourceRatio given", func() {
				var allowedRatios cv1b1.ResourceRatios
				BeforeEach(func() {
					allowedRatios = cv1b1.ResourceRatios{
						{
							corev1.ResourceCPU:    resource.MustParse("125m"),
							corev1.ResourceMemory: resource.MustParse("250M"),
						},
					}
				})

				DescribeTable("and pod specifies resource request",
					func(c RatioCase) {
						RatioTestFunc(allowedRatios, c)
					},
					Entry(RatioTestDescription(), RatioCase{
						// It's also valid to define no Resources -> will get defaulted
						specifiedResourceRequests: corev1.ResourceList{},
						isValid:                   true,
					}),
					Entry(RatioTestDescription(), RatioCase{
						// It's not valid to partially define resources
						specifiedResourceRequests: corev1.ResourceList{
							corev1.ResourceCPU: resource.MustParse("125m"),
						},
						isValid: false,
					}),
					Entry(RatioTestDescription(), RatioCase{
						specifiedResourceRequests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("125m"),
							corev1.ResourceMemory: resource.MustParse("250M"),
						},
						isValid: true,
					}),
					Entry(RatioTestDescription(), RatioCase{
						specifiedResourceRequests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("250m"),
							corev1.ResourceMemory: resource.MustParse("500M"),
						},
						isValid: false,
					}),
					Entry(RatioTestDescription(), RatioCase{
						specifiedResourceRequests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("500m"),
							corev1.ResourceMemory: resource.MustParse("1024M"),
						},
						isValid: false,
					}),
					Entry(RatioTestDescription(), RatioCase{
						specifiedResourceRequests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("500m"),
							corev1.ResourceMemory: resource.MustParse("500M"),
						},
						isValid: false,
					}),
					Entry(RatioTestDescription(), RatioCase{
						specifiedResourceRequests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("0"),
							corev1.ResourceMemory: resource.MustParse("0"),
						},
						isValid: false,
					}),
				)
			})

			Context("and there are multiple ResourceRatios given", func() {
				var allowedRatios cv1b1.ResourceRatios
				BeforeEach(func() {
					allowedRatios = cv1b1.ResourceRatios{
						{
							corev1.ResourceCPU:    resource.MustParse("125m"),
							corev1.ResourceMemory: resource.MustParse("250M"),
						},
						{
							corev1.ResourceCPU:    resource.MustParse("250m"),
							corev1.ResourceMemory: resource.MustParse("500M"),
						},
						{
							corev1.ResourceCPU:    resource.MustParse("500m"),
							corev1.ResourceMemory: resource.MustParse("1024M"),
						},
					}
				})

				DescribeTable("and pod specifies resource request",
					func(c RatioCase) {
						RatioTestFunc(allowedRatios, c)
					},
					Entry(RatioTestDescription(), RatioCase{
						// It's also valid to define no Resources -> will get defaulted
						specifiedResourceRequests: corev1.ResourceList{},
						isValid:                   true,
					}),
					Entry(RatioTestDescription(), RatioCase{
						// It's not valid to partially define resources
						specifiedResourceRequests: corev1.ResourceList{
							corev1.ResourceCPU: resource.MustParse("125m"),
						},
						isValid: false,
					}),
					Entry(RatioTestDescription(), RatioCase{
						specifiedResourceRequests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("125m"),
							corev1.ResourceMemory: resource.MustParse("250M"),
						},
						isValid: true,
					}),
					Entry(RatioTestDescription(), RatioCase{
						specifiedResourceRequests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("250m"),
							corev1.ResourceMemory: resource.MustParse("500M"),
						},
						isValid: true,
					}),
					Entry(RatioTestDescription(), RatioCase{
						specifiedResourceRequests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("500m"),
							corev1.ResourceMemory: resource.MustParse("1024M"),
						},
						isValid: true,
					}),
					Entry(RatioTestDescription(), RatioCase{
						specifiedResourceRequests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("500m"),
							corev1.ResourceMemory: resource.MustParse("500M"),
						},
						isValid: false,
					}),
					Entry(RatioTestDescription(), RatioCase{
						specifiedResourceRequests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("0"),
							corev1.ResourceMemory: resource.MustParse("0"),
						},
						isValid: false,
					}),
				)
			})
		})
	})
})

type RatioCase struct {
	specifiedResourceRequests corev1.ResourceList
	isValid                   bool
}

func RatioTestDescription() func(RatioCase) string {
	return func(testcase RatioCase) string {
		if len(testcase.specifiedResourceRequests) == 0 {
			return "with no resources defined"
		}
		return fmt.Sprintf("with CPU: %s and memory: %s",
			testcase.specifiedResourceRequests.Cpu().String(),
			testcase.specifiedResourceRequests.Cpu().String(),
		)
	}
}

func RatioTestFunc(allowedRatios cv1b1.ResourceRatios, c RatioCase) {
	jobDef := &cv1b1.JobDefinition{
		Spec: cv1b1.JobDefinitionSpec{
			ArraySpec:  pointers.StringPtr("0"),
			RetryLimit: pointers.Int64Ptr(2),
			Template: cv1b1.JobPodTemplate{
				Containers: []corev1.Container{
					{
						Name:  "fake-container",
						Image: "fake-image",
						Resources: corev1.ResourceRequirements{
							Requests: c.specifiedResourceRequests,
						},
					},
				},
			},
		},
	}

	validationConf := cv1b1.ValidationConfig{
		ContainerLimitRangeItem: nil,
		ResourceRatios:          allowedRatios,
	}
	err := jobDef.Validate(validationConf)

	if c.isValid {
		Expect(err).Should(HaveLen(0), "no validation errors expected")
	} else {
		Expect(errors.AggregateErrorList(err)).To(
			SatisfyAll(
				ContainSubstring("resources.requests"),
				ContainSubstring("Unsupported resource specification"),
			))
	}
}
