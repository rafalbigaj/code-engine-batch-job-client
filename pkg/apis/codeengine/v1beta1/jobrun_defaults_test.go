/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cv1b1 "github.ibm.com/coligo/batch-job-controller/pkg/apis/codeengine/v1beta1"
	"github.ibm.com/coligo/batch-job-controller/pkg/utils/pointers"
)

var _ = Describe("JobRunDefaults", func() {
	Context("JobRun.Defaults", func() {
		var (
			jobRun *cv1b1.JobRun
		)

		BeforeEach(func() {
			jobRun = &cv1b1.JobRun{
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
											corev1.ResourceCPU:    *resource.NewQuantity(10, resource.DecimalSI),
											corev1.ResourceMemory: *resource.NewQuantity(128, resource.BinarySI),
										},
									},
								},
							},
						},
					},
				},
			}
		})

		It("sets a label for jobDefinition", func() {
			copyJobRun := jobRun.DeepCopy()
			copyJobRun.Spec.JobDefinitionRef = "fake-job-definition"
			copyJobRun.SetDefaults()
		})
	})

	Context("JobRunSpec.Defaults", func() {
		var (
			jobRunSpec *cv1b1.JobRunSpec
		)

		BeforeEach(func() {
			jobRunSpec = &cv1b1.JobRunSpec{
				JobDefinitionSpec: cv1b1.JobDefinitionSpec{
					Template: cv1b1.JobPodTemplate{
						Containers: []corev1.Container{
							{
								Name:  "fake-container",
								Image: "fake-image",
								Resources: corev1.ResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceCPU:    *resource.NewQuantity(10, resource.DecimalSI),
										corev1.ResourceMemory: *resource.NewQuantity(128, resource.BinarySI),
									},
								},
							},
						},
					},
				},
			}
		})

		It("sets defaults for jobRunSpec", func() {
			jobRunSpec.SetDefaults()
			By("JobRunSpec.ArraySpec is '0'")
			Expect(*jobRunSpec.JobDefinitionSpec.ArraySpec).Should(Equal("0"))
			By("JobRunSpec.RetryLimit is 3")
			Expect(*jobRunSpec.JobDefinitionSpec.RetryLimit).Should(Equal(int64(3)))
			By("JobRunSpec.MaxExecutionTime is set to 7200 seconds (2h)")
			Expect(*jobRunSpec.JobDefinitionSpec.MaxExecutionTime).Should(Equal(int64(7200)))
		})
	})

	Context("When calling Jobrun.SetDefaultsFromJobDefinition", func() {
		var (
			baseJobrun         *cv1b1.JobRun
			referredJobDefSpec cv1b1.JobDefinition
		)

		BeforeEach(func() {
			baseJobrun = &cv1b1.JobRun{
				Spec: cv1b1.JobRunSpec{
					JobDefinitionRef: "fake-job-definition",
					JobDefinitionSpec: cv1b1.JobDefinitionSpec{
						Template: cv1b1.JobPodTemplate{
							Containers: []corev1.Container{
								{
									Resources: corev1.ResourceRequirements{
										Requests: corev1.ResourceList{
											corev1.ResourceCPU:    *resource.NewQuantity(10, resource.DecimalSI),
											corev1.ResourceMemory: *resource.NewQuantity(128, resource.BinarySI),
										},
									},
								},
							},
						},
					},
				},
			}

			referredJobDefSpec = cv1b1.JobDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name: "fake-job-definition",
					UID:  "fake-uuid",
				},
				Spec: cv1b1.JobDefinitionSpec{
					Template: cv1b1.JobPodTemplate{
						Containers: []corev1.Container{
							{
								Name:  "fake-container",
								Image: "fake-image",
								Resources: corev1.ResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceCPU:    *resource.NewQuantity(10, resource.DecimalSI),
										corev1.ResourceMemory: *resource.NewQuantity(128, resource.BinarySI),
									},
								},
							},
						},
					},
				},
			}
		})

		JustBeforeEach(func() {
			cv1b1.SetDefaultsFromJobDefinition(baseJobrun, referredJobDefSpec)
		})
		Specify("Sets labels with the referenced JobDefinition", func() {
			Expect(baseJobrun.Labels).Should(SatisfyAll(
				HaveKeyWithValue(cv1b1.LabelJobDefName, "fake-job-definition"),
				HaveKeyWithValue(cv1b1.LabelJobDefUUID, "fake-uuid"),
			))
		})

		Specify("Sets the container name from the given JobDefinition container ", func() {
			Expect(baseJobrun.Spec.JobDefinitionSpec.Template.Containers[0].Name).Should(Equal("fake-container"))
		})

		Specify("Sets the container image from the given JobDefinition container", func() {
			Expect(baseJobrun.Spec.JobDefinitionSpec.Template.Containers[0].Image).Should(Equal("fake-image"))
		})
	})

	Context("JobRunSpec.RequiresDefaultingFromJobDefinition", func() {
		var (
			jobRunSpec *cv1b1.JobRunSpec
		)

		BeforeEach(func() {
			jobRunSpec = &cv1b1.JobRunSpec{
				JobDefinitionRef: "fake-jobdef-ref",
				JobDefinitionSpec: cv1b1.JobDefinitionSpec{
					Template: cv1b1.JobPodTemplate{
						Containers: []corev1.Container{
							{
								Resources: corev1.ResourceRequirements{
									Requests: corev1.ResourceList{
										corev1.ResourceCPU:    *resource.NewQuantity(10, resource.DecimalSI),
										corev1.ResourceMemory: *resource.NewQuantity(128, resource.BinarySI),
									},
								},
							},
						},
					},
				},
			}
		})

		It("check if the container needs to sets defaults in JobRun Spec.", func() {
			needDefaults := jobRunSpec.RequiresDefaultingFromJobDefinition()
			By("Container needs to set default")
			Expect(needDefaults).Should(BeTrue())

			jobRunSpec.JobDefinitionRef = ""
			needDefaults = jobRunSpec.RequiresDefaultingFromJobDefinition()
			By("Container doesn't to set default")
			Expect(needDefaults).Should(BeFalse())
		})
	})
})
