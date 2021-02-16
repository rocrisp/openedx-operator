package controllers

import (
	"context"

	"github.com/prometheus/common/log"
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const memcachedImage = "docker.io/memcached:1.4.38"
const memcachedPort = 11211

func memcachedDeploymentName(instance *cachev1.Openedx) string {
	return instance.Name + "-memcached"
}

func memcachedServiceName(instance *cachev1.Openedx) string {
	return "memcached"
}

func (r *OpenedxReconciler) memcachedDeployment(instance *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(instance, "memcached")
	size := instance.Spec.Size

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      memcachedDeploymentName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
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
						Image: memcachedImage,
						Name:  memcachedServiceName(instance),
						Ports: []corev1.ContainerPort{{
							ContainerPort: memcachedPort,
							Name:          "memcached",
						}},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(instance, deployment, r.Scheme)
	return deployment
}

func (r *OpenedxReconciler) memcachedService(instance *cachev1.Openedx) *corev1.Service {
	labels := labels(instance, "memcached")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      memcachedServiceName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       memcachedPort,
				TargetPort: intstr.FromInt(memcachedPort),
			}},
		},
	}

	controllerutil.SetControllerReference(instance, service, r.Scheme)
	return service
}

// Returns whether or not the MySQL deployment is running
func (r *OpenedxReconciler) isMemcachedUp(instance *cachev1.Openedx) bool {

	deployment := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      memcachedDeploymentName(instance),
		Namespace: instance.Namespace,
	}, deployment)

	if err != nil {
		log.Error(err, "Deployment memcached not found")
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
