/*
Copyright 2020 The Knative Authors

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

// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "github.com/rafal-bigaj/code-engine-batch-job-client/pkg/apis/codeengine/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// JobRunLister helps list JobRuns.
// All objects returned here must be treated as read-only.
type JobRunLister interface {
	// List lists all JobRuns in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.JobRun, err error)
	// JobRuns returns an object that can list and get JobRuns.
	JobRuns(namespace string) JobRunNamespaceLister
	JobRunListerExpansion
}

// jobRunLister implements the JobRunLister interface.
type jobRunLister struct {
	indexer cache.Indexer
}

// NewJobRunLister returns a new JobRunLister.
func NewJobRunLister(indexer cache.Indexer) JobRunLister {
	return &jobRunLister{indexer: indexer}
}

// List lists all JobRuns in the indexer.
func (s *jobRunLister) List(selector labels.Selector) (ret []*v1beta1.JobRun, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.JobRun))
	})
	return ret, err
}

// JobRuns returns an object that can list and get JobRuns.
func (s *jobRunLister) JobRuns(namespace string) JobRunNamespaceLister {
	return jobRunNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// JobRunNamespaceLister helps list and get JobRuns.
// All objects returned here must be treated as read-only.
type JobRunNamespaceLister interface {
	// List lists all JobRuns in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.JobRun, err error)
	// Get retrieves the JobRun from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.JobRun, error)
	JobRunNamespaceListerExpansion
}

// jobRunNamespaceLister implements the JobRunNamespaceLister
// interface.
type jobRunNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all JobRuns in the indexer for a given namespace.
func (s jobRunNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.JobRun, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.JobRun))
	})
	return ret, err
}

// Get retrieves the JobRun from the indexer for a given namespace and name.
func (s jobRunNamespaceLister) Get(name string) (*v1beta1.JobRun, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("jobrun"), name)
	}
	return obj.(*v1beta1.JobRun), nil
}