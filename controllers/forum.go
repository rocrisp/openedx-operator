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

const forumPort = 4567
const forumImage = "docker.io/overhangio/openedx-forum:10.4.0"

func forumDeploymentName(instance *cachev1.Openedx) string {
	return instance.Name + "-forum"
}

func forumServiceName(instance *cachev1.Openedx) string {
	return instance.Name + "-forum-service"
}

func forumAuthName() string {
	return "forum-auth"
}

func (r *OpenedxReconciler) forumAuthSecret(d *cachev1.Openedx) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      forumAuthName(),
			Namespace: d.Namespace,
		},
		Type: "Opaque",
		StringData: map[string]string{
			"username": "root",
			"password": "cakephp",
		},
	}
	controllerutil.SetControllerReference(d, secret, r.Scheme)
	return secret
}

func (r *OpenedxReconciler) forumDeployment(d *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(d, "forum")
	size := d.Spec.Size

	// userSecret := &corev1.EnvVarSource{
	// 	SecretKeyRef: &corev1.SecretKeySelector{
	// 		LocalObjectReference: corev1.LocalObjectReference{Name: forumAuthName()},
	// 		Key:                  "username",
	// 	},
	// }

	// passwordSecret := &corev1.EnvVarSource{
	// 	SecretKeyRef: &corev1.SecretKeySelector{
	// 		LocalObjectReference: corev1.LocalObjectReference{Name: forumAuthName()},
	// 		Key:                  "password",
	// 	},
	// }

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      forumDeploymentName(d),
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

					Containers: []corev1.Container{{
						Image: forumImage,
						Name:  "forum",
						Ports: []corev1.ContainerPort{{
							ContainerPort: forumPort,
							Name:          "forum",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "SEARCH_SERVER",
								Value: "http://elasticsearch:9200",
							},
							{
								Name:  "MONGODB_AUTH",
								Value: "",
							},
							{
								Name:  "MONGODB_HOST",
								Value: "mongodb",
							},
							{
								Name:  "MONGODB_PORT",
								Value: "27017",
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(d, dep, r.Scheme)
	return dep
}

func (r *OpenedxReconciler) forumService(instance *cachev1.Openedx) *corev1.Service {
	labels := labels(instance, "forum")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      forumServiceName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       forumPort,
				TargetPort: intstr.FromInt(forumPort),
				NodePort:   30040,
			}},
		},
	}

	controllerutil.SetControllerReference(instance, service, r.Scheme)
	return service
}

// Returns whether or not the forum deployment is running
func (r *OpenedxReconciler) isforumUp(d *cachev1.Openedx) bool {
	deployment := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      forumDeploymentName(d),
		Namespace: d.Namespace,
	}, deployment)

	if err != nil {
		log.Error(err, "Deployment forum not found")
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
