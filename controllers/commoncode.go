package controllers

import (
	"context"

	"github.com/prometheus/common/log"
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *OpenedxReconciler) ensureDeployment(request reconcile.Request,
	instance *cachev1.Openedx,
	dep *appsv1.Deployment,
) (*reconcile.Result, error) {

	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		log.Info("Creating a new Deployment")
		log.Info("Deployment Namespace : ", dep.Namespace)
		log.Info("Deployment Name : ", dep.Name)

		err = r.Client.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			log.Info("Failed to create new Deployment")
			log.Info("Deployment Namespace : ", dep.Namespace)
			log.Info("Deployment Name : ", dep.Name)
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		log.Error(err, "Failed to get Deployment")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *OpenedxReconciler) ensureService(request reconcile.Request,
	instance *cachev1.Openedx,
	s *corev1.Service,
) (*reconcile.Result, error) {
	found := &corev1.Service{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new Service")
		log.Info("Service Namespace : ", s.Namespace)
		log.Info("Service Name : ", s.Name)
		err = r.Client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get Service")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func labels(instance *cachev1.Openedx, app string) map[string]string {
	return map[string]string{
		"app":        "openedx",
		"instance":   instance.Name,
		"managed-by": "rose",
		"name":       app,
		"part-of":    "openedx",
		"version":    "10.4.0",
	}
}
