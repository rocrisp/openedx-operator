package controllers

import (
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const caddyImage = "docker.io/caddy:2.2.1"
const caddyPort1 = 80
const caddyPort2 = 443

func caddyDeploymentName(instance *cachev1.Openedx) string {
	return instance.Name + "-caddy"
}

func caddyServiceName(instance *cachev1.Openedx) string {
	return "caddy"
}

func (r *OpenedxReconciler) caddyDeployment(instance *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(instance, "caddy")
	size := instance.Spec.Size

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      caddyDeploymentName(instance),
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
					Volumes: []corev1.Volume{
						{
							Name: "data",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "caddy",
								},
							},
						}, {
							Name: "caddyConfig",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "caddy-config",
									},
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Image: caddyImage,
						Name:  "caddy",
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: caddyPort1,
								Name:          "caddy1",
							},
							{
								ContainerPort: caddyPort2,
								Name:          "caddy2",
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "config",
								MountPath: "/etc/caddy/",
							},
							{
								Name:      "data",
								MountPath: "/data/",
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(instance, deployment, r.Scheme)
	return deployment
}

func (r *OpenedxReconciler) caddyService(instance *cachev1.Openedx) *corev1.Service {
	labels := labels(instance, "caddy")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      caddyServiceName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "LoadBalancer",
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Protocol: corev1.ProtocolTCP,
					Port:     caddyPort1,
				},
				{
					Name:     "https",
					Protocol: corev1.ProtocolTCP,
					Port:     caddyPort2,
				},
			},
		},
	}

	controllerutil.SetControllerReference(instance, service, r.Scheme)
	return service
}
