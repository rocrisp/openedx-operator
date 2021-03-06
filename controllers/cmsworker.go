package controllers

import (
	"context"
	"github.com/prometheus/common/log"
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const cmsworkerImage = "docker.io/overhangio/openedx:11.2.3"
const cmsworkerPort = 8000

func cmsworkerDeploymentName(cr *cachev1.Openedx) string {
	return cr.Name + "-cmsworker"
}

func (r *OpenedxReconciler) cmsworkerDeployment(cr *cachev1.Openedx) *appsv1.Deployment {
	labels := labels(cr, "cmsworker")
	size := cr.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmsworkerDeploymentName(cr),
			Namespace: cr.Namespace,
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
						Args: []string{
							"celery",
							"worker",
							"--app=cms.celery",
							"--loglevel=info",
							"--hostname=edx.cms.core.default.%%h",
							"--maxtasksperchild",
							"100",
							"--exclude-queues=edx.lms.core.default",
						},
						Image: cmsworkerImage,
						Name:  "cms-worker",
						Ports: []corev1.ContainerPort{{
							ContainerPort: cmsPort,
							Name:          "cmsworker",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "SERVICE_VARIANT",
								Value: "cms",
							},
							{
								Name:  "C_FORCE_ROOT",
								Value: "1",
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

	controllerutil.SetControllerReference(cr, dep, r.Scheme)
	return dep
}

//Returns whether or not the cms deployment is running
func (r *OpenedxReconciler) isCmsworkerUp(cr *cachev1.Openedx) bool {

	deployment := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      cmsworkerDeploymentName(cr),
		Namespace: cr.Namespace,
	}, deployment)

	if err != nil {
		log.Error(err, "Deployment cmsworker not found")
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
