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

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	time "time"

	codeenginev1beta1 "github.com/rafalbigaj/code-engine-batch-job-client/pkg/apis/codeengine/v1beta1"
	versioned "github.com/rafalbigaj/code-engine-batch-job-client/pkg/client/clientset/versioned"
	internalinterfaces "github.com/rafalbigaj/code-engine-batch-job-client/pkg/client/informers/externalversions/internalinterfaces"
	v1beta1 "github.com/rafalbigaj/code-engine-batch-job-client/pkg/client/listers/codeengine/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// JobDefinitionInformer provides access to a shared informer and lister for
// JobDefinitions.
type JobDefinitionInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.JobDefinitionLister
}

type jobDefinitionInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewJobDefinitionInformer constructs a new informer for JobDefinition type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewJobDefinitionInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredJobDefinitionInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredJobDefinitionInformer constructs a new informer for JobDefinition type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredJobDefinitionInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CodeengineV1beta1().JobDefinitions(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CodeengineV1beta1().JobDefinitions(namespace).Watch(context.TODO(), options)
			},
		},
		&codeenginev1beta1.JobDefinition{},
		resyncPeriod,
		indexers,
	)
}

func (f *jobDefinitionInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredJobDefinitionInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *jobDefinitionInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&codeenginev1beta1.JobDefinition{}, f.defaultInformer)
}

func (f *jobDefinitionInformer) Lister() v1beta1.JobDefinitionLister {
	return v1beta1.NewJobDefinitionLister(f.Informer().GetIndexer())
}
