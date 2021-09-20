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
	cv1b1 "github.ibm.com/coligo/batch-job-controller/pkg/apis/codeengine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("jobrun_lifecycle", func() {
	Context("When calling JobRun.SetOwner", func() {

		var (
			jr                        *cv1b1.JobRun
			jd                        *cv1b1.JobDefinition
			jdKind, jdVersion, jdName string
			jdUID                     types.UID
		)

		BeforeEach(func() {
			jdKind = "JobDefinition"
			jdVersion = "codeengine.cloud.ibm.com/v1beta1"
			jdName = "JobDefName"
			jdUID = "JobDefUID"

			jr = &cv1b1.JobRun{}
			jd = &cv1b1.JobDefinition{}

			jd.Name = jdName
			jd.UID = jdUID
		})

		JustBeforeEach(func() {
			jr.SetOwner(jd)
		})

		Specify("JobRun has only one OwnerReference", func() {
			Expect(jr.GetOwnerReferences()).To(HaveLen(1))
		})

		Specify("OwnerReference has correct APIVersion", func() {
			Expect(jr.GetOwnerReferences()[0].APIVersion).To(Equal(jdVersion))
		})

		Specify("OwnerReference has correct Kind", func() {
			Expect(jr.GetOwnerReferences()[0].Kind).To(Equal(jdKind))
		})

		Specify("OwnerReference has JobDefinitionName as name", func() {
			Expect(jr.GetOwnerReferences()[0].Name).To(Equal(jdName))
		})

		Specify("OwnerReference has JobDefinitionUID as uid", func() {
			Expect(jr.GetOwnerReferences()[0].UID).To(Equal(jdUID))
		})

		Specify("OwnerReference.BlockOwnerDeletion is enabled (this ensures Jobruns gets deleted when the correspending JobDefinition gets deleted)", func() {
			Expect(jr.GetOwnerReferences()[0].BlockOwnerDeletion).ToNot(BeNil())
			Expect(*jr.GetOwnerReferences()[0].BlockOwnerDeletion).To(BeTrue())
		})

		Specify("OwnerReference.Controller is enabled (this ensures Jobruns gets deleted when the correspending JobDefinition gets deleted)", func() {
			Expect(jr.GetOwnerReferences()[0].Controller).ToNot(BeNil())
			Expect(*jr.GetOwnerReferences()[0].Controller).To(BeTrue())
		})

		Context("when there are other owner references", func() {
			BeforeEach(func() {
				jr.OwnerReferences = []metav1.OwnerReference{
					{
						Name: "something",
					},
				}
			})

			Specify("maintains other owner references", func() {
				Expect(jr.GetOwnerReferences()).To(HaveLen(2))
				Expect(jr.GetOwnerReferences()).To(ContainElement(metav1.OwnerReference{
					Name: "something",
				}))
			})
		})
	})

	Context("When calling JobRun.AddLabel", func() {

		var (
			jr         *cv1b1.JobRun
			key, value string
			overwrite  bool
		)

		BeforeEach(func() {
			jr = &cv1b1.JobRun{}
			key = "someKey"
			value = "someValue"
		})
		JustBeforeEach(func() {
			jr.AddLabel(key, value, overwrite)
		})
		Context("with overwrite enabled", func() {
			BeforeEach(func() {
				overwrite = true
			})
			Context("Label does not exists", func() {
				Specify("Given label will be added to the Jobrun", func() {
					Expect(jr.Labels).ToNot(BeNil())
					Expect(jr.Labels[key]).To(Equal(value))
				})
			})
			Context("Label already exists", func() {
				BeforeEach(func() {
					jr.Labels = map[string]string{
						key: "someFancyValue",
					}
				})
				Specify("Existing label will be overwritten by the new label", func() {
					Expect(jr.Labels[key]).ToNot(Equal("someFancyValue"))
					Expect(jr.Labels[key]).To(Equal(value))
				})
			})
		})
		Context("with overwrite disabled", func() {
			BeforeEach(func() {
				overwrite = false
			})
			Context("Label does not exists", func() {
				Specify("Given label will be added to the Jobrun", func() {
					Expect(jr.Labels).ToNot(BeNil())
					Expect(jr.Labels[key]).To(Equal(value))
				})
			})
			Context("Label already exists", func() {
				BeforeEach(func() {
					jr.Labels = map[string]string{
						key: "someFancyValue",
					}
				})
				Specify("Existing label will be overwritten by the new label", func() {
					Expect(jr.Labels[key]).To(Equal("someFancyValue"))
					Expect(jr.Labels[key]).ToNot(Equal(value))
				})
			})
		})
	})
})
