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
 * © Copyright IBM Corp. 2020, 2021
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers/internalinterfaces"

	"github.com/rafalbigaj/code-engine-batch-job-client/pkg/apis/codeengine"
)

var (
	// LabelJobIndex is the label key for job index
	LabelJobIndex = fmt.Sprintf("%s/job-index", codeengine.GroupName)
	// LabelPodType is the label key for pod type
	LabelPodType = fmt.Sprintf("%s/pod-type", codeengine.GroupName)
	// LabelJobRun is the label key for job run
	LabelJobRun = fmt.Sprintf("%s/job-run", codeengine.GroupName)
	// LabelJobDefName is the label key for job definition name
	LabelJobDefName = fmt.Sprintf("%s/job-definition-name", codeengine.GroupName)
	// LabelJobDefUUID is the label key for job definition uuid
	LabelJobDefUUID = fmt.Sprintf("%s/job-definition-uuid", codeengine.GroupName)

	// AnnotationRetryTimes is the annotation key for retry times
	AnnotationRetryTimes = fmt.Sprintf("%s/retry-times", codeengine.GroupName)
	// AnnotationPodExpectations is the annotation key for counting pod expectations
	AnnotationPodExpectations = fmt.Sprintf("%s/pod-expectations", codeengine.GroupName)
)

const (
	// JobRunType pod type for jobRun.
	JobRunType = "jobrun"
	// JobIndex env name for job index.
	JobIndex = "JOB_INDEX"

	// CodeEngine Domain
	CEDomain = "CE_DOMAIN"
	// CodeEngine SubDomain
	CESubDomain = "CE_SUBDOMAIN"
	// CodeEngine Job
	CEJob = "CE_JOB"
	// CodeEngine Job Run
	CEJobRun = "CE_JOBRUN"
	// Mode for "daemon" Job Run
	CEExecutionMode      = "CE_EXECUTION_MODE"
	CEExecutionModeValue = "DAEMON"
	// Note:
	//   When changing 'MaxIndexValue' and/or 'maxArraySize' consider to keep documentation
	//   in sync with the new values. This includes kubectl CRDs and ibmcloud-cli plugin 'code-engine'.
	//

	// MaxIndexValue the max value of each index
	MaxIndexValue = 9999999
	// maxArraySize the max number jobruns that can be run in parallel
	maxArraySize = 1000
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JobRun is a specification for a jobRun resource
type JobRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JobRunSpec   `json:"spec"`
	Status JobRunStatus `json:"status"`
}

// JobRunSpec describes how the jobRun execution will look like
type JobRunSpec struct {

	// The name of the jobDefinition that this jobRun refers
	JobDefinitionRef string `json:"jobDefinitionRef"`

	// The spec for a jobDefinition resource
	JobDefinitionSpec JobDefinitionSpec `json:"jobDefinitionSpec,omitempty"`
}

// JobRunStatus is the current status of a jobRun resource
type JobRunStatus struct {

	// Represents time when the job was acknowledged by the job controller.
	// It is not guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// Represents time when the job was completed. It is not guaranteed to
	// be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// List of JobRun indices that failed.
	// It contains
	// - single JobRun index notation: 1,3,5
	// - range JobRun index notation: 1-3,6-9
	// - mixed Jobrun index notation: 1,3-4,6,8-9
	FailedIndices *string `json:"failedIndices,omitempty"`

	// List of JobRun indices that succeeded.
	// List can be a comma-separated list of index ranges.
	SucceededIndices *string `json:"succeededIndices,omitempty"`

	// The latest available observations of an object's current state.
	// +optional
	Conditions []JobRunCondition `json:"conditions,omitempty"`

	// The number of pods which reached phase Unknown.
	// +optional
	Unknown int64 `json:"unknown,omitempty"`

	// The number of pods which reached phase Pending.
	// +optional
	Pending int64 `json:"pending,omitempty"`

	// The number of pods which reached phase Running.
	// +optional
	Running int64 `json:"running,omitempty"`

	// The number of pods which reached phase Succeeded.
	// +optional
	Succeeded int64 `json:"succeeded,omitempty"`

	// The number of pods which reached phase Failed.
	// +optional
	Failed int64 `json:"failed,omitempty"`

	// The number of pods which are requested but not created.
	// +optional
	Requested int64 `json:"requested,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JobRunList is a list of JobRunList resources
type JobRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []JobRun `json:"items"`
}

// GetCondition fetches the condition of the specified type.
func (s *JobRunStatus) GetCondition(t JobRunConditionType) *JobRunCondition {
	for _, cond := range s.Conditions {
		if cond.Type == t {
			return &cond
		}
	}
	return nil
}

// GetLatestCondition returns latest condition or nil, if the jobRun does not have a condition yet.
func (s *JobRunStatus) GetLatestCondition() *JobRunCondition {
	if len(s.Conditions) > 0 {
		return s.Conditions[len(s.Conditions)-1].DeepCopy()
	}
	return nil
}

// AddCondition sets or updates new condition on conditions,
// put down condition when update.
func (s *JobRunStatus) AddCondition(new JobRunCondition) {
	var newConditions []JobRunCondition
	for _, c := range s.Conditions {
		if c.Type != new.Type {
			newConditions = append(newConditions, c)
		} else {
			// Ignore duplicated pod events if happening.
			new.LastProbeTime = c.LastProbeTime
			new.LastTransitionTime = c.LastTransitionTime
			if reflect.DeepEqual(new, c) {
				return
			}
		}
	}

	new.LastProbeTime = metav1.Now()
	new.LastTransitionTime = metav1.Now()
	newConditions = append(newConditions, new)

	s.Conditions = newConditions
}

// CheckCondition check if the condition exist in jobRun
func (j *JobRun) CheckCondition(t JobRunConditionType) bool {
	for _, c := range j.Status.Conditions {
		if c.Type == t {
			return true
		}
	}
	return false
}

// IsJobRunFinished check if the jobRun is finished
func (j *JobRun) IsJobRunFinished() bool {
	return j.CheckCondition(JobComplete) || j.CheckCondition(JobFailed)
}

// UpdateStatusCounts updates job pods status counts.
func (j *JobRun) UpdateStatusCounts(total int64, podStatus []corev1.PodPhase) {
	var unknown, pending, running, succeeded, failed int64
	for _, ps := range podStatus {
		switch ps {
		case corev1.PodUnknown:
			unknown++
		case corev1.PodPending:
			pending++
		case corev1.PodRunning:
			running++
		case corev1.PodSucceeded:
			succeeded++
		case corev1.PodFailed:
			failed++
		}
	}

	j.Status.Unknown = unknown
	j.Status.Pending = pending
	j.Status.Running = running
	j.Status.Succeeded = succeeded
	j.Status.Failed = failed
	j.Status.Requested = total - int64(len(podStatus))
}

// UpdateFailedIndices updates jr.Status.FailedIndices.
func (j *JobRun) UpdateFailedIndices(podSnapshots map[int64]corev1.PodPhase) {
	// Don't count failed indices when no pods exists or JobRun is running in daemon mode
	if len(podSnapshots) == 0 || j.IsRunningInDaemonMode() {
		return
	}
	// indexRange is a string containing consecutive indices, such as: 2-5.
	indexRanges := generateIndexRanges(podSnapshots, false)

	failedIndices := strings.Join(indexRanges, ",")

	if len(failedIndices) != 0 && (j.Status.FailedIndices == nil || failedIndices != *j.Status.FailedIndices) {
		j.Status.FailedIndices = &failedIndices
	}
}

// GetSucceededIndices get succeeded indices map from jr.Status.SucceededIndices.
func (j *JobRun) GetSucceededIndices() map[int64]corev1.PodPhase {
	indexStatus := map[int64]corev1.PodPhase{}
	if j.Status.SucceededIndices == nil || *j.Status.SucceededIndices == "" {
		return indexStatus
	}

	succeededIndices := strings.Split(*j.Status.SucceededIndices, ",")

	return j.GetIndicesWithStatus(corev1.PodSucceeded, succeededIndices)
}

// GetFailedIndices get failed indices map from jr.Status.FailedIndices.
func (j *JobRun) GetFailedIndices() map[int64]corev1.PodPhase {
	indexStatus := map[int64]corev1.PodPhase{}
	if j.Status.FailedIndices == nil || *j.Status.FailedIndices == "" {
		return indexStatus
	}

	failedIndices := strings.Split(*j.Status.FailedIndices, ",")

	return j.GetIndicesWithStatus(corev1.PodFailed, failedIndices)
}

func (j *JobRun) GetIndicesWithStatus(status corev1.PodPhase, indices []string) map[int64]corev1.PodPhase {
	indexStatus := map[int64]corev1.PodPhase{}
	for _, ci := range indices {
		indices := strings.Split(ci, "-")
		if len(indices) == 1 {
			idx, err := strconv.ParseInt(indices[0], 10, 64)
			if err != nil {
				continue
			}
			indexStatus[idx] = status
		} else if len(indices) == 2 {
			startIdx, err := strconv.ParseInt(indices[0], 10, 64)
			if err != nil {
				continue
			}
			indexStatus[startIdx] = status

			endIdx, err := strconv.ParseInt(indices[1], 10, 64)
			if err != nil {
				continue
			}
			indexStatus[endIdx] = status
		}
	}

	return indexStatus
}

// UpdateSucceededIndices updates jr.Status.SucceededIndices.
func (j *JobRun) UpdateSucceededIndices(podSnapshots map[int64]corev1.PodPhase) {
	indexRanges := generateIndexRanges(podSnapshots, true)

	succeededIndices := strings.Join(indexRanges, ",")

	if len(succeededIndices) != 0 && (j.Status.SucceededIndices == nil || succeededIndices != *j.Status.SucceededIndices) {
		j.Status.SucceededIndices = &succeededIndices
	}
}

func generateIndexRanges(podSnapshots map[int64]corev1.PodPhase, expectSucceeded bool) (indexRanges []string) {
	var indexKeys []int64
	for idx := range podSnapshots {
		indexKeys = append(indexKeys, idx)
	}

	sort.Slice(indexKeys, func(l, u int) bool { return indexKeys[l] < indexKeys[u] })

	// Initialize first indexRange.
	for _, idx := range indexKeys {
		key := strconv.FormatInt(idx, 10)
		// Pick up index by XOR operator between non-succeeded phase and expected succeeded phase.
		if ps := podSnapshots[idx]; (ps != corev1.PodSucceeded) != expectSucceeded {

			// Initialize first indexRange.
			if len(indexRanges) == 0 {
				indexRanges = append(indexRanges, key)
				continue
			}

			// Merge consecutive indices.
			lastIndexRange := indexRanges[len(indexRanges)-1]
			// Slice `indices` at least have one index number.
			indices := strings.Split(lastIndexRange, "-")
			lastIndex, err := strconv.ParseInt(indices[len(indices)-1], 10, 64)
			if err == nil && lastIndex+1 == idx {
				replaceRange := fmt.Sprintf("%s-%d", indices[0], idx)
				indexRanges = append(indexRanges[:len(indexRanges)-1], replaceRange)
			} else {
				indexRanges = append(indexRanges, key)
			}
		}
	}

	return
}

// JobRunPodLabelFilterOption generate a TweakListOptionsFunc for jobRun pods
func JobRunPodLabelFilterOption() internalinterfaces.TweakListOptionsFunc {
	return func(l *metav1.ListOptions) {
		if l == nil {
			l = &metav1.ListOptions{}
		}
		l.LabelSelector = fmt.Sprintf("%s=%s", LabelPodType, JobRunType)
	}
}

func (j *JobRun) IsRunningInDaemonMode() bool {
	for _, envVar := range j.Spec.JobDefinitionSpec.Template.Containers[0].Env {
		if envVar.Name == CEExecutionMode {
			return envVar.Value == CEExecutionModeValue
		}
	}
	return false
}
