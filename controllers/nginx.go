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

const nginxImage = "docker.io/nginx:1.13"
const nginxPort = 80

func nginxDeploymentName(instance *cachev1.Openedx) string {
	return instance.Name + "-nginx"
}

func nginxServiceName(instance *cachev1.Openedx) string {
	return instance.Name + "-nginx-service"
}

func (r *OpenedxReconciler) nginxDeployment(instance *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(instance, "nginx")
	size := instance.Spec.Size

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nginxDeploymentName(instance),
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
						Image: nginxImage,
						Name:  "nginx",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: nginxPort,
								Name:          "nginx1",
							},
							{
								ContainerPort: 443,
								Name:          "nginx2",
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "config",
								MountPath: "/etc/nginx/conf.d/",
							},
						},
					}},
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "openedx-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	controllerutil.SetControllerReference(instance, deployment, r.Scheme)
	return deployment
}

func (r *OpenedxReconciler) nginxService(instance *cachev1.Openedx) *corev1.Service {
	labels := labels(instance, "nginx")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nginxServiceName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       80,
				TargetPort: intstr.FromInt(nginxPort),
				NodePort:   30080,
			}},
		},
	}

	controllerutil.SetControllerReference(instance, service, r.Scheme)
	return service
}

// Returns whether or not the MySQL deployment is running
func (r *OpenedxReconciler) isNginxUp(instance *cachev1.Openedx) bool {

	deployment := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      nginxDeploymentName(instance),
		Namespace: instance.Namespace,
	}, deployment)

	if err != nil {
		log.Error(err, "Deployment nginx not found")
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
