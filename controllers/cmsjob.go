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

const cmsJobPort = 8000

func cmsJobName(instance *cachev1.Openedx) string {
	return instance.Name + "-cmsjob"
}

func getCmsContainerEnv(cr *cachev1.Openedx) []corev1.EnvVar {
	env := make([]corev1.EnvVar, 0)

	env = append(env, corev1.EnvVar{
		Name:  "SERVICE_VARIANT",
		Value: "cms",
	})
	return env
}

// newJob returns a new Job instance.
func newCmsJob(instance *cachev1.Openedx) *batchv1.Job {
	labels := labels(instance, "cmsjob")
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmsJobName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
	}
}

func newCmsPodSpec(cr *cachev1.Openedx) corev1.PodSpec {
	pod := corev1.PodSpec{}

	pod.Containers = []corev1.Container{{
		Args: []string{
			"./manage.py",
			"cms",
			"migrate",
		},
		Env:             getCmsContainerEnv(cr),
		Image:           "docker.io/overhangio/openedx:11.2.1",
		ImagePullPolicy: corev1.PullAlways,
		Name:            "cms",
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

func newCmsPodTemplateSpec(cr *cachev1.Openedx) corev1.PodTemplateSpec {
	labels := labels(cr, "job")

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmsJobName(cr),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: newCmsPodSpec(cr),
	}
}

func (r *OpenedxReconciler) cmsJob(instance *cachev1.Openedx) *batchv1.Job {
	job := newCmsJob(instance)
	job.Spec.Template = newCmsPodTemplateSpec(instance)

	controllerutil.SetControllerReference(instance, job, r.Scheme)
	return job
}

func (r *OpenedxReconciler) isCmsJobDone(instance *cachev1.Openedx) bool {

	job := &batchv1.Job{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      cmsJobName(instance),
		Namespace: instance.Namespace,
	}, job)

	if err != nil {
		log.Error(err, "cmsjob not found")
		return false
	}

	if job.Status.Succeeded > 0 {
		return true
	}

	return false // Job not complete.

}
