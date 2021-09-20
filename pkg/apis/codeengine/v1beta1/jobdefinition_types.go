/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*******************************************************************************
 * Portions of this file are subject to the following notice:
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JobDefinition is a specification for a JobDefinition resource
type JobDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JobDefinitionSpec   `json:"spec"`
	Status JobDefinitionStatus `json:"status,omitempty"`
}

// JobDefinitionSpec is the spec for a jobDefinition resource
type JobDefinitionSpec struct {
	// Specifies the indices of pods to be created
	// White spaces are allowed
	// Example values are:
	// 1,3,6,9
	// 1-5, 7 - 8, 10
	ArraySpec *string `json:"arraySpec,omitempty"`

	// Number of retries before marking this job failed.
	// The retry times will be RetryLimit + 1 if not specified.
	// Default value is 3 if not set explicitly.
	RetryLimit *int64 `json:"retryLimit,omitempty"`

	// Specifies the duration in seconds relative to the startTime that the job may be active
	// before the system tries to terminate it. Value must be positive integer
	MaxExecutionTime *int64 `json:"maxExecutionTime,omitempty"`

	// Specifies the template for creating copies of a pod
	Template JobPodTemplate `json:"template,omitempty"`
}

// JobDefinitionStatus is the current status of a jobDefinition resource
type JobDefinitionStatus struct {
	// Address holds the information needed for a Route to be the target of an event.
	// Read-only, User cannot set it or update it
	Address *duckv1.Addressable `json:"address,omitempty"`
}

// JobPodTemplate is the template of pod resource
type JobPodTemplate struct {
	// List of containers belonging to the pod.
	// Containers cannot currently be added or removed.
	// There must be at least one container in a Pod.
	// Cannot be updated.
	Containers []corev1.Container `json:"containers,omitempty"`

	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use. For example,
	// in the case of docker, only DockerConfig type secrets are honored.
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// ServiceAccountName is an optional string of the service account name which the generated pods should belong to
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JobDefinitionList is a list of JobDefinition resources
type JobDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []JobDefinition `json:"items"`
}

// IsEmpty returns if the JobDefinitionSpec is empty
func (jds *JobDefinitionSpec) IsEmpty() bool {
	empty := &JobDefinitionSpec{}
	return equality.Semantic.DeepEqual(jds, empty)
}

// IsEmpty returns if the JobPodTemplate is empty
func (t *JobPodTemplate) IsEmpty() bool {
	empty := &JobPodTemplate{}
	return equality.Semantic.DeepEqual(t, empty)
}

// IsEmpty returns if the JobDefinitionStatus is empty
func (jds *JobDefinitionStatus) IsEmpty() bool {
	empty := &JobDefinitionStatus{}
	return equality.Semantic.DeepEqual(jds, empty)
}
