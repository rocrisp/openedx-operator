package controllers

import (
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const smtpPort = 25
const smtpImage = "docker.io/namshi/smtp:latest"

func smtpDeploymentName(instance *cachev1.Openedx) string {
	return instance.Name + "-smtp"
}

func smtpServiceName(instance *cachev1.Openedx) string {
	return instance.Name + "-smtp-service"
}

func (r *OpenedxReconciler) smtpDeployment(instance *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(instance, "smtp")
	size := instance.Spec.Size

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      smtpDeploymentName(instance),
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
						Image: smtpImage,
						Name:  "smtp-server",
						Ports: []corev1.ContainerPort{{
							ContainerPort: smtpPort,
							Name:          "smtp",
						}},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(instance, deployment, r.Scheme)
	return deployment
}

func (r *OpenedxReconciler) smtpService(instance *cachev1.Openedx) *corev1.Service {
	labels := labels(instance, "smtp")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      smtpServiceName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       25,
				TargetPort: intstr.FromInt(smtpPort),
				NodePort:   0,
			}},
		},
	}

	controllerutil.SetControllerReference(instance, service, r.Scheme)
	return service
}
