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

const elasticsearchPort = 9200
const elasticsearchImage = "docker.io/elasticsearch:1.5.2"

func elasticsearchDeploymentName() string {
	return "elasticsearch"
}

func elasticsearchServiceName() string {
	return "elasticsearch"
}

func elasticsearchAuthName() string {
	return "elasticsearch-auth"
}

func (r *OpenedxReconciler) elasticsearchAuthSecret(d *cachev1.Openedx) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      elasticsearchAuthName(),
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

func (r *OpenedxReconciler) elasticsearchDeployment(d *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(d, "elasticsearch")
	size := d.Spec.Size

	// userSecret := &corev1.EnvVarSource{
	// 	SecretKeyRef: &corev1.SecretKeySelector{
	// 		LocalObjectReference: corev1.LocalObjectReference{Name: elasticsearchAuthName()},
	// 		Key:                  "username",
	// 	},
	// }

	// passwordSecret := &corev1.EnvVarSource{
	// 	SecretKeyRef: &corev1.SecretKeySelector{
	// 		LocalObjectReference: corev1.LocalObjectReference{Name: elasticsearchAuthName()},
	// 		Key:                  "password",
	// 	},
	// }

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      elasticsearchDeploymentName(),
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
						Name: "elasticsearch-data",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: "elasticsearch",
							},
						},
					}},
					Containers: []corev1.Container{{
						Image: elasticsearchImage,
						Name:  "elasticsearch",
						Ports: []corev1.ContainerPort{{
							ContainerPort: elasticsearchPort,
							Name:          "elasticsearch",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "elasticsearch-data",
							MountPath: "/usr/share/elasticsearch/data",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "ES_JAVA_OPTS",
								Value: "-Xms1g -Xmx1g",
							},
							{
								Name:  "cluster.name",
								Value: "openedx",
							},
							{
								Name:  "bootstrap.memory_lock",
								Value: "true",
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

func (r *OpenedxReconciler) elasticsearchService(d *cachev1.Openedx) *corev1.Service {
	labels := labels(d, "elasticsearch")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      elasticsearchServiceName(),
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

// Returns whether or not the elasticsearch deployment is running
func (r *OpenedxReconciler) iselasticsearchUp(d *cachev1.Openedx) bool {
	deployment := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      elasticsearchDeploymentName(),
		Namespace: d.Namespace,
	}, deployment)

	if err != nil {
		log.Error(err, "Deployment elasticsearch not found")
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
