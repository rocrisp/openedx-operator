package controllers

import (
	"context"
	"github.com/prometheus/common/log"
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const lmsJobPort = 8000

func lmsJobName(instance *cachev1.Openedx) string {
	return instance.Name + "-lmsjob"
}

// newJob returns a new Job instance.
func newLmsJob(instance *cachev1.Openedx) *batchv1.Job {
	labels := labels(instance, "lmsjob")
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      lmsJobName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
	}
}

func newLmsPodSpec(cr *cachev1.Openedx) corev1.PodSpec {
	pod := corev1.PodSpec{}

	pod.Containers = []corev1.Container{{
		Args: []string{
			"./manage.py",
			"lms",
			"migrate",
		},
		Image:           "docker.io/overhangio/openedx:11.2.0",
		ImagePullPolicy: corev1.PullAlways,
		Name:            "lms",
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
	}}
	pod.RestartPolicy = corev1.RestartPolicyOnFailure

	pod.Volumes = []corev1.Volume{
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
	}

	return pod
}

func newLmsPodTemplateSpec(cr *cachev1.Openedx) corev1.PodTemplateSpec {
	labels := labels(cr, "job")

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      lmsJobName(cr),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: newCmsPodSpec(cr),
	}
}

func (r *OpenedxReconciler) lmsJob(instance *cachev1.Openedx) *batchv1.Job {
	job := newLmsJob(instance)
	job.Spec.Template = newLmsPodTemplateSpec(instance)

	controllerutil.SetControllerReference(instance, job, r.Scheme)
	return job
}

func (r *OpenedxReconciler) isLmsJobDone(instance *cachev1.Openedx) bool {

	job := &batchv1.Job{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      lmsJobName(instance),
		Namespace: instance.Namespace,
	}, job)

	if err != nil {
		log.Error(err, "lmsjob not found")
		return false
	}

	if job.Status.Succeeded > 0 {
		return true
	}

	return false // Job not complete.

}
