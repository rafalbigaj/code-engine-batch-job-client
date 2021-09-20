/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddLabel Add one label(key: value) to JobRun.
func (jr *JobRun) AddLabel(key, value string, overwrite bool) {
	if jr.Labels == nil {
		jr.Labels = map[string]string{}
	}
	_, ok := jr.Labels[key]
	if !ok || overwrite {
		jr.Labels[key] = value
	}
}

// SetOwner sets the given JobDefiniton as OwnerReference
// With enabled BlockOwnerDeletion
func (jr *JobRun) SetOwner(jd *JobDefinition) {
	jr.OwnerReferences = append(jr.OwnerReferences, *metav1.NewControllerRef(jd, SchemeGroupVersion.WithKind("JobDefinition")))
}
