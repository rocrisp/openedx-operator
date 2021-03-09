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

const redisImage = "docker.io/redis:6.0.9"
const redisPort = 6379

func redisDeploymentName(instance *cachev1.Openedx) string {
	return instance.Name + "-redis"
}

func redisServiceName(instance *cachev1.Openedx) string {
	return "redis"
}

func (r *OpenedxReconciler) redisDeployment(instance *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(instance, "redis")
	size := instance.Spec.Size

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redisDeploymentName(instance),
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
					Volumes: []corev1.Volume{
						{
							Name: "data",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "redis",
								},
							},
						}, {
							Name: "redis-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "redis-config",
									},
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Args: []string{
							"redis-server",
							"/openedx/redis/config/redis.conf",
						},
						WorkingDir: "/openedx/redis/data",
						Image:      redisImage,
						Name:       redisServiceName(instance),
						Ports: []corev1.ContainerPort{{
							ContainerPort: redisPort,
							Name:          "redis",
						}},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "redis-config",
								MountPath: "/openedx/redis/config/",
							},
							{
								Name:      "data",
								MountPath: "/openedx/redis/data",
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

func (r *OpenedxReconciler) redisService(instance *cachev1.Openedx) *corev1.Service {
	labels := labels(instance, "redis")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      redisServiceName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       redisPort,
				TargetPort: intstr.FromInt(redisPort),
			}},
		},
	}

	controllerutil.SetControllerReference(instance, service, r.Scheme)
	return service
}

// Returns whether or not the Memcached deployment is running
func (r *OpenedxReconciler) isRedisdUp(instance *cachev1.Openedx) bool {

	deployment := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      redisDeploymentName(instance),
		Namespace: instance.Namespace,
	}, deployment)

	if err != nil {
		log.Error(err, "Deployment Redis not found")
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
