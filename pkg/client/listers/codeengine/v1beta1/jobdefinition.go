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

// JobDefinitionLister helps list JobDefinitions.
// All objects returned here must be treated as read-only.
type JobDefinitionLister interface {
	// List lists all JobDefinitions in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.JobDefinition, err error)
	// JobDefinitions returns an object that can list and get JobDefinitions.
	JobDefinitions(namespace string) JobDefinitionNamespaceLister
	JobDefinitionListerExpansion
}

// jobDefinitionLister implements the JobDefinitionLister interface.
type jobDefinitionLister struct {
	indexer cache.Indexer
}

// NewJobDefinitionLister returns a new JobDefinitionLister.
func NewJobDefinitionLister(indexer cache.Indexer) JobDefinitionLister {
	return &jobDefinitionLister{indexer: indexer}
}

// List lists all JobDefinitions in the indexer.
func (s *jobDefinitionLister) List(selector labels.Selector) (ret []*v1beta1.JobDefinition, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.JobDefinition))
	})
	return ret, err
}

// JobDefinitions returns an object that can list and get JobDefinitions.
func (s *jobDefinitionLister) JobDefinitions(namespace string) JobDefinitionNamespaceLister {
	return jobDefinitionNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// JobDefinitionNamespaceLister helps list and get JobDefinitions.
// All objects returned here must be treated as read-only.
type JobDefinitionNamespaceLister interface {
	// List lists all JobDefinitions in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.JobDefinition, err error)
	// Get retrieves the JobDefinition from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.JobDefinition, error)
	JobDefinitionNamespaceListerExpansion
}

// jobDefinitionNamespaceLister implements the JobDefinitionNamespaceLister
// interface.
type jobDefinitionNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all JobDefinitions in the indexer for a given namespace.
func (s jobDefinitionNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.JobDefinition, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.JobDefinition))
	})
	return ret, err
}

// Get retrieves the JobDefinition from the indexer for a given namespace and name.
func (s jobDefinitionNamespaceLister) Get(name string) (*v1beta1.JobDefinition, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("jobdefinition"), name)
	}
	return obj.(*v1beta1.JobDefinition), nil
}