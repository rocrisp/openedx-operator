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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
)

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

	// == LMS  ==========
	result, err = r.ensureDeployment(req, instance, r.lmsDeployment(instance))
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

	// == MONODB ========
	result, err = r.ensureDeployment(req, instance, r.mongodbDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == NGINX ========
	result, err = r.ensureDeployment(req, instance, r.nginxDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == MEMCACHED ========
	result, err = r.ensureDeployment(req, instance, r.memcachedDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == RABBITMQ ========
	result, err = r.ensureDeployment(req, instance, r.rabbitmqDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == SMTP ========
	result, err = r.ensureDeployment(req, instance, r.smtpDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == ELASTICSEARCH ========
	result, err = r.ensureDeployment(req, instance, r.elasticsearchDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == FORUM ========
	result, err = r.ensureDeployment(req, instance, r.forumDeployment(instance))
	if result != nil {
		return *result, err
	}

	// == Finish ==========
	// Everything went fine, don't requeue
	return ctrl.Result{}, nil
}

// comment

func (r *OpenedxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1.Openedx{}).
		Complete(r)
}
