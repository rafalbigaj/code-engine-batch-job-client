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
	"github.ibm.com/coligo/batch-job-controller/pkg/utils/pointers"
)

var _ = Describe("JobDefinitionLifecycle", func() {
	var (
		jds *cv1b1.JobDefinitionSpec
	)
	Context("when CalculateArrayIndices", func() {
		var (
			Indices map[int64]interface{}
			err     error
		)

		JustBeforeEach(func() {
			Indices, err = jds.CalculateArrayIndices()
		})

		Context("JobDefinitionSpec.ArraySpec is set", func() {
			BeforeEach(func() {
				jds = &cv1b1.JobDefinitionSpec{
					ArraySpec: pointers.StringPtr("1,3,5,7-9"),
				}
			})

			It("does not error", func() {
				Expect(err).To(Not(HaveOccurred()))
			})

			It("it returns the correct count of JobRunIndices", func() {
				Expect(Indices).ToNot(BeNil())
				Expect(len(Indices)).Should(Equal(6))
			})
		})
		Context("JobDefinitionSpec.ArraySpec contains an invalid values", func() {
			BeforeEach(func() {
				jds = &cv1b1.JobDefinitionSpec{
					ArraySpec: pointers.StringPtr("-1"),
				}
			})

			It("throws an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
		Context("JobDefinitionSpec.ArraySpec specifies overlapping index ranges", func() {
			BeforeEach(func() {
				jds = &cv1b1.JobDefinitionSpec{
					ArraySpec: pointers.StringPtr("5-8,2-4,2-11"),
				}
			})

			It("does not error", func() {
				Expect(err).To(Not(HaveOccurred()))
			})

			It("return the correct count of JobIndices", func() {
				Expect(len(Indices)).Should(BeNumerically("==", 10))
			})
		})
		Context("JobDefinitionSpec.ArraySpec is invalid with overlapping ranges and invalid fields", func() {
			BeforeEach(func() {
				jds = &cv1b1.JobDefinitionSpec{
					ArraySpec: pointers.StringPtr("5-8,2-4,-11"),
				}
			})

			It("throws an error ", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
