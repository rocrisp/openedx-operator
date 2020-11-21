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

const rabbitmqPort = 5672
const rabbitmqImage = "docker.io/rabbitmq:3.6.10-management-alpine"

func rabbitmqDeploymentName() string {
	return "rabbitmq"
}

func rabbitmqServiceName() string {
	return "rabbitmq"
}

func (r *OpenedxReconciler) rabbitmqDeployment(d *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(d, "rabbitmq")
	size := d.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rabbitmqDeploymentName(),
			Namespace: d.Namespace,
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
						Name: "rabbitmq-data",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: "rabbitmq",
							},
						},
					}},
					Containers: []corev1.Container{{
						Image: rabbitmqImage,
						Name:  "rabbitmq-server",
						Ports: []corev1.ContainerPort{{
							ContainerPort: rabbitmqPort,
							Name:          "rabbitmq",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "rabbitmq-data",
							MountPath: "/var/lib/rabbitmq",
						}},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(d, dep, r.Scheme)
	return dep
}

func (r *OpenedxReconciler) rabbitmqService(d *cachev1.Openedx) *corev1.Service {
	labels := labels(d, "rabbitmq")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rabbitmqServiceName(),
			Namespace: d.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Port:       3306,
				TargetPort: intstr.FromInt(sqlPort),
			}},
		},
	}

	controllerutil.SetControllerReference(d, s, r.Scheme)
	return s
}

// Returns whether or not the rabbitmq deployment is running
func (r *OpenedxReconciler) israbbitmqUp(d *cachev1.Openedx) bool {
	deployment := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      rabbitmqDeploymentName(),
		Namespace: d.Namespace,
	}, deployment)

	if err != nil {
		log.Error(err, "Deployment rabbitmq not found")
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
