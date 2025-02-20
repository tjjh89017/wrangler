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

package v1

import (
	"context"
	"time"

	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

// PersistentVolumeController interface for managing PersistentVolume resources.
type PersistentVolumeController interface {
	generic.ControllerMeta
	PersistentVolumeClient

	// OnChange runs the given handler when the controller detects a resource was changed.
	OnChange(ctx context.Context, name string, sync PersistentVolumeHandler)

	// OnRemove runs the given handler when the controller detects a resource was changed.
	OnRemove(ctx context.Context, name string, sync PersistentVolumeHandler)

	// Enqueue adds the resource with the given name to the worker queue of the controller.
	Enqueue(name string)

	// EnqueueAfter runs Enqueue after the provided duration.
	EnqueueAfter(name string, duration time.Duration)

	// Cache returns a cache for the resource type T.
	Cache() PersistentVolumeCache
}

// PersistentVolumeClient interface for managing PersistentVolume resources in Kubernetes.
type PersistentVolumeClient interface {
	// Create creates a new object and return the newly created Object or an error.
	Create(*v1.PersistentVolume) (*v1.PersistentVolume, error)

	// Update updates the object and return the newly updated Object or an error.
	Update(*v1.PersistentVolume) (*v1.PersistentVolume, error)
	// UpdateStatus updates the Status field of a the object and return the newly updated Object or an error.
	// Will always return an error if the object does not have a status field.
	UpdateStatus(*v1.PersistentVolume) (*v1.PersistentVolume, error)

	// Delete deletes the Object in the given name.
	Delete(name string, options *metav1.DeleteOptions) error

	// Get will attempt to retrieve the resource with the specified name.
	Get(name string, options metav1.GetOptions) (*v1.PersistentVolume, error)

	// List will attempt to find multiple resources.
	List(opts metav1.ListOptions) (*v1.PersistentVolumeList, error)

	// Watch will start watching resources.
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	// Patch will patch the resource with the matching name.
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.PersistentVolume, err error)
}

// PersistentVolumeCache interface for retrieving PersistentVolume resources in memory.
type PersistentVolumeCache interface {
	// Get returns the resources with the specified name from the cache.
	Get(name string) (*v1.PersistentVolume, error)

	// List will attempt to find resources from the Cache.
	List(selector labels.Selector) ([]*v1.PersistentVolume, error)

	// AddIndexer adds  a new Indexer to the cache with the provided name.
	// If you call this after you already have data in the store, the results are undefined.
	AddIndexer(indexName string, indexer PersistentVolumeIndexer)

	// GetByIndex returns the stored objects whose set of indexed values
	// for the named index includes the given indexed value.
	GetByIndex(indexName, key string) ([]*v1.PersistentVolume, error)
}

// PersistentVolumeHandler is function for performing any potential modifications to a PersistentVolume resource.
type PersistentVolumeHandler func(string, *v1.PersistentVolume) (*v1.PersistentVolume, error)

// PersistentVolumeIndexer computes a set of indexed values for the provided object.
type PersistentVolumeIndexer func(obj *v1.PersistentVolume) ([]string, error)

// PersistentVolumeGenericController wraps wrangler/pkg/generic.NonNamespacedController so that the function definitions adhere to PersistentVolumeController interface.
type PersistentVolumeGenericController struct {
	generic.NonNamespacedControllerInterface[*v1.PersistentVolume, *v1.PersistentVolumeList]
}

// OnChange runs the given resource handler when the controller detects a resource was changed.
func (c *PersistentVolumeGenericController) OnChange(ctx context.Context, name string, sync PersistentVolumeHandler) {
	c.NonNamespacedControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v1.PersistentVolume](sync))
}

// OnRemove runs the given object handler when the controller detects a resource was changed.
func (c *PersistentVolumeGenericController) OnRemove(ctx context.Context, name string, sync PersistentVolumeHandler) {
	c.NonNamespacedControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v1.PersistentVolume](sync))
}

// Cache returns a cache of resources in memory.
func (c *PersistentVolumeGenericController) Cache() PersistentVolumeCache {
	return &PersistentVolumeGenericCache{
		c.NonNamespacedControllerInterface.Cache(),
	}
}

// PersistentVolumeGenericCache wraps wrangler/pkg/generic.NonNamespacedCache so the function definitions adhere to PersistentVolumeCache interface.
type PersistentVolumeGenericCache struct {
	generic.NonNamespacedCacheInterface[*v1.PersistentVolume]
}

// AddIndexer adds  a new Indexer to the cache with the provided name.
// If you call this after you already have data in the store, the results are undefined.
func (c PersistentVolumeGenericCache) AddIndexer(indexName string, indexer PersistentVolumeIndexer) {
	c.NonNamespacedCacheInterface.AddIndexer(indexName, generic.Indexer[*v1.PersistentVolume](indexer))
}

type PersistentVolumeStatusHandler func(obj *v1.PersistentVolume, status v1.PersistentVolumeStatus) (v1.PersistentVolumeStatus, error)

type PersistentVolumeGeneratingHandler func(obj *v1.PersistentVolume, status v1.PersistentVolumeStatus) ([]runtime.Object, v1.PersistentVolumeStatus, error)

func FromPersistentVolumeHandlerToHandler(sync PersistentVolumeHandler) generic.Handler {
	return generic.FromObjectHandlerToHandler(generic.ObjectHandler[*v1.PersistentVolume](sync))
}

func RegisterPersistentVolumeStatusHandler(ctx context.Context, controller PersistentVolumeController, condition condition.Cond, name string, handler PersistentVolumeStatusHandler) {
	statusHandler := &persistentVolumeStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromPersistentVolumeHandlerToHandler(statusHandler.sync))
}

func RegisterPersistentVolumeGeneratingHandler(ctx context.Context, controller PersistentVolumeController, apply apply.Apply,
	condition condition.Cond, name string, handler PersistentVolumeGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &persistentVolumeGeneratingHandler{
		PersistentVolumeGeneratingHandler: handler,
		apply:                             apply,
		name:                              name,
		gvk:                               controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterPersistentVolumeStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type persistentVolumeStatusHandler struct {
	client    PersistentVolumeClient
	condition condition.Cond
	handler   PersistentVolumeStatusHandler
}

func (a *persistentVolumeStatusHandler) sync(key string, obj *v1.PersistentVolume) (*v1.PersistentVolume, error) {
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

type persistentVolumeGeneratingHandler struct {
	PersistentVolumeGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *persistentVolumeGeneratingHandler) Remove(key string, obj *v1.PersistentVolume) (*v1.PersistentVolume, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.PersistentVolume{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *persistentVolumeGeneratingHandler) Handle(obj *v1.PersistentVolume, status v1.PersistentVolumeStatus) (v1.PersistentVolumeStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.PersistentVolumeGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
