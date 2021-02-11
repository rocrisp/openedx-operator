package controllers

import (
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const cmsImage = "docker.io/overhangio/openedx:10.4.0"
const cmsPort = 8000

func cmsDeploymentName(instance *cachev1.Openedx) string {
	return instance.Name + "-cms"
}
func cmsServiceName(instance *cachev1.Openedx) string {
	return instance.Name + "-cms-service"
}

func (r *OpenedxReconciler) cmsDeployment(cms *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(cms, "cms")
	size := cms.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmsDeploymentName(cms),
			Namespace: cms.Namespace,
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
						Image: cmsImage,
						Name:  "cms",
						Ports: []corev1.ContainerPort{{
							ContainerPort: cmsPort,
							Name:          "cms",
						}},

						Env: []corev1.EnvVar{
							{
								Name:  "SERVICE_VARIANT",
								Value: "cms",
							},
						},

						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "settings-lms",
								MountPath: "/openedx/edx-platform/lms/envs/tutor/",
							},
							{
								Name:      "settings-cms",
								MountPath: "/openedx/edx-platform/cms/envs/tutor/",
							},
							{
								Name:      "config",
								MountPath: "/openedx/config",
							},
						},
					}},

					Volumes: []corev1.Volume{
						{
							Name: "settings-lms",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "openedx-settings-lms",
									},
								},
							},
						}, {
							Name: "settings-cms",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "openedx-settings-cms",
									},
								},
							},
						}, {
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "openedx-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	controllerutil.SetControllerReference(cms, dep, r.Scheme)
	return dep
}

func (r *OpenedxReconciler) cmsService(instance *cachev1.Openedx) *corev1.Service {
	labels := labels(instance, "cms")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmsServiceName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       cmsPort,
				TargetPort: intstr.FromInt(cmsPort),
			}},
		},
	}

	controllerutil.SetControllerReference(instance, service, r.Scheme)
	return service
}
