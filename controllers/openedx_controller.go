/*


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

package controllers

import (
	"context"
	// "fmt"
	// "time"

	"github.com/go-logr/logr"
	// "github.com/prometheus/common/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
)

// blank assignment to verify that Openedx implements reconcile.Reconciler
var _ reconcile.Reconciler = &OpenedxReconciler{}

// OpenedxReconciler reconciles a Openedx object
// comment
type OpenedxReconciler struct {
	Client client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.operatortrain.me,resources=openedxes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.operatortrain.me,resources=openedxes/status,verbs=get;update;patch
// comment

func (r *OpenedxReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("openedx", req.NamespacedName)

	// Fetch the Openedx instance
	instance := &cachev1.Openedx{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	var result *reconcile.Result

	// == ConfigMap ========

	result, err = r.ensureConfigMap(req, instance, r.ConfigMap(instance))
	if result != nil {
		return *result, err
	}

	// == LMS  ==========

	result, err = r.ensureDeployment(req, instance, r.lmsDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.lmsService(instance))
	if result != nil {
		return *result, err
	}

	// == LMS WORKER ==========
	result, err = r.ensureDeployment(req, instance, r.lmsworkerDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == CMS  ==========
	result, err = r.ensureDeployment(req, instance, r.cmsDeployment(instance))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureService(req, instance, r.cmsService(instance))
	if result != nil {
		return *result, err
	}

	// == CMS WORKER ==========
	result, err = r.ensureDeployment(req, instance, r.cmsworkerDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == MYSQL ========
	result, err = r.ensureDeployment(req, instance, r.mysqlDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.mysqlService(instance))
	if result != nil {
		return *result, err
	}
	// mysqlRunning := r.isMysqlUp(instance)

	// if !mysqlRunning {
	// 	// If nginx isn't running yet, requeue the reconcile
	// 	// to run again after a delay
	// 	delay := time.Second * time.Duration(5)

	// 	log.Info(fmt.Sprintf("MYSQL isn't running, waiting for %s", delay))
	// 	return reconcile.Result{RequeueAfter: delay}, nil
	// }

	// == MONODB ========
	result, err = r.ensureDeployment(req, instance, r.mongodbDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.mongodbService(instance))
	if result != nil {
		return *result, err
	}
	// mongodbRunning := r.isMongodbUp(instance)

	// if !mongodbRunning {
	// 	// If nginx isn't running yet, requeue the reconcile
	// 	// to run again after a delay
	// 	delay := time.Second * time.Duration(5)

	// 	log.Info(fmt.Sprintf("MONGODB isn't running, waiting for %s", delay))
	// 	return reconcile.Result{RequeueAfter: delay}, nil
	// }

	// == NGINX ========
	result, err = r.ensureDeployment(req, instance, r.nginxDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.nginxService(instance))
	if result != nil {
		return *result, err
	}
	// nginxRunning := r.isNginxUp(instance)

	// if !nginxRunning {
	// 	// If nginx isn't running yet, requeue the reconcile
	// 	// to run again after a delay
	// 	delay := time.Second * time.Duration(5)

	// 	log.Info(fmt.Sprintf("NGINX isn't running, waiting for %s", delay))
	// 	return reconcile.Result{RequeueAfter: delay}, nil
	// }

	// == MEMCACHED ========
	result, err = r.ensureDeployment(req, instance, r.memcachedDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.memcachedService(instance))
	if result != nil {
		return *result, err
	}
	// memcachedRunning := r.isMemcachedUp(instance)

	// if !memcachedRunning {
	// 	// If nginx isn't running yet, requeue the reconcile
	// 	// to run again after a delay
	// 	delay := time.Second * time.Duration(5)

	// 	log.Info(fmt.Sprintf("NGINX isn't running, waiting for %s", delay))
	// 	return reconcile.Result{RequeueAfter: delay}, nil
	// }

	// == RABBITMQ ========
	result, err = r.ensureDeployment(req, instance, r.rabbitmqDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.rabbitmqService(instance))
	if result != nil {
		return *result, err
	}

	// == SMTP ========
	result, err = r.ensureDeployment(req, instance, r.smtpDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.smtpService(instance))
	if result != nil {
		return *result, err
	}

	// == ELASTICSEARCH ========
	result, err = r.ensureDeployment(req, instance, r.elasticsearchDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.elasticsearchService(instance))
	if result != nil {
		return *result, err
	}

	// == FORUM ========
	result, err = r.ensureDeployment(req, instance, r.forumDeployment(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.forumService(instance))
	if result != nil {
		return *result, err
	}

	// == Finish ==========
	// Everything went fine, don't requeue
	return ctrl.Result{}, nil
}

// add comment

func (r *OpenedxReconciler) SetupWithManager(mgr ctrl.Manager) error {

	// Create a new controller
	c, err := controller.New("openedx-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Openedx
	err = c.Watch(&source.Kind{Type: &cachev1.Openedx{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to Deployment
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1.Openedx{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to Service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1.Openedx{},
	})
	if err != nil {
		return err
	}

	// Watch for change to pods
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1.Openedx{},
	})
	if err != nil {
		return err
	}

	return nil
}
