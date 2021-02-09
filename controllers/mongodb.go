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

const mongodbPort = 27017

func mongodbDeploymentName(instance *cachev1.Openedx) string {
	return instance.Name + "-mongodb"
}

func mongodbServiceName(instance *cachev1.Openedx) string {
	return instance.Name + "-mongodb-service"
}

func (r *OpenedxReconciler) mongodbDeployment(instance *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(instance, "mongodb")
	size := instance.Spec.Size

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mongodbDeploymentName(instance),
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
						Args: []string{
							"mongod",
							"--smallfiles",
							"--nojournal",
							"--storageEngine",
							"wiredTiger",
						},
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

	controllerutil.SetControllerReference(instance, deployment, r.Scheme)
	return deployment
}

func (r *OpenedxReconciler) mongodbService(instance *cachev1.Openedx) *corev1.Service {
	labels := labels(instance, "mongodb")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mongodbServiceName(instance),
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
				NodePort:   30060,
			}},
		},
	}

	controllerutil.SetControllerReference(instance, service, r.Scheme)
	return service
}

// Returns whether or not the MySQL deployment is running
func (r *OpenedxReconciler) isMongodbUp(instance *cachev1.Openedx) bool {

	deployment := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      nginxDeploymentName(instance),
		Namespace: instance.Namespace,
	}, deployment)

	if err != nil {
		log.Error(err, "Deployment mongodb not found")
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
