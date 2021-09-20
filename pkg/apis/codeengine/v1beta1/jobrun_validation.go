/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.ibm.com/coligo/batch-job-controller/pkg/utils/errors"
)

// jobRunNameMaxLen is the maximum length of a jobRun name.
//
// As of Kubernetes, the max limit on a name is 63 chars. We reserve 10 for
// controller to add data{arraySpec, retryLimit}.
const jobRunNameMaxLen = 53

// Validate ensures JobRun is properly configured.
func (jr *JobRun) Validate(validationConfig ValidationConfig) field.ErrorList {
	errs := field.ErrorList{}

	if len(jr.Name) > jobRunNameMaxLen {
		errs = append(errs, errors.Invalid(field.NewPath("metadata", "name"), fmt.Sprintf("%q", jr.Name), fmt.Sprintf("name exceeded max length of %d", jobRunNameMaxLen)))
	}

	return append(errs, ValidateJobDefinitionSpec(&jr.Spec.JobDefinitionSpec, validationConfig, field.NewPath("spec", "jobDefinitionSpec"))...)
}
