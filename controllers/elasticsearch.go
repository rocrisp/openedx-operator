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

func elasticsearchDeploymentName(elasticsearch *cachev1.Openedx) string {
	return elasticsearch.Name + "-elasticsearch"
}

func elasticsearchServiceName(elasticsearch *cachev1.Openedx) string {
	return elasticsearch.Name + "-elasticsearch-service"
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

func (r *OpenedxReconciler) elasticsearchDeployment(instance *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(instance, "elasticsearch")
	size := instance.Spec.Size

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      elasticsearchDeploymentName(instance),
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

	controllerutil.SetControllerReference(instance, deployment, r.Scheme)
	return deployment
}

func (r *OpenedxReconciler) elasticsearchService(instance *cachev1.Openedx) *corev1.Service {
	labels := labels(instance, "elasticsearch")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      elasticsearchServiceName(instance),
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Port:       9200,
				TargetPort: intstr.FromInt(sqlPort),
			}},
		},
	}

	controllerutil.SetControllerReference(instance, service, r.Scheme)
	return service
}

// Returns whether or not the elasticsearch deployment is running
func (r *OpenedxReconciler) iselasticsearchUp(instance *cachev1.Openedx) bool {

	deployment := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      elasticsearchDeploymentName(instance),
		Namespace: instance.Namespace,
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
