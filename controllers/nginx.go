package controllers

import (
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const nginxImage = "docker.io/nginx:1.13"
const nginxPort = 8000

func nginxDeploymentName(nginx *cachev1.Openedx) string {
	return nginx.Name + "-nginx"
}

func (r *OpenedxReconciler) nginxDeployment(nginx *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(nginx, "nginx")
	size := nginx.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nginxDeploymentName(nginx),
			Namespace: nginx.Namespace,
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
								ContainerPort: 80,
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

	controllerutil.SetControllerReference(nginx, dep, r.Scheme)
	return dep
}
