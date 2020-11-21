package controllers

import (
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const mongodbPort = 27017

func mongodbDeploymentName() string {
	return "mongodb"
}

func mongodbServiceName() string {
	return "mongodb"
}

func (r *OpenedxReconciler) mongodbDeployment(instance *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(instance, "mongodb")
	size := instance.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mongodbDeploymentName(),
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
					Volumes: []corev1.Volume{{
						Name: "mongodb-data",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: "mongodb",
							},
						},
					}},
					Containers: []corev1.Container{{
						Image: "docker.io/mongo:3.6.18",
						Name:  "mongodb-server",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 27017,
							Name:          "mongodb",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "mongodb-data",
							MountPath: "/data/db",
						}},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(instance, dep, r.Scheme)
	return dep
}
