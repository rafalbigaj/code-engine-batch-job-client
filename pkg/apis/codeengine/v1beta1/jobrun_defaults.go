/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1

import (
	"github.ibm.com/coligo/batch-job-controller/pkg/utils/pointers"
)

// SetDefaults sets defaults for JobRun.
func (jr *JobRun) SetDefaults() {
	jr.Spec.SetDefaults()
}

// SetDefaults sets defaults for JobRun Spec.
func (js *JobRunSpec) SetDefaults() {
	// Set defaults for standalone jobRun only.
	if js.JobDefinitionRef != "" {
		return
	}

	if js.JobDefinitionSpec.ArraySpec == nil {
		js.JobDefinitionSpec.ArraySpec = pointers.StringPtr("0")
	}

	if js.JobDefinitionSpec.RetryLimit == nil {
		js.JobDefinitionSpec.RetryLimit = pointers.Int64Ptr(3)
	}

	if js.JobDefinitionSpec.MaxExecutionTime == nil {
		js.JobDefinitionSpec.MaxExecutionTime = pointers.Int64Ptr(7200)
	}
}

// RequiresDefaultingFromJobDefinition return true if JobRun refers to a jobDefinition.
func (js *JobRunSpec) RequiresDefaultingFromJobDefinition() bool {
	return js.JobDefinitionRef != ""
}

// SetDefaultsFromJobDefinition set defaults in place from its JobDefinitionRef:
// - Labels
// - Container name and image
func SetDefaultsFromJobDefinition(jr *JobRun, referredJD JobDefinition) {
	jr.AddLabel(LabelJobDefName, jr.Spec.JobDefinitionRef, false)
	jr.AddLabel(LabelJobDefUUID, string(referredJD.UID), false)

	if len(referredJD.Spec.Template.Containers) == 1 && len(jr.Spec.JobDefinitionSpec.Template.Containers) == 1 {
		if referredJD.Spec.Template.Containers[0].Name != "" && jr.Spec.JobDefinitionSpec.Template.Containers[0].Name == "" {
			jr.Spec.JobDefinitionSpec.Template.Containers[0].Name = referredJD.Spec.Template.Containers[0].Name
		}
		if referredJD.Spec.Template.Containers[0].Image != "" && jr.Spec.JobDefinitionSpec.Template.Containers[0].Image == "" {
			jr.Spec.JobDefinitionSpec.Template.Containers[0].Image = referredJD.Spec.Template.Containers[0].Image
		}
	}
}
