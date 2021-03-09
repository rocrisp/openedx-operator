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

const sqlImage = "docker.io/mysql:5.7.32"
const sqlPort = 3306

func mysqlDeploymentName(instance *cachev1.Openedx) string {
	return instance.Name + "-mysql"
}

func mysqlServiceName(instance *cachev1.Openedx) string {
	return "mysql"
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

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mysqlDeploymentName(instance),
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
									ClaimName: "mysql",
								},
							},
						}, {
							Name: "mysql-initdb",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "mysql-initdb-config",
									},
								},
							},
						},
					},

					Containers: []corev1.Container{{
						Args: []string{
							"mysqld",
							"--character-set-server=utf8",
							"--collation-server=utf8_general_ci",
							"--ignore-db-dir=lost+found",
						},
						Image: sqlImage,
						Name:  "mysql-server",
						Ports: []corev1.ContainerPort{{
							ContainerPort: sqlPort,
							Name:          "mysql",
						}},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "data",
								MountPath: "/var/lib/mysql",
							},
							{
								Name:      "mysql-initdb",
								MountPath: "/docker-entrypoint-initdb.d",
							},
						},
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_ROOT_PASSWORD",
								Value: "mQh8ZJz4",
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
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			Ports: []corev1.ServicePort{{
				Port:       sqlPort,
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
