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
	"fmt"
	"time"

	"github.com/go-logr/logr"
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
	openedx := &cachev1.Openedx{}

	err := r.Client.Get(context.TODO(), req.NamespacedName, openedx)

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
	r.Log.Info("This operator only works with openedx namespace")

	namespaceName := openedx.Name
	if namespaceName == "openedx" {
		r.Log.Info("Using openedx namespace")
	} else {
		r.Log.Info("Using ", namespaceName, " namespace")
	}

	var result *reconcile.Result

	// == namespace ======================

	result, err = r.ensureNamespace(req, openedx, r.namespace(openedx))
	if result != nil {
		return *result, err
	}

	// == Persistent Volume Claim ========

	result, err = r.ensurePVC(req, openedx, r.persistencevolumeclaim("elasticsearch", "2Gi", openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensurePVC(req, openedx, r.persistencevolumeclaim("caddy", "1Gi", openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensurePVC(req, openedx, r.persistencevolumeclaim("mongodb", "5Gi", openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensurePVC(req, openedx, r.persistencevolumeclaim("mysql", "5Gi", openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensurePVC(req, openedx, r.persistencevolumeclaim("redis", "1Gi", openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensurePVC(req, openedx, r.persistencevolumeclaim("demo", "1Gi", openedx))
	if result != nil {
		return *result, err
	}

	// == ConfigMap ========

	result, err = r.ensureConfigMap(req, openedx, r.openedxConfig(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureConfigMap(req, openedx, r.openedxSettingsCmsConfig(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureConfigMap(req, openedx, r.openedxSettingsLmsConfig(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureConfigMap(req, openedx, r.nginxConfig(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureConfigMap(req, openedx, r.mysqlInitdbConfig(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureConfigMap(req, openedx, r.caddyConfig(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureConfigMap(req, openedx, r.redisConfig(openedx))
	if result != nil {
		return *result, err
	}

	// == SERVICE ========

	result, err = r.ensureService(req, openedx, r.cmsService(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, openedx, r.elasticsearchService(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, openedx, r.forumService(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, openedx, r.lmsService(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, openedx, r.redisService(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, openedx, r.mongodbService(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, openedx, r.mysqlService(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, openedx, r.nginxService(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, openedx, r.caddyService(openedx))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, openedx, r.smtpService(openedx))
	if result != nil {
		return *result, err
	}

	// == CADDY ========
	result, err = r.ensureDeployment(req, openedx, r.caddyDeployment(openedx))
	if result != nil {
		return *result, err
	}

	// == CMS WORKER ==========
	result, err = r.ensureDeployment(req, openedx, r.cmsworkerDeployment(openedx))
	if result != nil {
		return *result, err
	}

	// == CMS  ==========
	result, err = r.ensureDeployment(req, openedx, r.cmsDeployment(openedx))
	if result != nil {
		return *result, err
	}

	// == ELASTICSEARCH ========
	result, err = r.ensureDeployment(req, openedx, r.elasticsearchDeployment(openedx))
	if result != nil {
		return *result, err
	}

	// == FORUM ========
	result, err = r.ensureDeployment(req, openedx, r.forumDeployment(openedx))
	if result != nil {
		return *result, err
	}

	// == LMS WORKER ==========
	result, err = r.ensureDeployment(req, openedx, r.lmsworkerDeployment(openedx))
	if result != nil {
		return *result, err
	}
	// == LMS  ==========

	result, err = r.ensureDeployment(req, openedx, r.lmsDeployment(openedx))
	if result != nil {
		return *result, err
	}

	// == REDIS ========
	result, err = r.ensureDeployment(req, openedx, r.redisDeployment(openedx))
	if result != nil {
		return *result, err
	}

	// == MONGODB ========
	result, err = r.ensureDeployment(req, openedx, r.mongodbDeployment(openedx))
	if result != nil {
		return *result, err
	}

	// == MYSQL ========
	result, err = r.ensureDeployment(req, openedx, r.mysqlDeployment(openedx))
	if result != nil {
		return *result, err
	}

	// == NGINX ========
	result, err = r.ensureDeployment(req, openedx, r.nginxDeployment(openedx))
	if result != nil {
		return *result, err
	}

	// == SMTP ========
	result, err = r.ensureDeployment(req, openedx, r.smtpDeployment(openedx))
	if result != nil {
		return *result, err
	}

	// == JOB =======

	//== LMS Job ========
	result, err = r.ensureJob(req, openedx, r.lmsJob(openedx))
	if result != nil {
		return *result, err
	}

	lmsjobComplete := r.isLmsJobDone(openedx)

	if !lmsjobComplete {
		// If lmsJob isn't complete, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(15)

		r.Log.Info(fmt.Sprintf("Lms Job isn't Complete, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	//== CMS Job ========
	result, err = r.ensureJob(req, openedx, r.cmsJob(openedx))
	if result != nil {
		return *result, err
	}

	cmsjobComplete := r.isCmsJobDone(openedx)

	if !cmsjobComplete {
		// If cmsJob isn't complete, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(15)

		r.Log.Info(fmt.Sprintf("Cms Job isn't Complete, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	//== Forum Job ======
	result, err = r.ensureJob(req, openedx, r.forumJob(openedx))
	if result != nil {
		return *result, err
	}

	forumjobComplete := r.isForumJobDone(openedx)

	if !forumjobComplete {
		// If forumJob isn't complete, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(15)

		r.Log.Info(fmt.Sprintf("Forum Job isn't Complete, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	// == INGRESS ==========

	result, err = r.ensureIngress(req, openedx, r.ingress("web", openedx))
	if result != nil {
		return *result, err
	}

	//== Demo Job ========
	result, err = r.ensureJob(req, openedx, r.demoJob(openedx))
	if result != nil {
		return *result, err
	}

	demojobComplete := r.isDemoJobDone(openedx)

	if !demojobComplete {
		// If demoJob isn't complete, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(15)

		r.Log.Info(fmt.Sprintf("Demo Job isn't Complete, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
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
