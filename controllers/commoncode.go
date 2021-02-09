package controllers

import (
	"context"

	routev1 "github.com/openshift/api/route/v1"
	"github.com/prometheus/common/log"
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
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

func (r *OpenedxReconciler) ensureConfigMap(request reconcile.Request,
	instance *cachev1.Openedx,
	cm *corev1.ConfigMap,
) (*reconcile.Result, error) {
	found := &corev1.ConfigMap{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      cm.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the configMap
		log.Info("Creating a new ConfigMap")
		log.Info("ConfigMap Namespace : ", cm.Namespace)
		log.Info("COnfigMap Name : ", cm.Name)
		err = r.Client.Create(context.TODO(), cm)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the ConfigMap not existing
		log.Error(err, "Failed to get ConfigMap")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *OpenedxReconciler) ensurePVC(request reconcile.Request,
	instance *cachev1.Openedx,
	pvc *corev1.PersistentVolumeClaim,
) (*reconcile.Result, error) {
	found := &corev1.PersistentVolumeClaim{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      pvc.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the pvc
		log.Info("Creating a new pvc")
		log.Info("pvc Namespace : ", pvc.Namespace)
		log.Info("pvc Name : ", pvc.Name)
		err = r.Client.Create(context.TODO(), pvc)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new pvc", "pvc Namespace", pvc.Namespace, "pvc Name", pvc.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the pvc not existing
		log.Error(err, "Failed to get pvc")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *OpenedxReconciler) ensureJob(request reconcile.Request,
	instance *cachev1.Openedx,
	j *batchv1.Job,
) (*reconcile.Result, error) {
	found := &batchv1.Job{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      j.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the configMap
		log.Info("Creating a new Job")
		log.Info("Job Namespace : ", j.Namespace)
		log.Info("Job Name : ", j.Name)
		err = r.Client.Create(context.TODO(), j)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Job", "Job.Namespace", j.Namespace, "Job.Name", j.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the ConfigMap not existing
		log.Error(err, "Failed to get Job")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

//See if route exist, create one if it doesn't already exist
func (r *OpenedxReconciler) ensureRoute(request reconcile.Request,
	instance *cachev1.Openedx,
	rs *routev1.Route,
) (*reconcile.Result, error) {
	found := &routev1.Route{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      rs.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the route
		log.Info("Creating a new Route", "Route.Namespace", rs.Namespace, "Route.Name", rs.Name)
		err = r.Client.Create(context.TODO(), rs)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Route", "Route.Namespace", rs.Namespace, "Route.Name", rs.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get Route")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *OpenedxReconciler) ensureIngress(request reconcile.Request,
	instance *cachev1.Openedx,
	ing *extv1beta1.Ingress,
) (*reconcile.Result, error) {
	found := &extv1beta1.Ingress{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      ing.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the configMap
		log.Info("Creating a new Ingress")
		log.Info("Ingress Namespace : ", ing.Namespace)
		log.Info("Ingress Name : ", ing.Name)
		err = r.Client.Create(context.TODO(), ing)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Ingress", "Ingress.Namespace", ing.Namespace, "Ingress.Name", ing.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the ConfigMap not existing
		log.Error(err, "Failed to get Ingress")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func labels(instance *cachev1.Openedx, app string) map[string]string {
	return map[string]string{
		"app":        "OpenedX",
		"instance":   instance.Name,
		"managed-by": "rose",
		"tier":       app,
		"part-of":    "openedx",
		"version":    "10.4.0",
	}
}
