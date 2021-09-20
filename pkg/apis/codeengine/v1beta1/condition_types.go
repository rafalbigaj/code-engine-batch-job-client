/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JobRunConditionType are the condition types of JobRun
type JobRunConditionType string

// These are valid conditions of a job.
const (
	// JobPending means the job has submitted into the system.
	JobPending JobRunConditionType = "Pending"
	// JobRunning means at least one job's pod running.
	JobRunning JobRunConditionType = "Running"
	// JobComplete means the job has completed its execution.
	JobComplete JobRunConditionType = "Complete"
	// JobFailed means the job has failed its execution.
	JobFailed JobRunConditionType = "Failed"
)

// JobRunCondition describes current state of a jobRun.
type JobRunCondition struct {
	// Type of JobRun condition, Complete or Failed.
	Type JobRunConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=JobConditionType"`
	// Status of the condition, one of True, False, Unknown.
	Status v1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`
	// Last time the condition was checked.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty" protobuf:"bytes,3,opt,name=lastProbeTime"`
	// Last time the condition transit from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// (brief) reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// Human readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}
