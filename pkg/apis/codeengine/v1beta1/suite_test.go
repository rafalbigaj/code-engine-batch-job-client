/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1_test

import (
	"testing"

	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.ibm.com/coligo/batch-job-controller/pkg/ctxlog"
	"github.ibm.com/coligo/batch-job-controller/test"
)

var (
	ctx context.Context
)

var _ = BeforeSuite(func() {
	ctx = ctxlog.NewParentContext(test.CreateGinkgoLogger())
})

func TestTypes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Types Suite")
}
