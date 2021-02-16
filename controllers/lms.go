package controllers

import (
	routev1 "github.com/openshift/api/route/v1"

	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const lmsImage = "docker.io/overhangio/openedx:10.4.0"
const lmsPort = 8000

func lmsDeploymentName(lms *cachev1.Openedx) string {
	return lms.Name + "-lms"
}

func lmsServiceName(lms *cachev1.Openedx) string {
	return "lms"
}

func (r *OpenedxReconciler) lmsDeployment(instance *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(instance, "lms")
	size := instance.Spec.Size

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      lmsDeploymentName(instance),
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
					Containers: []corev1.Container{{
						Image: lmsImage,
						Name:  "lms",
						Ports: []corev1.ContainerPort{{
							ContainerPort: lmsPort,
							Name:          "lms",
						}},

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

	controllerutil.SetControllerReference(instance, deployment, r.Scheme)
	return deployment
}

func (r *OpenedxReconciler) lmsService(instance *cachev1.Openedx) *corev1.Service {
	labels := labels(instance, "lms")

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      lmsServiceName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "NodePort",
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       lmsPort,
				TargetPort: intstr.FromInt(lmsPort),
			}},
		},
	}

	controllerutil.SetControllerReference(instance, service, r.Scheme)
	return service
}

func (r *OpenedxReconciler) lmsRoute(instance *cachev1.Openedx) *routev1.Route {
	labels := labels(instance, "lms")

	route := &routev1.Route{
		ObjectMeta: metav1.ObjectMeta{
			Name:      lmsServiceName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: lmsServiceName(instance),
			},
		},
	}
	controllerutil.SetControllerReference(instance, route, r.Scheme)
	return route
}
