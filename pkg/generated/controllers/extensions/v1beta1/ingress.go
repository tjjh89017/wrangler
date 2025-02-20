/*
Copyright The Kubernetes Authors.

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

// Code generated by main. DO NOT EDIT.

package v1beta1

import (
	"context"
	"time"

	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

// IngressController interface for managing Ingress resources.
type IngressController interface {
	generic.ControllerMeta
	IngressClient

	// OnChange runs the given handler when the controller detects a resource was changed.
	OnChange(ctx context.Context, name string, sync IngressHandler)

	// OnRemove runs the given handler when the controller detects a resource was changed.
	OnRemove(ctx context.Context, name string, sync IngressHandler)

	// Enqueue adds the resource with the given name to the worker queue of the controller.
	Enqueue(namespace, name string)

	// EnqueueAfter runs Enqueue after the provided duration.
	EnqueueAfter(namespace, name string, duration time.Duration)

	// Cache returns a cache for the resource type T.
	Cache() IngressCache
}

// IngressClient interface for managing Ingress resources in Kubernetes.
type IngressClient interface {
	// Create creates a new object and return the newly created Object or an error.
	Create(*v1beta1.Ingress) (*v1beta1.Ingress, error)

	// Update updates the object and return the newly updated Object or an error.
	Update(*v1beta1.Ingress) (*v1beta1.Ingress, error)
	// UpdateStatus updates the Status field of a the object and return the newly updated Object or an error.
	// Will always return an error if the object does not have a status field.
	UpdateStatus(*v1beta1.Ingress) (*v1beta1.Ingress, error)

	// Delete deletes the Object in the given name.
	Delete(namespace, name string, options *metav1.DeleteOptions) error

	// Get will attempt to retrieve the resource with the specified name.
	Get(namespace, name string, options metav1.GetOptions) (*v1beta1.Ingress, error)

	// List will attempt to find multiple resources.
	List(namespace string, opts metav1.ListOptions) (*v1beta1.IngressList, error)

	// Watch will start watching resources.
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)

	// Patch will patch the resource with the matching name.
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.Ingress, err error)
}

// IngressCache interface for retrieving Ingress resources in memory.
type IngressCache interface {
	// Get returns the resources with the specified name from the cache.
	Get(namespace, name string) (*v1beta1.Ingress, error)

	// List will attempt to find resources from the Cache.
	List(namespace string, selector labels.Selector) ([]*v1beta1.Ingress, error)

	// AddIndexer adds  a new Indexer to the cache with the provided name.
	// If you call this after you already have data in the store, the results are undefined.
	AddIndexer(indexName string, indexer IngressIndexer)

	// GetByIndex returns the stored objects whose set of indexed values
	// for the named index includes the given indexed value.
	GetByIndex(indexName, key string) ([]*v1beta1.Ingress, error)
}

// IngressHandler is function for performing any potential modifications to a Ingress resource.
type IngressHandler func(string, *v1beta1.Ingress) (*v1beta1.Ingress, error)

// IngressIndexer computes a set of indexed values for the provided object.
type IngressIndexer func(obj *v1beta1.Ingress) ([]string, error)

// IngressGenericController wraps wrangler/pkg/generic.Controller so that the function definitions adhere to IngressController interface.
type IngressGenericController struct {
	generic.ControllerInterface[*v1beta1.Ingress, *v1beta1.IngressList]
}

// OnChange runs the given resource handler when the controller detects a resource was changed.
func (c *IngressGenericController) OnChange(ctx context.Context, name string, sync IngressHandler) {
	c.ControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v1beta1.Ingress](sync))
}

// OnRemove runs the given object handler when the controller detects a resource was changed.
func (c *IngressGenericController) OnRemove(ctx context.Context, name string, sync IngressHandler) {
	c.ControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v1beta1.Ingress](sync))
}

// Cache returns a cache of resources in memory.
func (c *IngressGenericController) Cache() IngressCache {
	return &IngressGenericCache{
		c.ControllerInterface.Cache(),
	}
}

// IngressGenericCache wraps wrangler/pkg/generic.Cache so the function definitions adhere to IngressCache interface.
type IngressGenericCache struct {
	generic.CacheInterface[*v1beta1.Ingress]
}

// AddIndexer adds  a new Indexer to the cache with the provided name.
// If you call this after you already have data in the store, the results are undefined.
func (c IngressGenericCache) AddIndexer(indexName string, indexer IngressIndexer) {
	c.CacheInterface.AddIndexer(indexName, generic.Indexer[*v1beta1.Ingress](indexer))
}

type IngressStatusHandler func(obj *v1beta1.Ingress, status v1beta1.IngressStatus) (v1beta1.IngressStatus, error)

type IngressGeneratingHandler func(obj *v1beta1.Ingress, status v1beta1.IngressStatus) ([]runtime.Object, v1beta1.IngressStatus, error)

func FromIngressHandlerToHandler(sync IngressHandler) generic.Handler {
	return generic.FromObjectHandlerToHandler(generic.ObjectHandler[*v1beta1.Ingress](sync))
}

func RegisterIngressStatusHandler(ctx context.Context, controller IngressController, condition condition.Cond, name string, handler IngressStatusHandler) {
	statusHandler := &ingressStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromIngressHandlerToHandler(statusHandler.sync))
}

func RegisterIngressGeneratingHandler(ctx context.Context, controller IngressController, apply apply.Apply,
	condition condition.Cond, name string, handler IngressGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &ingressGeneratingHandler{
		IngressGeneratingHandler: handler,
		apply:                    apply,
		name:                     name,
		gvk:                      controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterIngressStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type ingressStatusHandler struct {
	client    IngressClient
	condition condition.Cond
	handler   IngressStatusHandler
}

func (a *ingressStatusHandler) sync(key string, obj *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type ingressGeneratingHandler struct {
	IngressGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *ingressGeneratingHandler) Remove(key string, obj *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1beta1.Ingress{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *ingressGeneratingHandler) Handle(obj *v1beta1.Ingress, status v1beta1.IngressStatus) (v1beta1.IngressStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.IngressGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
