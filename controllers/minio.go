package controllers

import (
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const minioImage = "docker.io/overhangio/minio:10.1.3"
const minioPort = 9000

func minioDeploymentName(instance *cachev1.Openedx) string {
	return instance.Name + "-minio"
}

func minioServiceName(instance *cachev1.Openedx) string {
	return instance.Name + "-minio-service"
}

func (r *OpenedxReconciler) minioDeployment(instance *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(instance, "minio")
	size := instance.Spec.Size

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      minioDeploymentName(instance),
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
						Name: "data",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: "minio",
							},
						},
					}},
					Containers: []corev1.Container{{
						Args: []string{
							"server",
							"--address",
							":9000",
							"/data",
						},
						Image: minioImage,
						Name:  "minio",
						Ports: []corev1.ContainerPort{{
							ContainerPort: minioPort,
							Name:          "minio",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "MINIO_ACCESS_KEY",
								Value: "openedx",
							},
							{
								Name:  "MINIO_SECRET_KEY",
								Value: "1ARGphpXC0Xwpv59qhrk8PYm",
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "data",
								MountPath: "/data",
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

func (r *OpenedxReconciler) minioService(instance *cachev1.Openedx) *corev1.Service {
	labels := labels(instance, "minio")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      minioServiceName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       minioPort,
				TargetPort: intstr.FromInt(minioPort),
				NodePort:   0,
			}},
		},
	}

	controllerutil.SetControllerReference(instance, service, r.Scheme)
	return service
}
