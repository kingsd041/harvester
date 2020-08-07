/*
Copyright 2020 Rancher Labs, Inc.

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

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/rancher/harvester/pkg/apis/vm.cattle.io/v1alpha1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ImageHandler func(string, *v1alpha1.Image) (*v1alpha1.Image, error)

type ImageController interface {
	generic.ControllerMeta
	ImageClient

	OnChange(ctx context.Context, name string, sync ImageHandler)
	OnRemove(ctx context.Context, name string, sync ImageHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ImageCache
}

type ImageClient interface {
	Create(*v1alpha1.Image) (*v1alpha1.Image, error)
	Update(*v1alpha1.Image) (*v1alpha1.Image, error)
	UpdateStatus(*v1alpha1.Image) (*v1alpha1.Image, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.Image, error)
	List(namespace string, opts metav1.ListOptions) (*v1alpha1.ImageList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Image, err error)
}

type ImageCache interface {
	Get(namespace, name string) (*v1alpha1.Image, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha1.Image, error)

	AddIndexer(indexName string, indexer ImageIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha1.Image, error)
}

type ImageIndexer func(obj *v1alpha1.Image) ([]string, error)

type imageController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewImageController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ImageController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &imageController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromImageHandlerToHandler(sync ImageHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha1.Image
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha1.Image))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *imageController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha1.Image))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateImageDeepCopyOnChange(client ImageClient, obj *v1alpha1.Image, handler func(obj *v1alpha1.Image) (*v1alpha1.Image, error)) (*v1alpha1.Image, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *imageController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *imageController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *imageController) OnChange(ctx context.Context, name string, sync ImageHandler) {
	c.AddGenericHandler(ctx, name, FromImageHandlerToHandler(sync))
}

func (c *imageController) OnRemove(ctx context.Context, name string, sync ImageHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromImageHandlerToHandler(sync)))
}

func (c *imageController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *imageController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *imageController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *imageController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *imageController) Cache() ImageCache {
	return &imageCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *imageController) Create(obj *v1alpha1.Image) (*v1alpha1.Image, error) {
	result := &v1alpha1.Image{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *imageController) Update(obj *v1alpha1.Image) (*v1alpha1.Image, error) {
	result := &v1alpha1.Image{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *imageController) UpdateStatus(obj *v1alpha1.Image) (*v1alpha1.Image, error) {
	result := &v1alpha1.Image{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *imageController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *imageController) Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.Image, error) {
	result := &v1alpha1.Image{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *imageController) List(namespace string, opts metav1.ListOptions) (*v1alpha1.ImageList, error) {
	result := &v1alpha1.ImageList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *imageController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *imageController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1alpha1.Image, error) {
	result := &v1alpha1.Image{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type imageCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *imageCache) Get(namespace, name string) (*v1alpha1.Image, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1alpha1.Image), nil
}

func (c *imageCache) List(namespace string, selector labels.Selector) (ret []*v1alpha1.Image, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Image))
	})

	return ret, err
}

func (c *imageCache) AddIndexer(indexName string, indexer ImageIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha1.Image))
		},
	}))
}

func (c *imageCache) GetByIndex(indexName, key string) (result []*v1alpha1.Image, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1alpha1.Image, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha1.Image))
	}
	return result, nil
}

type ImageStatusHandler func(obj *v1alpha1.Image, status v1alpha1.ImageStatus) (v1alpha1.ImageStatus, error)

type ImageGeneratingHandler func(obj *v1alpha1.Image, status v1alpha1.ImageStatus) ([]runtime.Object, v1alpha1.ImageStatus, error)

func RegisterImageStatusHandler(ctx context.Context, controller ImageController, condition condition.Cond, name string, handler ImageStatusHandler) {
	statusHandler := &imageStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromImageHandlerToHandler(statusHandler.sync))
}

func RegisterImageGeneratingHandler(ctx context.Context, controller ImageController, apply apply.Apply,
	condition condition.Cond, name string, handler ImageGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &imageGeneratingHandler{
		ImageGeneratingHandler: handler,
		apply:                  apply,
		name:                   name,
		gvk:                    controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterImageStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type imageStatusHandler struct {
	client    ImageClient
	condition condition.Cond
	handler   ImageStatusHandler
}

func (a *imageStatusHandler) sync(key string, obj *v1alpha1.Image) (*v1alpha1.Image, error) {
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
		var newErr error
		obj.Status = newStatus
		obj, newErr = a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
	}
	return obj, err
}

type imageGeneratingHandler struct {
	ImageGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *imageGeneratingHandler) Remove(key string, obj *v1alpha1.Image) (*v1alpha1.Image, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1alpha1.Image{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *imageGeneratingHandler) Handle(obj *v1alpha1.Image, status v1alpha1.ImageStatus) (v1alpha1.ImageStatus, error) {
	objs, newStatus, err := a.ImageGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
