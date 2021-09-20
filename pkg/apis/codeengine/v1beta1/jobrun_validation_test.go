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
	"k8s.io/apimachinery/pkg/util/validation/field"

	cv1b1 "github.ibm.com/coligo/batch-job-controller/pkg/apis/codeengine/v1beta1"
	"github.ibm.com/coligo/batch-job-controller/pkg/utils/errors"
	"github.ibm.com/coligo/batch-job-controller/pkg/utils/pointers"
)

var _ = Describe("JobRun Validation", func() {

	When("validating a JobRun", func() {
		var (
			jobRun *cv1b1.JobRun
			err    field.ErrorList
		)

		BeforeEach(func() {
			jobRun = &cv1b1.JobRun{}
		})

		JustBeforeEach(func() {
			err = jobRun.Validate(cv1b1.ValidationConfig{})
		})

		Context("that is valid", func() {
			BeforeEach(func() {
				jobRun = &cv1b1.JobRun{
					Spec: cv1b1.JobRunSpec{
						JobDefinitionSpec: cv1b1.JobDefinitionSpec{
							ArraySpec:  pointers.StringPtr("0"),
							RetryLimit: pointers.Int64Ptr(2),
							Template: cv1b1.JobPodTemplate{
								Containers: []corev1.Container{
									{
										Name:  "container-name",
										Image: "image-name",
									},
								},
							},
						},
					},
				}
			})

			It("does not error", func() {
				Expect(len(err)).Should(Equal(0))
			})
		})

		Context("which name exceeds the maximum length", func() {
			BeforeEach(func() {
				jobRun.Name = "jobrun-invalid-long-name-aaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			})

			It("a meaningful error message is returned", func() {
				Expect(errors.AggregateErrorList(err)).To(
					ContainSubstring("name: Invalid value: \"jobrun-invalid-long-name-aaaaaaaaaaaaaaaaaaaaaaaaaaaaa\": name exceeded max length of 53"),
				)
			})
		})
	})

	When("validating JobRun.Spec", func() {
		var (
			jobRunSpec *cv1b1.JobRunSpec
			err        field.ErrorList
		)

		BeforeEach(func() {
			jobRunSpec = &cv1b1.JobRunSpec{
				JobDefinitionSpec: cv1b1.JobDefinitionSpec{
					ArraySpec:  pointers.StringPtr("0"),
					RetryLimit: pointers.Int64Ptr(2),
					Template: cv1b1.JobPodTemplate{
						ServiceAccountName: "fake-service-account-name",
						Containers: []corev1.Container{
							{
								Name:  "container",
								Image: "image",
							},
						},
					},
				},
			}
		})

		JustBeforeEach(func() {
			err = cv1b1.ValidateJobDefinitionSpec(&jobRunSpec.JobDefinitionSpec, cv1b1.ValidationConfig{}, field.NewPath("jobDefinitionSpec"))
		})

		Context("that is valid", func() {
			It("does not error", func() {
				Expect(len(err)).Should(Equal(0))
			})
		})

		Context("that is invalid", func() {
			BeforeEach(func() {
				jobRunSpec.JobDefinitionSpec.Template.Containers = []corev1.Container{
					{
						Name: "container1",
					},
					{
						Name: "container2",
					},
				}
			})

			It("a meaningful error message gets returned", func() {
				Expect(errors.AggregateErrorList(err)).To(
					ContainSubstring("jobDefinitionSpec.template.containers: Unsupported container count: 2: there must be exactly one container"),
				)
			})
		})
	})
})
