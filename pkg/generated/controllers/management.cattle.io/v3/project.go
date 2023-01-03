/*
Copyright 2023 Rancher Labs, Inc.

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

package v3

import (
	"context"
	"time"

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
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

type ProjectHandler func(string, *v3.Project) (*v3.Project, error)

type ProjectController interface {
	generic.ControllerMeta
	ProjectClient

	OnChange(ctx context.Context, name string, sync ProjectHandler)
	OnRemove(ctx context.Context, name string, sync ProjectHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ProjectCache
}

type ProjectClient interface {
	Create(*v3.Project) (*v3.Project, error)
	Update(*v3.Project) (*v3.Project, error)
	UpdateStatus(*v3.Project) (*v3.Project, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.Project, error)
	List(namespace string, opts metav1.ListOptions) (*v3.ProjectList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.Project, err error)
}

type ProjectCache interface {
	Get(namespace, name string) (*v3.Project, error)
	List(namespace string, selector labels.Selector) ([]*v3.Project, error)

	AddIndexer(indexName string, indexer ProjectIndexer)
	GetByIndex(indexName, key string) ([]*v3.Project, error)
}

type ProjectIndexer func(obj *v3.Project) ([]string, error)

type projectController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewProjectController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ProjectController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &projectController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromProjectHandlerToHandler(sync ProjectHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.Project
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.Project))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *projectController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.Project))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateProjectDeepCopyOnChange(client ProjectClient, obj *v3.Project, handler func(obj *v3.Project) (*v3.Project, error)) (*v3.Project, error) {
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

func (c *projectController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *projectController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *projectController) OnChange(ctx context.Context, name string, sync ProjectHandler) {
	c.AddGenericHandler(ctx, name, FromProjectHandlerToHandler(sync))
}

func (c *projectController) OnRemove(ctx context.Context, name string, sync ProjectHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromProjectHandlerToHandler(sync)))
}

func (c *projectController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *projectController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *projectController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *projectController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *projectController) Cache() ProjectCache {
	return &projectCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *projectController) Create(obj *v3.Project) (*v3.Project, error) {
	result := &v3.Project{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *projectController) Update(obj *v3.Project) (*v3.Project, error) {
	result := &v3.Project{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *projectController) UpdateStatus(obj *v3.Project) (*v3.Project, error) {
	result := &v3.Project{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *projectController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *projectController) Get(namespace, name string, options metav1.GetOptions) (*v3.Project, error) {
	result := &v3.Project{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *projectController) List(namespace string, opts metav1.ListOptions) (*v3.ProjectList, error) {
	result := &v3.ProjectList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *projectController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *projectController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.Project, error) {
	result := &v3.Project{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type projectCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *projectCache) Get(namespace, name string) (*v3.Project, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.Project), nil
}

func (c *projectCache) List(namespace string, selector labels.Selector) (ret []*v3.Project, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.Project))
	})

	return ret, err
}

func (c *projectCache) AddIndexer(indexName string, indexer ProjectIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.Project))
		},
	}))
}

func (c *projectCache) GetByIndex(indexName, key string) (result []*v3.Project, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.Project, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.Project))
	}
	return result, nil
}

type ProjectStatusHandler func(obj *v3.Project, status v3.ProjectStatus) (v3.ProjectStatus, error)

type ProjectGeneratingHandler func(obj *v3.Project, status v3.ProjectStatus) ([]runtime.Object, v3.ProjectStatus, error)

func RegisterProjectStatusHandler(ctx context.Context, controller ProjectController, condition condition.Cond, name string, handler ProjectStatusHandler) {
	statusHandler := &projectStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromProjectHandlerToHandler(statusHandler.sync))
}

func RegisterProjectGeneratingHandler(ctx context.Context, controller ProjectController, apply apply.Apply,
	condition condition.Cond, name string, handler ProjectGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &projectGeneratingHandler{
		ProjectGeneratingHandler: handler,
		apply:                    apply,
		name:                     name,
		gvk:                      controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterProjectStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type projectStatusHandler struct {
	client    ProjectClient
	condition condition.Cond
	handler   ProjectStatusHandler
}

func (a *projectStatusHandler) sync(key string, obj *v3.Project) (*v3.Project, error) {
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

type projectGeneratingHandler struct {
	ProjectGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *projectGeneratingHandler) Remove(key string, obj *v3.Project) (*v3.Project, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.Project{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *projectGeneratingHandler) Handle(obj *v3.Project, status v3.ProjectStatus) (v3.ProjectStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.ProjectGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
