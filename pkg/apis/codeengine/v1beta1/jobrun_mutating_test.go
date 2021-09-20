/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1_test

import (
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	cv1b1 "github.ibm.com/coligo/batch-job-controller/pkg/apis/codeengine/v1beta1"
	"github.ibm.com/coligo/batch-job-controller/pkg/utils/errors"
	"github.ibm.com/coligo/batch-job-controller/pkg/utils/pointers"
)

type ByName []corev1.EnvVar

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

var _ = Describe("Mutating JobRun", func() {
	Describe("JobRun without a reference to a JobDefinition", func() {
		var (
			jobRun               *cv1b1.JobRun
			jobRunBeforeMutation *cv1b1.JobRun
		)

		BeforeEach(func() {
			jobRun = &cv1b1.JobRun{}
			jobRunBeforeMutation = &cv1b1.JobRun{}
		})

		JustBeforeEach(func() {
			jobRun.DeepCopyInto(jobRunBeforeMutation)
			jobRun.MutateWithDefaults()
		})

		Describe("ArraySpec", func() {
			Context("with a JobRun that did not specify ArraySpec", func() {
				BeforeEach(func() {
					jobRun.Spec.JobDefinitionSpec.ArraySpec = nil
				})

				It("JobDefinitionSpec contains the default value for ArraySpec", func() {
					Expect(*jobRun.Spec.JobDefinitionSpec.ArraySpec).To(Equal("0"))
				})
			})

			Context("with a JobRun that specified ArraySpec", func() {
				BeforeEach(func() {
					jobRun.Spec.JobDefinitionSpec.ArraySpec = pointers.StringPtr("0-1")
				})

				It("JobDefinitionSpec contains the ArraySpec value from jobRun.Spec.JobDefinitionSpec", func() {
					Expect(*jobRun.Spec.JobDefinitionSpec.ArraySpec).To(Equal("0-1"))
				})
			})
		})

		Describe("RetryLimit", func() {
			Context("with a JobRun that did not specify RetryLimit", func() {
				BeforeEach(func() {
					jobRun.Spec.JobDefinitionSpec.RetryLimit = nil
				})

				It("JobDefinitionSpec contains the default value for RetryLimit", func() {
					Expect(*jobRun.Spec.JobDefinitionSpec.RetryLimit).To(BeNumerically("==", 3))
				})
			})

			Context("with a JobRun that specified RetryLimit", func() {
				BeforeEach(func() {
					jobRun.Spec.JobDefinitionSpec.RetryLimit = pointers.Int64Ptr(4)
				})

				It("JobDefinitionSpec contains the RetryLimit value from jobRun.Spec.JobDefinitionSpec", func() {
					Expect(*jobRun.Spec.JobDefinitionSpec.RetryLimit).To(BeNumerically("==", 4))
				})
			})
		})

		Describe("MaxExecutionTime", func() {
			Context("with a JobRun that did not specify MaxExecutionTime", func() {
				BeforeEach(func() {
					jobRun.Spec.JobDefinitionSpec.MaxExecutionTime = nil
				})

				It("JobDefinitionSpec contains the default value for MaxExecutionTime", func() {
					Expect(*jobRun.Spec.JobDefinitionSpec.MaxExecutionTime).To(BeNumerically("==", 7200))
				})
			})

			Context("with a JobRun that specified MaxExecutionTime", func() {
				BeforeEach(func() {
					jobRun.Spec.JobDefinitionSpec.MaxExecutionTime = pointers.Int64Ptr(5)
				})

				It("JobDefinitionSpec contains the MaxExecutionTime value from jobRun.Spec.JobDefinitionSpec", func() {
					Expect(*jobRun.Spec.JobDefinitionSpec.MaxExecutionTime).To(BeNumerically("==", 5))
				})
			})
		})

		Describe("Jobrun GenerateName", func() {
			When("the length of jobrun generateName exceeds 'maxGenerateNameLength'", func() {
				BeforeEach(func() {
					jobRun.ObjectMeta.GenerateName = "jobrun-with-generate-name-that-has-more-characters-than-jobRunNameMaxLen-and-hence-should-get-trimmed"
				})

				It("jobRun generateName should get trimmed to 'maxGenerateNameLength'", func() {
					Expect(jobRun.ObjectMeta.GenerateName).To(HaveLen(48))
				})
			})

			When("the length of jobrun generateName is below 'maxGenerateNameLength'", func() {
				var validJobRunGenerateName = "jobrun-with-short-name"

				BeforeEach(func() {
					jobRun.ObjectMeta.GenerateName = validJobRunGenerateName
				})

				It("jobRun name should remain the same", func() {
					Expect(jobRun.ObjectMeta.GenerateName).To(Equal(validJobRunGenerateName))
				})
			})
		})

		It("passes the template unmodified", func() {
			Expect(jobRun.Spec.JobDefinitionSpec.Template).To(Equal(jobRunBeforeMutation.Spec.JobDefinitionSpec.Template))
		})
	})

	Describe("JobRun with an explicit reference to a JobDefinition", func() {
		var (
			jobRun *cv1b1.JobRun
			jobDef *cv1b1.JobDefinition
			err    field.ErrorList
		)

		BeforeEach(func() {
			jobRun = &cv1b1.JobRun{}
			jobDef = &cv1b1.JobDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name: "JobDef-Fake-Name",
					UID:  "JobDef-Fake-UID",
				},
				Spec: cv1b1.JobDefinitionSpec{
					ArraySpec:        pointers.StringPtr("0-1"),
					MaxExecutionTime: pointers.Int64Ptr(5000),
					RetryLimit:       pointers.Int64Ptr(4),
					Template: cv1b1.JobPodTemplate{
						Containers: []corev1.Container{
							{
								Name:  "Foo",
								Image: "Bar",
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

		JustBeforeEach(func() {
			err = jobRun.MutateWithJobDefinition(jobDef.DeepCopy())
		})

		Describe("Meta data", func() {
			It("the Jobrun has the JobDefinition's Name label", func() {
				Expect(jobRun.Labels[cv1b1.LabelJobDefName]).To(Equal("JobDef-Fake-Name"))
			})

			It("the Jobrun has the JobDefinition's UID label", func() {
				Expect(jobRun.Labels[cv1b1.LabelJobDefUUID]).To(Equal("JobDef-Fake-UID"))
			})

			It("the Jobrun has the JobDefinition as controller", func() {
				Expect(string(metav1.GetControllerOf(jobRun).UID)).To(Equal("JobDef-Fake-UID"))
			})
		})

		Describe("JobDefinitionSpec", func() {
			Context("with a JobRun that has a JobDefinitionSpec with zero value", func() {
				BeforeEach(func() {
					jobRun.Spec.JobDefinitionSpec = cv1b1.JobDefinitionSpec{}
				})

				It("mutates", func() {
					Expect(len(err)).Should(Equal(0))
				})

				It("JobDefinitionSpec equals the jobDef.Spec", func() {
					Expect(jobRun.Spec.JobDefinitionSpec).To(Equal(jobDef.Spec))
				})
			})

			Context("with a JobRun that has a non-empty JobDefinitionSpec", func() {
				BeforeEach(func() {
					jobRun.Spec.JobDefinitionSpec = cv1b1.JobDefinitionSpec{
						ArraySpec:        pointers.StringPtr("0-3"),
						MaxExecutionTime: pointers.Int64Ptr(7234),
						RetryLimit:       pointers.Int64Ptr(0),
						Template: cv1b1.JobPodTemplate{
							Containers: []corev1.Container{
								{
									Resources: corev1.ResourceRequirements{
										Requests: corev1.ResourceList{
											corev1.ResourceCPU:    resource.MustParse("2"),
											corev1.ResourceMemory: resource.MustParse("128Mi"),
										},
									},
								},
							},
						},
					}
				})

				Describe("ArraySpec", func() {
					//jobDef.Spec will always have a Arrayspec value because we use CRD defaulting.
					//Therefore we will fill it here explicitly with some value
					BeforeEach(func() {
						jobDef.Spec.ArraySpec = pointers.StringPtr("0-42")
					})

					Context("with a JobRun that did not specify jobRun.Spec.JobDefinitionSpec.ArraySpec", func() {
						BeforeEach(func() {
							jobRun.Spec.JobDefinitionSpec.ArraySpec = nil
						})

						It("mutates", func() {
							Expect(len(err)).Should(Equal(0))
						})

						It("JobDefinitionSpec contains given value from jobDef.Spec.ArraySpec", func() {
							Expect(*jobRun.Spec.JobDefinitionSpec.ArraySpec).To(Equal("0-42"))
						})
					})

					Context("where the JobRun specified jobRun.Spec.JobDefinitionSpec.ArraySpec", func() {
						BeforeEach(func() {
							jobRun.Spec.JobDefinitionSpec.ArraySpec = pointers.StringPtr("11")
						})

						It("mutates", func() {
							Expect(len(err)).Should(Equal(0))
						})

						It("JobDefinitionSpec contains the ArraySpec value from jobRun.Spec.JobDefinitionSpec", func() {
							Expect(*jobRun.Spec.JobDefinitionSpec.ArraySpec).To(Equal("11"))
						})
					})
				})

				Describe("RetryLimit", func() {
					//jobDef.Spec will always have a RetryLimit value because we use CRD defaulting.
					//Therefore we have to fill it here explicitly with some value
					BeforeEach(func() {
						jobDef.Spec.RetryLimit = pointers.Int64Ptr(5)
					})

					Context("with a JobRun that did not specify jobRun.Spec.JobDefinitionSpec.RetryLimit", func() {
						BeforeEach(func() {
							jobRun.Spec.JobDefinitionSpec.RetryLimit = nil
						})

						It("mutates", func() {
							Expect(len(err)).Should(Equal(0))
						})

						It("JobDefinitionSpec contains given value from jobDef.Spec.RetryLimit", func() {
							Expect(*jobRun.Spec.JobDefinitionSpec.RetryLimit).To(BeNumerically("==", 5))
						})
					})

					Context("with a JobRun that specified jobRun.Spec.JobDefinitionSpec.RetryLimit", func() {
						BeforeEach(func() {
							jobRun.Spec.JobDefinitionSpec.RetryLimit = pointers.Int64Ptr(1)
						})

						It("mutates", func() {
							Expect(len(err)).Should(Equal(0))
						})

						It("JobDefinitionSpec contains the RetryLimit value from jobRun.Spec.JobDefinitionSpec", func() {
							Expect(*jobRun.Spec.JobDefinitionSpec.RetryLimit).To(BeNumerically("==", 1))
						})
					})
				})

				Describe("MaxExecutionTime", func() {
					//jobDef.Spec will always have a Arrayspec value because we use CRD defaulting.
					//Therefore we will fill it here explicitly with some value
					BeforeEach(func() {
						jobDef.Spec.MaxExecutionTime = pointers.Int64Ptr(1234)
					})

					Context("with a JobRun that did not specify jobRun.Spec.JobDefinitionSpec.MaxExecutionTime", func() {
						BeforeEach(func() {
							jobRun.Spec.JobDefinitionSpec.MaxExecutionTime = nil
						})

						It("mutates", func() {
							Expect(len(err)).Should(Equal(0))
						})

						It("JobDefinitionSpec contains given value from jobDef.Spec.MaxExecutionTime", func() {
							Expect(*jobRun.Spec.JobDefinitionSpec.MaxExecutionTime).To(BeNumerically("==", 1234))
						})
					})

					Context("with a JobRun that specified jobRun.Spec.JobDefinitionSpec.MaxExecutionTime", func() {
						BeforeEach(func() {
							jobRun.Spec.JobDefinitionSpec.MaxExecutionTime = pointers.Int64Ptr(4321)
						})

						It("mutates", func() {
							Expect(len(err)).Should(Equal(0))
						})

						It("JobDefinitionSpec contains the MaxExecutionTime value from jobRun.Spec.JobDefinitionSpec", func() {
							Expect(*jobRun.Spec.JobDefinitionSpec.MaxExecutionTime).To(BeNumerically("==", 4321))
						})
					})
				})

				Describe("Jobrun GenerateName", func() {
					When("the length of jobrun GenerateName exceeds 'maxGenerateNameLength'", func() {
						BeforeEach(func() {
							jobRun.ObjectMeta.GenerateName = "jobrun-with-generate-name-that-has-more-characters-than-jobRunNameMaxLen-and-hence-should-get-trimmed"
						})

						It("jobRun GenerateName should get trimmed to 'maxGenerateNameLength'", func() {
							Expect(jobRun.ObjectMeta.GenerateName).To(HaveLen(48))
						})
					})

					When("the length of jobrun GenerateName is below 'maxGenerateNameLength'", func() {
						var validJobRunGenerateName = "jobrun-with-short-name"

						BeforeEach(func() {
							jobRun.ObjectMeta.GenerateName = validJobRunGenerateName
						})

						It("jobRun GenerateName should remain the same", func() {
							Expect(jobRun.ObjectMeta.GenerateName).To(Equal(validJobRunGenerateName))
						})
					})
				})

				Describe("Template", func() {
					Describe("ImagePullSecrets", func() {
						Context("with a JobRun that did not specify jobRun.Spec.JobDefinitionSpec.Template.ImagePullSecrets", func() {
							BeforeEach(func() {
								jobRun.Spec.JobDefinitionSpec.Template.ImagePullSecrets = nil
							})

							Context("where jobDef.Spec.Template has ImagePullSecrets", func() {
								BeforeEach(func() {
									jobDef.Spec.Template.ImagePullSecrets = []corev1.LocalObjectReference{
										{
											Name: "JobDefImagePullSecret",
										},
									}
								})

								It("mutates", func() {
									Expect(len(err)).Should(Equal(0))
								})

								It("JobDefinitionSpec contains given ImagePullSecrets from jobDef.Spec.Template.ImagePullSecrets", func() {
									Expect(jobRun.Spec.JobDefinitionSpec.Template.ImagePullSecrets).To(Equal([]corev1.LocalObjectReference{
										{
											Name: "JobDefImagePullSecret",
										},
									}))
								})
							})

							Context("where jobDef.Spec.Template has no ImagePullSecrets", func() {
								BeforeEach(func() {
									jobDef.Spec.Template.ImagePullSecrets = nil
								})

								It("mutates", func() {
									Expect(len(err)).Should(Equal(0))
								})

								It("JobDefinitionSpec should not contain an ImagePullSecret", func() {
									Expect(jobRun.Spec.JobDefinitionSpec.Template.ImagePullSecrets).To(BeNil())
								})
							})
						})

						Context("with a JobRun that specified jobRun.Spec.JobDefinitionSpec.Template.ImagePullSecrets", func() {
							BeforeEach(func() {
								jobRun.Spec.JobDefinitionSpec.Template.ImagePullSecrets = []corev1.LocalObjectReference{
									{
										Name: "JobRunImagePullSecret",
									},
								}
							})

							It("does not mutate", func() {
								Expect(len(err)).ToNot(Equal(0))
							})

							It("fails with meaningful error message", func() {
								Expect(errors.AggregateErrorList(err)).To(ContainSubstring("spec.jobDefinitionSpec.template.imagePullSecrets: Unsupported field: disallow to set with referenced jobDefinition"))
							})
						})
					})

					Describe("ServiceAccountName", func() {
						Context("with a JobRun that did not specify jobRun.Spec.JobDefinitionSpec.Template.ServiceAccountName", func() {
							BeforeEach(func() {
								jobRun.Spec.JobDefinitionSpec.Template.ServiceAccountName = ""
							})

							Context("where jobDef.Spec.Template has ServiceAccountName", func() {
								BeforeEach(func() {
									jobDef.Spec.Template.ServiceAccountName = "JobDefServiceAccountName"
								})

								It("mutates", func() {
									Expect(len(err)).Should(Equal(0))
								})

								It("JobDefinitionSpec contains given ServiceAccountName from jobDef.Spec.Template.ServiceAccountName", func() {
									Expect(jobRun.Spec.JobDefinitionSpec.Template.ServiceAccountName).To(Equal("JobDefServiceAccountName"))
								})
							})

							Context("where jobDef.Spec.Template has no ServiceAccountName", func() {
								BeforeEach(func() {
									jobDef.Spec.Template.ServiceAccountName = ""
								})

								It("mutates", func() {
									Expect(len(err)).Should(Equal(0))
								})

								It("JobDefinitionSpec should have an empty ServiceAccountName", func() {
									Expect(jobRun.Spec.JobDefinitionSpec.Template.ServiceAccountName).To(Equal(""))
								})
							})
						})

						Context("with a JobRun that specified jobRun.Spec.JobDefinitionSpec.Template.ServiceAccountName", func() {
							BeforeEach(func() {
								jobRun.Spec.JobDefinitionSpec.Template.ServiceAccountName = "JobRunServiceAccount"
							})

							It("does not mutate", func() {
								Expect(len(err)).NotTo(Equal(0))
							})

							It("fails with meaningful error message", func() {
								Expect(errors.AggregateErrorList(err)).To(ContainSubstring("spec.jobDefinitionSpec.template.serviceAccountName: Unsupported field: disallow to set with referenced jobDefinition"))
							})
						})
					})

					Describe("Containers", func() {
						// Currently we only allow one container and therefore we always pick the first containers.
						// Grep for d31add5e to see other references to this constraint

						Context("with a JobRun that did not specify containers", func() {
							//jobDef.Spec.Template will always have a Container because of our webhook validation.
							//Therefore we will fill it here explicitly with some value
							BeforeEach(func() {
								jobDef.Spec.Template.Containers = []corev1.Container{
									{
										Name:  "Foo",
										Image: "Bar",
										Resources: corev1.ResourceRequirements{
											Requests: corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse("1"),
												corev1.ResourceMemory: resource.MustParse("64Mi"),
											},
										},
									},
								}
								jobRun.Spec.JobDefinitionSpec.Template.Containers = nil
							})

							It("mutates", func() {
								Expect(len(err)).Should(Equal(0))
							})

							It("JobDefinitionSpec contains given container from jobDef.Spec.Template", func() {
								Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers).To(Equal([]corev1.Container{
									{
										Name:  "Foo",
										Image: "Bar",
										Resources: corev1.ResourceRequirements{
											Requests: corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse("1"),
												corev1.ResourceMemory: resource.MustParse("64Mi"),
											},
										},
									},
								}))
							})
						})

						Context("with a JobRun that specifies one container", func() {
							//jobDef.Spec.Template will always have a Container because of our webhook validation.
							//Therefore we will fill it here explicitly with some value
							BeforeEach(func() {
								//this would be a minimal version of a container that a jobdefiniton will have
								jobDef.Spec.Template.Containers = []corev1.Container{
									{
										Name:  "ContainerName",
										Image: "Image",
										Resources: corev1.ResourceRequirements{
											Requests: corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse("1"),
												corev1.ResourceMemory: resource.MustParse("64Mi"),
											},
										},
									},
								}

								jobRun.Spec.JobDefinitionSpec.Template.Containers = []corev1.Container{
									{
										Name:  "ContainerName",
										Image: "Image",
										Resources: corev1.ResourceRequirements{
											Requests: corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse("2"),
												corev1.ResourceMemory: resource.MustParse("128Mi"),
											},
										},
									},
								}
							})

							Describe("Image", func() {
								//The container of the referenced jobdef will always have a image
								//Therefore we will fill it here explicitly with some value
								BeforeEach(func() {
									jobDef.Spec.Template.Containers[0].Image = "Some-Image:JobDef"
								})

								Context("where the container did not specify Image", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Image = ""
									})

									It("mutates", func() {
										Expect(len(err)).Should(Equal(0))
									})

									It("JobDefinitionSpec contains the Image from referenced JobDefiniton", func() {
										Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Image).To(Equal("Some-Image:JobDef"))
									})
								})

								Context("where the container specified Image", func() {
									Context("and the container image is the same as the image from the referenced JobDefinition", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Image = "Some-Image:Same"
											jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Image = "Some-Image:Same"
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the image", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Image).To(Equal("Some-Image:Same"))
										})
									})
									Context("and the container image is not the same as the image from the referenced JobDefinition", func() {
										BeforeEach(func() {
											jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Image = "Some-Image:JobRun"
										})

										It("does not mutate", func() {
											Expect(len(err)).NotTo(Equal(0))
										})

										It("fails with meanigful error message", func() {
											Expect(errors.AggregateErrorList(err)).To(ContainSubstring("spec.jobDefinitionSpec.template.containers[0].image: Invalid value: \"Some-Image:JobDef\": must be the same as referenced jobDefinition.spec.template.containers[0].image"))
										})
									})
								})
							})

							Describe("Name", func() {
								BeforeEach(func() {
									//The container of the referenced jobdef will always have a name
									//Therefore we will fill it here explicitly with some value
									jobDef.Spec.Template.Containers[0].Name = "Some-Name-From-JobDef"
								})

								Context("where the jobrun container did not specify Name", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Name = ""
									})

									It("mutates", func() {
										Expect(len(err)).Should(Equal(0))
									})

									It("JobDefinitionSpec contains the name from referenced JobDefiniton", func() {
										Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Name).To(Equal("Some-Name-From-JobDef"))
									})
								})

								Context("where the JobRun's container name is not empty", func() {
									Context("and the container name is the same as the name from the referenced JobDefinition", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Name = "Same-Name"
											jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Name = "Same-Name"
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the name", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Name).To(Equal("Same-Name"))
										})
									})
									Context("and the container name is not the same as the name from the referenced JobDefinition", func() {
										BeforeEach(func() {
											jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Name = "Some-Name-From-JobRun"
										})

										It("does not mutate", func() {
											Expect(len(err)).NotTo(Equal(0))
										})

										It("fails with meanigful error message", func() {
											Expect(errors.AggregateErrorList(err)).To(ContainSubstring("spec.jobDefinitionSpec.template.containers[0].name: Invalid value: \"Some-Name-From-JobDef\": must be the same as referenced jobDefinition.spec.template.containers[0].name"))
										})
									})
								})
							})

							Describe("Command", func() {
								Context("where the container specified Command", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Command = []string{"JobRunCommand"}
									})

									Context("where the referenced JobDefinition.Container has a Command", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Command = []string{"JobDefinitionCommand"}
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the Command from JobRun", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Command).To(Equal([]string{"JobRunCommand"}))
										})
									})
									Context("where the referenced JobDefinition specifies no Commands", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Command = nil
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the Command from JobRun", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Command).To(Equal([]string{"JobRunCommand"}))
										})
									})
								})

								Context("where the container did not specify Command", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Command = nil
									})

									Context("where the referenced JobDefinition specifies Commands", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Command = []string{"JobDefinitionCommand"}
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the Command from referenced JobDefinition", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Command).To(Equal([]string{"JobDefinitionCommand"}))
										})
									})

									Context("where the referenced JobDefinition specifies no Commands", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Command = nil
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains no Command at all", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Command).To(BeNil())
										})
									})
								})
							})

							Describe("Args", func() {
								Context("where the container specified Args", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Args = []string{"JobRunArg1", "JobRunArg2", "JobRunArg3"}
									})

									Context("where the referenced JobDefinition specified Args", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Args = []string{"JobDefinitionArgs1"}
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the container.Args from JobRun", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Args).To(Equal([]string{"JobRunArg1", "JobRunArg2", "JobRunArg3"}))
										})
									})

									Context("where the referenced JobDefinition specifies no Args", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Args = nil
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the JobRun.container.Args", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Args).To(Equal([]string{"JobRunArg1", "JobRunArg2", "JobRunArg3"}))
										})
									})
								})

								Context("where the container did not specify Args", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Args = nil
									})

									Context("where the referenced JobDefinition specifies Args", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Args = []string{"JobDefArg1", "JobDefArg2", "JobDefArg3"}
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the container.Args from referenced JobDefinition", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Args).To(Equal([]string{"JobDefArg1", "JobDefArg2", "JobDefArg3"}))
										})
									})

									Context("where the referenced JobDefinition specifies no Args", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Args = nil
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains no container.Args at all", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Args).To(BeNil())
										})
									})
								})
							})

							Describe("EnvFrom", func() {
								Context("where the container specified EnvFrom", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].EnvFrom = []corev1.EnvFromSource{
											{
												ConfigMapRef: &corev1.ConfigMapEnvSource{
													LocalObjectReference: corev1.LocalObjectReference{
														Name: "JobRunContainerEnvFrom",
													},
												},
											},
										}
									})

									Context("where the referenced JobDefinition specifies EnvFrom", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].EnvFrom = []corev1.EnvFromSource{
												{
													ConfigMapRef: &corev1.ConfigMapEnvSource{
														LocalObjectReference: corev1.LocalObjectReference{
															Name: "JobDefContainerEnvFrom",
														},
													},
												},
											}
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the container.EnvFrom from JobRun", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].EnvFrom).To(Equal(([]corev1.EnvFromSource{
												{
													ConfigMapRef: &corev1.ConfigMapEnvSource{
														LocalObjectReference: corev1.LocalObjectReference{
															Name: "JobRunContainerEnvFrom",
														},
													},
												},
											})))
										})
									})

									Context("where the referenced JobDefinition specifies no EnvFrom", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].EnvFrom = nil
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the container.EnvFrom from JobRun", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].EnvFrom).To(Equal(([]corev1.EnvFromSource{
												{
													ConfigMapRef: &corev1.ConfigMapEnvSource{
														LocalObjectReference: corev1.LocalObjectReference{
															Name: "JobRunContainerEnvFrom",
														},
													},
												},
											})))
										})
									})
								})

								Context("where the container did not specify EnvFrom", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].EnvFrom = nil
									})

									Context("where the referenced JobDefinition specifies EnvFrom", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].EnvFrom = []corev1.EnvFromSource{
												{
													ConfigMapRef: &corev1.ConfigMapEnvSource{
														LocalObjectReference: corev1.LocalObjectReference{
															Name: "JobDefContainerEnvFrom",
														},
													},
												},
											}
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the container.EnvFrom from referenced JobDefinition", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].EnvFrom).To(Equal(([]corev1.EnvFromSource{
												{
													ConfigMapRef: &corev1.ConfigMapEnvSource{
														LocalObjectReference: corev1.LocalObjectReference{
															Name: "JobDefContainerEnvFrom",
														},
													},
												},
											})))
										})
									})

									Context("where the referenced JobDefinition specifies no EnvFrom", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].EnvFrom = nil
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains no container.EnvFrom at all", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].EnvFrom).To(BeNil())
										})
									})
								})
							})

							Describe("EnvVars", func() {
								Context("where the container has one EnvVar", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env = []corev1.EnvVar{
											{
												Name:  "JobRunEnv1Name",
												Value: "JobRunEnv1Value",
											},
										}
									})

									Context("where the container of the referenced JobDefinition has one EnvVar", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Env = []corev1.EnvVar{
												{
													Name:  "JobDefEnv1Name",
													Value: "JobDefEnv1Value",
												},
											}
										})

										Context("and both specified container.Env[*].Name fields are different", func() {
											It("mutates", func() {
												Expect(len(err)).Should(Equal(0))
											})

											It("JobDefinitionSpec has both env var entries", func() {
												Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env).To(Equal([]corev1.EnvVar{
													{
														Name:  "JobDefEnv1Name",
														Value: "JobDefEnv1Value",
													},
													{
														Name:  "JobRunEnv1Name",
														Value: "JobRunEnv1Value",
													},
												}))
											})
										})

										Context("where the jobrun container and jobdefinition container has an env var with same name", func() {
											Context("where the value of these env vars are different", func() {
												BeforeEach(func() {
													jobDef.Spec.Template.Containers[0].Env = []corev1.EnvVar{
														{
															Name:  "ENV_NAME_DUPLICATE",
															Value: "JobDefinition",
														},
													}
													jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env = []corev1.EnvVar{
														{
															Name:  "ENV_NAME_DUPLICATE",
															Value: "JobRun",
														},
													}
												})

												It("mutates", func() {
													Expect(len(err)).Should(Equal(0))
												})

												Specify("JobDefinitionSpec has env var value of the jobRun", func() {
													Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env).To(Equal([]corev1.EnvVar{
														{
															Name:  "ENV_NAME_DUPLICATE",
															Value: "JobRun",
														},
													}))
												})
											})

											Context("where the value of the jobrun env is empty", func() {
												BeforeEach(func() {
													jobDef.Spec.Template.Containers[0].Env = []corev1.EnvVar{
														{
															Name:  "ENV_NAME_DUPLICATE",
															Value: "JobDefEnv1Value",
														},
													}
													jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env = []corev1.EnvVar{
														{
															Name:  "ENV_NAME_DUPLICATE",
															Value: "",
														},
													}
												})

												It("mutates", func() {
													Expect(len(err)).Should(Equal(0))
												})

												Specify("this env var entry does not exist in the JobDefinitionSpec", func() {
													Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env).To(HaveLen(0))
												})
											})

											Context("where the value of both is empty", func() {
												BeforeEach(func() {
													jobDef.Spec.Template.Containers[0].Env = []corev1.EnvVar{
														{
															Name:  "ENV_NAME_DUPLICATE",
															Value: "",
														},
													}
													jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env = []corev1.EnvVar{
														{
															Name:  "ENV_NAME_DUPLICATE",
															Value: "",
														},
													}
												})

												It("mutates", func() {
													Expect(len(err)).Should(Equal(0))
												})

												Specify("this env var entry does not exist in the JobDefinitionSpec", func() {
													Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env).To(HaveLen(0))
												})
											})

										})
									})

									Context("where the container of the referenced JobDefinition specifies no Env", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Env = nil
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										Specify("the EnvVar of JobDefinitionSpec contains the EnvVar of the JobRun", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env).To(Equal([]corev1.EnvVar{
												{
													Name:  "JobRunEnv1Name",
													Value: "JobRunEnv1Value",
												},
											}))
										})
									})
								})

								Context("where the container has no EnvVars", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env = nil
									})

									Context("where the referenced JobDefinition specifies EnvVars", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Env = []corev1.EnvVar{
												{
													Name:  "JobDefEnv1Name",
													Value: "JobDefEnv1Value",
												},
												{
													Name:  "JobDefEnv2Name",
													Value: "JobDefEnv2Value",
												},
											}
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										Specify("jobRun.Spec.JobDefinitionSpec.Template.Containers.Env contains jobDef.Spec.Template.Containers[0].Env", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env).To(Equal([]corev1.EnvVar{
												{
													Name:  "JobDefEnv1Name",
													Value: "JobDefEnv1Value",
												},
												{
													Name:  "JobDefEnv2Name",
													Value: "JobDefEnv2Value",
												},
											}))
										})
									})

									Context("where the referenced JobDefinition specifies no EnvVars", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Env = nil
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										Specify("jobRun.Spec.JobDefinitionSpec.Template.Containers.Env is nil", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env).To(BeNil())
										})
									})
								})

								Context("where the container has an EnvVar with empty value", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env = []corev1.EnvVar{
											{
												Name: "JobRunEnv1Name",
											},
										}
									})

									Context("and the referenced JobDefinition doesn't not have it", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Env = []corev1.EnvVar{}
										})

										It("JobDefinitionSpec doesn't have the env", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env).To(HaveLen(0))
										})
									})
								})

								Context("where the referenced JobDefinition container has a envvar with empty value", func() {
									BeforeEach(func() {
										jobDef.Spec.Template.Containers[0].Env = []corev1.EnvVar{
											{
												Name: "JobDefEnv1Name",
											},
										}
									})

									Context("and the jobrun does not have it", func() {
										BeforeEach(func() {
											jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env = []corev1.EnvVar{}
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec has the env var with an empty value", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env).To(Equal([]corev1.EnvVar{
												{
													Name: "JobDefEnv1Name",
												},
											}))
										})
									})
								})

								Context("where the jobrun container has a EnvVar with empty value", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env = []corev1.EnvVar{
											{
												Name: "JobRunEnv1Name",
											},
										}
									})

									Context("and the referenced JobDefinition does not have it", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Env = []corev1.EnvVar{}
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec doesn't have the env", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env).To(HaveLen(0))
										})
									})
								})

								Context("Multiple environment variables", func() {
									BeforeEach(func() {
										jobDef.Spec.Template.Containers[0].Env = []corev1.EnvVar{
											{Name: "JobDefEnv1Name", Value: "JobDefEnv1Value"},
											{Name: "JobDefEnv2Name", Value: "JobDefEnv2Value"},
											{Name: "EMPTY"},
											{Name: "DUPLICATE0", Value: "JobDef"},
											{Name: "DUPLICATE1", Value: "sadfads"},
										}

										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env = []corev1.EnvVar{
											{Name: "JobRunEnv1Name", Value: "JobRunEnv1Value"},
											{Name: "JobRunEnv2Name", Value: "JobRunEnv2Value"},
											{Name: "EmptyVar", Value: ""},
											{Name: "DUPLICATE0", Value: "JobRun"},
											{Name: "DUPLICATE1", Value: ""},
										}
									})

									It("mutates", func() {
										Expect(len(err)).Should(Equal(0))
									})

									It("merges together correctly", func() {
										expectedEnvVars := []corev1.EnvVar{
											{Name: "JobDefEnv1Name", Value: "JobDefEnv1Value"}, // from JobDef
											{Name: "JobDefEnv2Name", Value: "JobDefEnv2Value"}, // from JobDef
											{Name: "EMPTY", Value: ""},                         // from JobDef; "" is semantically the same as the zero value
											{Name: "DUPLICATE0", Value: "JobRun"},              // JobRun wins
											{Name: "JobRunEnv1Name", Value: "JobRunEnv1Value"}, // from JobRun
											{Name: "JobRunEnv2Name", Value: "JobRunEnv2Value"}, // from JobRun

											// DUPLICATE1 does not appear because being present in the JobRun with an
											// empty variable erases DUPLICATE1 from the JobDef

											// EmptyVar does not appear because it is empty, and the presence of a JobDef
											// changes the merge behavior for empty variables to be an override. In this case,
											// it means that EmptyVar simply does not appear in the effective JobDefinitionSpec.
										}

										// order of env vars must not matter
										sort.Sort(ByName(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env))
										sort.Sort(ByName(expectedEnvVars))

										Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Env).To(BeEquivalentTo(expectedEnvVars))
									})
								})

							})

							Describe("Resources.Requests", func() {
								Context("where the container specified Resources.Requests", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
											corev1.ResourceCPU:    resource.MustParse("2"),
											corev1.ResourceMemory: resource.MustParse("128Mi"),
										}
									})

									Context("where the referenced JobDefinition specifies Resources.Requests", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse("4"),
												corev1.ResourceMemory: resource.MustParse("264Mi"),
											}
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the container.Resources.Requests from JobRun", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests).To(Equal((corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse("2"),
												corev1.ResourceMemory: resource.MustParse("128Mi"),
											})))
										})
									})

									Context("where the referenced JobDefinition specifies no Resources.Requests", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Resources.Requests = nil
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the container.Resources.Requests from JobRun", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests).To(Equal((corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse("2"),
												corev1.ResourceMemory: resource.MustParse("128Mi"),
											})))
										})
									})
								})

								Context("where the container did not specify Resources.Requests", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests = nil
									})

									Context("where the referenced JobDefinition specifies Resources.Requests", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse("2"),
												corev1.ResourceMemory: resource.MustParse("128Mi"),
											}
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the container.Resources.Requests from referenced JobDefinition", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests).To(Equal((corev1.ResourceList{
												corev1.ResourceCPU:    resource.MustParse("2"),
												corev1.ResourceMemory: resource.MustParse("128Mi"),
											})))
										})
									})

									Context("where the referenced JobDefinition specifies no Resources.Requests", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].Resources.Requests = nil
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains no container.Resources.Requests at all", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests).To(BeNil())
										})
									})
								})
							})

							Describe("TerminationMessagePath", func() {
								Context("where the container specified TerminationMessagePath", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].TerminationMessagePath = "JobRunTerminationMessagePath"
									})

									It("does not mutate", func() {
										Expect(len(err)).NotTo(Equal(0))
									})

									It("fails with meanigful error message", func() {
										Expect(errors.AggregateErrorList(err)).To(ContainSubstring("spec.jobDefinitionSpec.template.containers[0].terminationMessagePath: Unsupported field: disallow to set with referenced jobDefinition"))
									})
								})

								Context("where the container did not specify TerminationMessagePath", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].TerminationMessagePath = ""
									})

									Context("where the referenced JobDefinition specifies TerminationMessagePath", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].TerminationMessagePath = "JobDefinitionTerminationMessagePath"
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the TerminationMessagePath from referenced JobDefinition", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].TerminationMessagePath).To(Equal("JobDefinitionTerminationMessagePath"))
										})
									})

									Context("where the referenced JobDefinition specifies no TerminationMessagePath", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].TerminationMessagePath = ""
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains no TerminationMessagePath at all", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].TerminationMessagePath).To(Equal(""))
										})
									})
								})
							})

							Describe("TerminationMessagePolicy", func() {
								Context("where the container specified TerminationMessagePolicy", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].TerminationMessagePolicy = corev1.TerminationMessageReadFile
									})

									It("does not mutate", func() {
										Expect(len(err)).NotTo(Equal(0))
									})

									It("fails with meaningful error message", func() {
										Expect(errors.AggregateErrorList(err)).To(ContainSubstring("spec.jobDefinitionSpec.template.containers[0].terminationMessagePolicy: Unsupported field: disallow to set with referenced jobDefinition"))
									})
								})

								Context("where the container did not specify TerminationMessagePolicy", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].TerminationMessagePolicy = ""
									})

									Context("where the referenced JobDefinition specifies TerminationMessagePolicy", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].TerminationMessagePolicy = corev1.TerminationMessageFallbackToLogsOnError
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the TerminationMessagePolicy from referenced JobDefinition", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].TerminationMessagePolicy).To(Equal(corev1.TerminationMessageFallbackToLogsOnError))
										})
									})

									Context("where the referenced JobDefinition specifies no TerminationMessagePolicy", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].TerminationMessagePolicy = ""
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains no TerminationMessagePolicy at all", func() {
											Expect(string(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].TerminationMessagePolicy)).To(Equal(""))
										})
									})
								})
							})

							Describe("WorkingDir", func() {
								Context("where the container specified WorkingDir", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].WorkingDir = "JobRunWorkingDir"
									})

									It("does not mutate", func() {
										Expect(len(err)).NotTo(Equal(0))
									})

									It("fails with meanigful error message", func() {
										Expect(errors.AggregateErrorList(err)).To(ContainSubstring("spec.jobDefinitionSpec.template.containers[0].workingDir: Unsupported field: disallow to set with referenced jobDefinition"))
									})
								})

								Context("where the container did not specify WorkingDir", func() {
									BeforeEach(func() {
										jobRun.Spec.JobDefinitionSpec.Template.Containers[0].WorkingDir = ""
									})

									Context("where the referenced JobDefinition specifies WorkingDir", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].WorkingDir = "JobDefinitionWorkingDir"
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains the WorkingDir from referenced JobDefinition", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].WorkingDir).To(Equal("JobDefinitionWorkingDir"))
										})
									})

									Context("where the referenced JobDefinition specifies no WorkingDir", func() {
										BeforeEach(func() {
											jobDef.Spec.Template.Containers[0].WorkingDir = ""
										})

										It("mutates", func() {
											Expect(len(err)).Should(Equal(0))
										})

										It("JobDefinitionSpec contains no WorkingDir at all", func() {
											Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].WorkingDir).To(Equal(""))
										})
									})
								})
							})
						})
					})
				})
			})
		})
	})

	Describe("MutateResourcesRequestsWithLimitRange()", func() {
		var (
			limitRangeItem       *corev1.LimitRangeItem
			err                  field.ErrorList
			jobRun               *cv1b1.JobRun
			jobRunBeforeMutation *cv1b1.JobRun
		)

		BeforeEach(func() {
			jobRun = &cv1b1.JobRun{}
		})

		JustBeforeEach(func() {
			jobRunBeforeMutation = &cv1b1.JobRun{}
			jobRun.DeepCopyInto(jobRunBeforeMutation)
			err = jobRun.MutateResourcesRequestsWithLimitRange(limitRangeItem)
		})

		When("the JobRun doesn't specify Template.Containers", func() {
			BeforeEach(func() {
				jobRun.Spec.JobDefinitionSpec.Template.Containers = nil
			})

			It("returns an error", func() {
				Expect(errors.AggregateErrorList(err)).To(SatisfyAll(
					ContainSubstring("there must be exactly one container"),
					ContainSubstring("jobDefinitionSpec.template.containers"),
				))
			})
		})

		When("the JobRun does specify Template.Containers", func() {
			BeforeEach(func() {
				jobRun.Spec.JobDefinitionSpec.Template.Containers = []corev1.Container{
					{
						Name: "fake-container",
					},
				}
			})

			When("there is no LimitRangeItem", func() {
				BeforeEach(func() {
					limitRangeItem = nil
				})

				It("does not throw an error", func() {
					Expect(len(err)).Should(Equal(0))
				})

				It("does not mutate the JobRun", func() {
					Expect(jobRun).To(Equal(jobRunBeforeMutation))
				})
			})

			When("there is a LimitRangeItem with DefaultRequests", func() {
				BeforeEach(func() {
					limitRangeItem = &corev1.LimitRangeItem{
						Type: corev1.LimitTypeContainer,
						DefaultRequest: corev1.ResourceList{
							corev1.ResourceCPU:              resource.MustParse("100m"),
							corev1.ResourceMemory:           resource.MustParse("32Mi"),
							corev1.ResourceEphemeralStorage: resource.MustParse("128Mi"),
						},
					}
				})

				Describe("CPU", func() {
					Context("and the JobRun specifies a CPU value", func() {
						BeforeEach(func() {
							jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceCPU: resource.MustParse("1"),
							}
						})

						It("the specified value in JobRun is used", func() {
							Expect(*jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests.Cpu()).To(
								Equal(resource.MustParse("1")),
							)
						})
					})

					Context("and the JobRun does not specify a CPU value", func() {
						BeforeEach(func() {
							jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests = corev1.ResourceList{}
						})

						It("the CPU value equals the DefaultRequest of given LimitRangeItem", func() {
							Expect(*jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests.Cpu()).To(
								Equal(resource.MustParse("100m")),
							)
						})
					})
				})

				Describe("Memory", func() {
					Context("and the JobRun specifies a memory value", func() {
						BeforeEach(func() {
							jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("100Mi"),
							}
						})

						It("the specified value in JobRun is used", func() {
							Expect(*jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests.Memory()).To(
								Equal(resource.MustParse("100Mi")),
							)
						})
					})

					Context("and the JobRun does not specify a memory value", func() {
						BeforeEach(func() {
							jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests = corev1.ResourceList{}
						})

						It("the CPU value equals the DefaultRequest of given LimitRangeItem", func() {
							Expect(*jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests.Memory()).To(
								Equal(resource.MustParse("32Mi")),
							)
						})
					})
				})

				Describe("EphemeralStorage", func() {
					Context("and the JobRun specifies a ephemeral-storage value", func() {
						BeforeEach(func() {
							jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests = corev1.ResourceList{
								corev1.ResourceEphemeralStorage: resource.MustParse("256Mi"),
							}
						})

						It("the specified value in JobRun is used", func() {
							Expect(*jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests.StorageEphemeral()).To(
								Equal(resource.MustParse("256Mi")),
							)
						})
					})

					Context("and the JobRun does not specify a ephemeral-storage value", func() {
						BeforeEach(func() {
							jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests = corev1.ResourceList{}
						})

						It("the CPU value equals the DefaultRequest of given LimitRangeItem", func() {
							Expect(*jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests.StorageEphemeral()).To(
								Equal(resource.MustParse("128Mi")),
							)
						})
					})
				})
			})
		})
	})

	Describe("MutateResourcesRequestsWithStaticLimits()", func() {
		var (
			err    field.ErrorList
			jobRun *cv1b1.JobRun
		)

		BeforeEach(func() {
			jobRun = &cv1b1.JobRun{
				ObjectMeta: metav1.ObjectMeta{
					Name: "fake-JobRun",
				},
				Spec: cv1b1.JobRunSpec{
					JobDefinitionSpec: cv1b1.JobDefinitionSpec{
						ArraySpec:  pointers.StringPtr("0"),
						RetryLimit: pointers.Int64Ptr(2),
						Template: cv1b1.JobPodTemplate{
							Containers: []corev1.Container{
								{
									Name:  "fake-container",
									Image: "fake-image",
									Resources: corev1.ResourceRequirements{
										Requests: corev1.ResourceList{
											corev1.ResourceCPU:              resource.MustParse("1"),
											corev1.ResourceMemory:           resource.MustParse("64Mi"),
											corev1.ResourceEphemeralStorage: resource.MustParse("256Mi"),
										},
									},
								},
							},
						},
					},
				},
			}
		})

		JustBeforeEach(func() {
			err = jobRun.MutateResourcesRequestsWithStaticLimits(ctx)
		})

		When("the JobRun doesn't specify a single element in Template.Containers", func() {
			BeforeEach(func() {
				jobRun.Spec.JobDefinitionSpec.Template.Containers = []corev1.Container{}
			})

			It("returns an error", func() {
				Expect(errors.AggregateErrorList(err)).To(SatisfyAll(
					ContainSubstring("there must be exactly one container"),
					ContainSubstring("jobDefinitionSpec.template.containers"),
				))
			})
		})

		When("the JobRun specifies a single Template.Container", func() {
			JustBeforeEach(func() {
				Expect(len(err)).Should(Equal(0))
			})

			When("the JobRun has a complete set of Resource.Requests", func() {
				It("passes", func() {
					Expect(len(err)).Should(Equal(0))
				})
			})

			When("the JobRun doesn't specify Resource.Requests", func() {
				BeforeEach(func() {
					jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources = corev1.ResourceRequirements{}
				})

				It("creates the missing structures and all mandatory Resource.Requests", func() {
					Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests).To(
						SatisfyAll(
							HaveLen(3),
							HaveKeyWithValue(corev1.ResourceCPU, resource.MustParse("1")),
							HaveKeyWithValue(corev1.ResourceMemory, resource.MustParse("4G")),
							HaveKeyWithValue(corev1.ResourceEphemeralStorage, resource.MustParse("4G")),
						))
				})
			})

			When("the JobRun defines partial Resources.Requests", func() {
				BeforeEach(func() {
					jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources = corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("64Mi"),
						},
					}
				})

				It("creates the missing elements and does not modify the defined ones", func() {
					Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests).To(
						SatisfyAll(
							HaveLen(3),
							HaveKeyWithValue(corev1.ResourceMemory, resource.MustParse("64Mi")),
						), "verify that predefined request.resources remain unchanged")

					Expect(jobRun.Spec.JobDefinitionSpec.Template.Containers[0].Resources.Requests).To(
						SatisfyAll(
							HaveLen(3),
							HaveKeyWithValue(corev1.ResourceCPU, resource.MustParse("1")),
							HaveKeyWithValue(corev1.ResourceEphemeralStorage, resource.MustParse("4G")),
						), "verify that non defined request.resources are extended with static resource values")
				})
			})
		})
	})
})
