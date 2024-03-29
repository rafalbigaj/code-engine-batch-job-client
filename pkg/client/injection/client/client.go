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

// Code generated by injection-gen. DO NOT EDIT.

package client

import (
	context "context"
	json "encoding/json"
	errors "errors"
	fmt "fmt"

	v1beta1 "github.com/rafalbigaj/code-engine-batch-job-client/pkg/apis/codeengine/v1beta1"
	versioned "github.com/rafalbigaj/code-engine-batch-job-client/pkg/client/clientset/versioned"
	typedcodeenginev1beta1 "github.com/rafalbigaj/code-engine-batch-job-client/pkg/client/clientset/versioned/typed/codeengine/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	discovery "k8s.io/client-go/discovery"
	dynamic "k8s.io/client-go/dynamic"
	rest "k8s.io/client-go/rest"
	injection "knative.dev/pkg/injection"
	dynamicclient "knative.dev/pkg/injection/clients/dynamicclient"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterClient(withClientFromConfig)
	injection.Default.RegisterClientFetcher(func(ctx context.Context) interface{} {
		return Get(ctx)
	})
	injection.Dynamic.RegisterDynamicClient(withClientFromDynamic)
}

// Key is used as the key for associating information with a context.Context.
type Key struct{}

func withClientFromConfig(ctx context.Context, cfg *rest.Config) context.Context {
	return context.WithValue(ctx, Key{}, versioned.NewForConfigOrDie(cfg))
}

func withClientFromDynamic(ctx context.Context) context.Context {
	return context.WithValue(ctx, Key{}, &wrapClient{dyn: dynamicclient.Get(ctx)})
}

// Get extracts the versioned.Interface client from the context.
func Get(ctx context.Context) versioned.Interface {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		if injection.GetConfig(ctx) == nil {
			logging.FromContext(ctx).Panic(
				"Unable to fetch github.com/rafalbigaj/code-engine-batch-job-client/pkg/client/clientset/versioned.Interface from context. This context is not the application context (which is typically given to constructors via sharedmain).")
		} else {
			logging.FromContext(ctx).Panic(
				"Unable to fetch github.com/rafalbigaj/code-engine-batch-job-client/pkg/client/clientset/versioned.Interface from context.")
		}
	}
	return untyped.(versioned.Interface)
}

type wrapClient struct {
	dyn dynamic.Interface
}

var _ versioned.Interface = (*wrapClient)(nil)

func (w *wrapClient) Discovery() discovery.DiscoveryInterface {
	panic("Discovery called on dynamic client!")
}

func convert(from interface{}, to runtime.Object) error {
	bs, err := json.Marshal(from)
	if err != nil {
		return fmt.Errorf("Marshal() = %w", err)
	}
	if err := json.Unmarshal(bs, to); err != nil {
		return fmt.Errorf("Unmarshal() = %w", err)
	}
	return nil
}

// CodeengineV1beta1 retrieves the CodeengineV1beta1Client
func (w *wrapClient) CodeengineV1beta1() typedcodeenginev1beta1.CodeengineV1beta1Interface {
	return &wrapCodeengineV1beta1{
		dyn: w.dyn,
	}
}

type wrapCodeengineV1beta1 struct {
	dyn dynamic.Interface
}

func (w *wrapCodeengineV1beta1) RESTClient() rest.Interface {
	panic("RESTClient called on dynamic client!")
}

func (w *wrapCodeengineV1beta1) JobDefinitions(namespace string) typedcodeenginev1beta1.JobDefinitionInterface {
	return &wrapCodeengineV1beta1JobDefinitionImpl{
		dyn: w.dyn.Resource(schema.GroupVersionResource{
			Group:    "codeengine.cloud.ibm.com",
			Version:  "v1beta1",
			Resource: "jobdefinitions",
		}),

		namespace: namespace,
	}
}

type wrapCodeengineV1beta1JobDefinitionImpl struct {
	dyn dynamic.NamespaceableResourceInterface

	namespace string
}

var _ typedcodeenginev1beta1.JobDefinitionInterface = (*wrapCodeengineV1beta1JobDefinitionImpl)(nil)

func (w *wrapCodeengineV1beta1JobDefinitionImpl) Create(ctx context.Context, in *v1beta1.JobDefinition, opts v1.CreateOptions) (*v1beta1.JobDefinition, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "codeengine.cloud.ibm.com",
		Version: "v1beta1",
		Kind:    "JobDefinition",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Create(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.JobDefinition{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapCodeengineV1beta1JobDefinitionImpl) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return w.dyn.Namespace(w.namespace).Delete(ctx, name, opts)
}

func (w *wrapCodeengineV1beta1JobDefinitionImpl) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	return w.dyn.Namespace(w.namespace).DeleteCollection(ctx, opts, listOpts)
}

func (w *wrapCodeengineV1beta1JobDefinitionImpl) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.JobDefinition, error) {
	uo, err := w.dyn.Namespace(w.namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.JobDefinition{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapCodeengineV1beta1JobDefinitionImpl) List(ctx context.Context, opts v1.ListOptions) (*v1beta1.JobDefinitionList, error) {
	uo, err := w.dyn.Namespace(w.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.JobDefinitionList{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapCodeengineV1beta1JobDefinitionImpl) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.JobDefinition, err error) {
	uo, err := w.dyn.Namespace(w.namespace).Patch(ctx, name, pt, data, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.JobDefinition{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapCodeengineV1beta1JobDefinitionImpl) Update(ctx context.Context, in *v1beta1.JobDefinition, opts v1.UpdateOptions) (*v1beta1.JobDefinition, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "codeengine.cloud.ibm.com",
		Version: "v1beta1",
		Kind:    "JobDefinition",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Update(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.JobDefinition{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapCodeengineV1beta1JobDefinitionImpl) UpdateStatus(ctx context.Context, in *v1beta1.JobDefinition, opts v1.UpdateOptions) (*v1beta1.JobDefinition, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "codeengine.cloud.ibm.com",
		Version: "v1beta1",
		Kind:    "JobDefinition",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).UpdateStatus(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.JobDefinition{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapCodeengineV1beta1JobDefinitionImpl) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("NYI: Watch")
}

func (w *wrapCodeengineV1beta1) JobRuns(namespace string) typedcodeenginev1beta1.JobRunInterface {
	return &wrapCodeengineV1beta1JobRunImpl{
		dyn: w.dyn.Resource(schema.GroupVersionResource{
			Group:    "codeengine.cloud.ibm.com",
			Version:  "v1beta1",
			Resource: "jobruns",
		}),

		namespace: namespace,
	}
}

type wrapCodeengineV1beta1JobRunImpl struct {
	dyn dynamic.NamespaceableResourceInterface

	namespace string
}

var _ typedcodeenginev1beta1.JobRunInterface = (*wrapCodeengineV1beta1JobRunImpl)(nil)

func (w *wrapCodeengineV1beta1JobRunImpl) Create(ctx context.Context, in *v1beta1.JobRun, opts v1.CreateOptions) (*v1beta1.JobRun, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "codeengine.cloud.ibm.com",
		Version: "v1beta1",
		Kind:    "JobRun",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Create(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.JobRun{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapCodeengineV1beta1JobRunImpl) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return w.dyn.Namespace(w.namespace).Delete(ctx, name, opts)
}

func (w *wrapCodeengineV1beta1JobRunImpl) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	return w.dyn.Namespace(w.namespace).DeleteCollection(ctx, opts, listOpts)
}

func (w *wrapCodeengineV1beta1JobRunImpl) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.JobRun, error) {
	uo, err := w.dyn.Namespace(w.namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.JobRun{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapCodeengineV1beta1JobRunImpl) List(ctx context.Context, opts v1.ListOptions) (*v1beta1.JobRunList, error) {
	uo, err := w.dyn.Namespace(w.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.JobRunList{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapCodeengineV1beta1JobRunImpl) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.JobRun, err error) {
	uo, err := w.dyn.Namespace(w.namespace).Patch(ctx, name, pt, data, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.JobRun{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapCodeengineV1beta1JobRunImpl) Update(ctx context.Context, in *v1beta1.JobRun, opts v1.UpdateOptions) (*v1beta1.JobRun, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "codeengine.cloud.ibm.com",
		Version: "v1beta1",
		Kind:    "JobRun",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Update(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.JobRun{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapCodeengineV1beta1JobRunImpl) UpdateStatus(ctx context.Context, in *v1beta1.JobRun, opts v1.UpdateOptions) (*v1beta1.JobRun, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "codeengine.cloud.ibm.com",
		Version: "v1beta1",
		Kind:    "JobRun",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).UpdateStatus(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1beta1.JobRun{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapCodeengineV1beta1JobRunImpl) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("NYI: Watch")
}
