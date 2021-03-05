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
	r.Log.Info("This operator only works with openedx namespace")

	namespaceName := instance.Name
	if namespaceName == "openedx" {
		r.Log.Info("Using openedx namespace")
	} else {
		r.Log.Info("Using ", namespaceName, " namespace")
	}

	var result *reconcile.Result

	// == namespace ======================

	result, err = r.ensureNamespace(req, instance, r.namespace(instance))
	if result != nil {
		return *result, err
	}

	// == Persistent Volume Claim ========

	result, err = r.ensurePVC(req, instance, r.persistencevolumeclaim("elasticsearch", "2Gi", instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensurePVC(req, instance, r.persistencevolumeclaim("caddy", "1Gi", instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensurePVC(req, instance, r.persistencevolumeclaim("mongodb", "5Gi", instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensurePVC(req, instance, r.persistencevolumeclaim("mysql", "5Gi", instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensurePVC(req, instance, r.persistencevolumeclaim("redis", "1Gi", instance))
	if result != nil {
		return *result, err
	}

	// == ConfigMap ========

	result, err = r.ensureConfigMap(req, instance, r.openedxConfig(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureConfigMap(req, instance, r.openedxSettingsCmsConfig(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureConfigMap(req, instance, r.openedxSettingsLmsConfig(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureConfigMap(req, instance, r.nginxConfig(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureConfigMap(req, instance, r.mysqlInitdbConfig(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureConfigMap(req, instance, r.caddyConfig(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureConfigMap(req, instance, r.redisConfig(instance))
	if result != nil {
		return *result, err
	}

	// == SERVICE ========

	result, err = r.ensureService(req, instance, r.cmsService(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.elasticsearchService(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.forumService(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.lmsService(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.redisService(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.mongodbService(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.mysqlService(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.nginxService(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.caddyService(instance))
	if result != nil {
		return *result, err
	}

	result, err = r.ensureService(req, instance, r.smtpService(instance))
	if result != nil {
		return *result, err
	}

	// == CMS WORKER ==========
	result, err = r.ensureDeployment(req, instance, r.cmsworkerDeployment(instance))
	if result != nil {
		return *result, err
	}

	cmsworkerRunning := r.isCmsworkerUp(instance)

	if !cmsworkerRunning {
		// If cmsworker isn't running yet, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(5)

		r.Log.Info(fmt.Sprintf("Cmsworker isn't running, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	// == CMS  ==========
	result, err = r.ensureDeployment(req, instance, r.cmsDeployment(instance))
	if result != nil {
		return *result, err
	}
	cmsRunning := r.isCmsUp(instance)

	if !cmsRunning {
		// If cms isn't running yet, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(5)

		r.Log.Info(fmt.Sprintf("CMS isn't running, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	// == ELASTICSEARCH ========
	result, err = r.ensureDeployment(req, instance, r.elasticsearchDeployment(instance))
	if result != nil {
		return *result, err
	}

	elRunning := r.iselasticsearchUp(instance)

	if !elRunning {
		// If elasticsearch isn't running yet, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(5)

		r.Log.Info(fmt.Sprintf("Elasticsearch isn't running, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	// == FORUM ========
	result, err = r.ensureDeployment(req, instance, r.forumDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == LMS WORKER ==========
	result, err = r.ensureDeployment(req, instance, r.lmsworkerDeployment(instance))
	if result != nil {
		return *result, err
	}
	// == LMS  ==========

	result, err = r.ensureDeployment(req, instance, r.lmsDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == REDIS ========
	result, err = r.ensureDeployment(req, instance, r.redisDeployment(instance))
	if result != nil {
		return *result, err
	}

	redisRunning := r.isRedisdUp(instance)

	if !redisRunning {
		// If redis isn't running yet, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(5)

		r.Log.Info(fmt.Sprintf("redis isn't running, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	// == MONGODB ========
	result, err = r.ensureDeployment(req, instance, r.mongodbDeployment(instance))
	if result != nil {
		return *result, err
	}

	mongodbRunning := r.isMongodbUp(instance)

	if !mongodbRunning {
		// If Mongodb isn't running yet, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(5)

		r.Log.Info(fmt.Sprintf("MONGODB isn't running, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	// == MYSQL ========
	result, err = r.ensureDeployment(req, instance, r.mysqlDeployment(instance))
	if result != nil {
		return *result, err
	}

	mysqlRunning := r.isMysqlUp(instance)

	if !mysqlRunning {
		// If mysql isn't running yet, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(5)

		r.Log.Info(fmt.Sprintf("Mysql isn't running, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	// == NGINX ========
	result, err = r.ensureDeployment(req, instance, r.nginxDeployment(instance))
	if result != nil {
		return *result, err
	}

	nginxRunning := r.isNginxUp(instance)

	if !nginxRunning {
		// If nginx isn't running yet, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(5)

		r.Log.Info(fmt.Sprintf("NGINX isn't running, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	// == CADDY ========
	result, err = r.ensureDeployment(req, instance, r.caddyDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == SMTP ========
	result, err = r.ensureDeployment(req, instance, r.smtpDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == JOB =======

	//== LMS Job ========
	result, err = r.ensureJob(req, instance, r.lmsJob(instance))
	if result != nil {
		return *result, err
	}

	lmsjobComplete := r.isLmsJobDone(instance)

	if !lmsjobComplete {
		// If lmsJob isn't complete, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(15)

		r.Log.Info(fmt.Sprintf("Lms Job isn't Complete, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	//== CMS Job ========
	result, err = r.ensureJob(req, instance, r.cmsJob(instance))
	if result != nil {
		return *result, err
	}

	cmsjobComplete := r.isCmsJobDone(instance)

	if !cmsjobComplete {
		// If cmsJob isn't complete, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(15)

		r.Log.Info(fmt.Sprintf("Cms Job isn't Complete, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	//== Forum Job ======
	result, err = r.ensureJob(req, instance, r.forumJob(instance))
	if result != nil {
		return *result, err
	}

	forumjobComplete := r.isForumJobDone(instance)

	if !forumjobComplete {
		// If forumJob isn't complete, requeue the reconcile
		// to run again after a delay
		delay := time.Second * time.Duration(15)

		r.Log.Info(fmt.Sprintf("Forum Job isn't Complete, waiting for %s", delay))
		return reconcile.Result{RequeueAfter: delay}, nil
	}

	// == INGRESS ==========

	result, err = r.ensureIngress(req, instance, r.ingress("web", instance))
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
