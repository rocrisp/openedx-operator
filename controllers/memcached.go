package controllers

import (
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const memcachedPort = 11211

func memcachedDeploymentName() string {
	return "memcached"
}

func memcachedServiceName() string {
	return "memcached"
}

func (r *OpenedxReconciler) memcachedDeployment(instance *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(instance, "memcached")
	size := instance.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      memcachedDeploymentName(),
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "docker.io/memcached:1.4.38",
						Name:  "memcached-server",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 27017,
							Name:          "memcached",
						}},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(instance, dep, r.Scheme)
	return dep
}
