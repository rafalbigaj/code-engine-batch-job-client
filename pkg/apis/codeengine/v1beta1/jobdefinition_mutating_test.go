/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cv1b1 "github.ibm.com/coligo/batch-job-controller/pkg/apis/codeengine/v1beta1"
)

var _ = Describe("Mutating JobDefinition", func() {
	Describe("JobDefinition Status", func() {
		var (
			jobDefinition *cv1b1.JobDefinition
		)

		BeforeEach(func() {
			jobDefinition = &cv1b1.JobDefinition{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "JobDef-Fake-Name",
					Namespace: "JobDef-Fake-Namespace",
				},
			}
		})

		JustBeforeEach(func() {
			jobDefinition.MutateJobDefStatus()
		})

		Describe("Address", func() {
			It("JobDefinition has the host address set correctly", func() {
				Expect(jobDefinition.Status.Address.URL.Host).To(Equal(fmt.Sprintf("%s.%s.svc.cluster.local", cv1b1.EventingJobrunnerName, cv1b1.EventingNamespace)))
			})
			It("JobDefinition has the path set correftly", func() {
				pathArray := strings.Split(jobDefinition.Status.Address.URL.Path, "/")
				pathWithoutUID := pathArray[1] + "/" + pathArray[2]
				Expect(pathWithoutUID).To(Equal("JobDef-Fake-Namespace/JobDef-Fake-Name"))
			})

		})

	})
})
