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

const sqlImage = "mysql:5.7"
const sqlPort = 3306

func mysqlDeploymentName(instance *cachev1.Openedx) string {
	return instance.Name + "-mysql"
}

func mysqlServiceName(instance *cachev1.Openedx) string {
	return instance.Name + "-mysql-service"
}

func mysqlAuthName() string {
	return "mysql-auth"
}

func (r *OpenedxReconciler) mysqlAuthSecret(instance *cachev1.Openedx) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlAuthName(),
			Namespace: instance.Namespace,
		},
		Type: "Opaque",
		StringData: map[string]string{
			"username": "root",
			"password": "cakephp",
		},
	}
	controllerutil.SetControllerReference(instance, secret, r.Scheme)
	return secret
}

func (r *OpenedxReconciler) mysqlDeployment(instance *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(instance, "mysql")
	size := instance.Spec.Size

	// userSecret := &corev1.EnvVarSource{
	// 	SecretKeyRef: &corev1.SecretKeySelector{
	// 		LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName()},
	// 		Key:                  "username",
	// 	},
	// }

	// passwordSecret := &corev1.EnvVarSource{
	// 	SecretKeyRef: &corev1.SecretKeySelector{
	// 		LocalObjectReference: corev1.LocalObjectReference{Name: mysqlAuthName()},
	// 		Key:                  "password",
	// 	},
	// }

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlDeploymentName(instance),
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
						Name: "mysql-data",
						VolumeSource: corev1.VolumeSource{
							EmptyDir: &corev1.EmptyDirVolumeSource{},
						},
					}},
					Containers: []corev1.Container{{
						Image: sqlImage,
						Name:  "mysql-server",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 3306,
							Name:          "mysql",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "mysql-data",
							MountPath: "/var/lib/mysql",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_ROOT_PASSWORD",
								Value: "cakephp",
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

func (r *OpenedxReconciler) mysqlService(instance *cachev1.Openedx) *corev1.Service {
	labels := labels(instance, "mysql")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlServiceName(instance),
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Port:       3306,
				TargetPort: intstr.FromInt(sqlPort),
			}},
		},
	}

	controllerutil.SetControllerReference(instance, service, r.Scheme)
	return service
}

// Returns whether or not the MySQL deployment is running
func (r *OpenedxReconciler) isMysqlUp(instance *cachev1.Openedx) bool {
	deployment := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      mysqlDeploymentName(instance),
		Namespace: instance.Namespace,
	}, deployment)

	if err != nil {
		log.Error(err, "Deployment mysql not found")
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
