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

const demoJobPort = 8000

func demoJobName(instance *cachev1.Openedx) string {
	return instance.Name + "-demojob"
}

func getDemoContainerEnv(cr *cachev1.Openedx) []corev1.EnvVar {
	env := make([]corev1.EnvVar, 0)

	env = append(env, corev1.EnvVar{
		Name:  "SERVICE_VARIANT",
		Value: "cms",
	})
	return env
}

// newJob returns a new Job instance.
func newDemoJob(instance *cachev1.Openedx) *batchv1.Job {
	labels := labels(instance, "cmsjob")
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      demoJobName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
	}
}

func newDemoPodSpec(cr *cachev1.Openedx) corev1.PodSpec {
	pod := corev1.PodSpec{}

	cmd := make([]string, 0)
	cmd = append(cmd, "/bin/sh")
	cmd = append(cmd, "-c")
	cmd = append(cmd, "git clone https://github.com/edx/edx-demo-course --branch open-release/koa.2 --depth 1 ../edx-demo-course")

	pod.InitContainers = []corev1.Container{
		{
			Command:         cmd,
			Env:             getDemoContainerEnv(cr),
			Image:           "docker.io/overhangio/openedx:11.2.3",
			ImagePullPolicy: corev1.PullAlways,
			Name:            "init-demo-part1",
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
		},
		{
			Args: []string{
				"./manage.py",
				"cms",
				"import",
				"../data",
				"../edx-demo-course",
			},
			Env:             getDemoContainerEnv(cr),
			Image:           "docker.io/overhangio/openedx:11.2.3",
			ImagePullPolicy: corev1.PullAlways,
			Name:            "init-demo-part2",
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
		},
	}

	pod.Containers = []corev1.Container{{
		Args: []string{
			"./manage.py",
			"cms",
			"reindex_course",
			"--all",
			"--setup",
		},
		Env:             getDemoContainerEnv(cr),
		Image:           "docker.io/overhangio/openedx:11.2.3",
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

func newDemoPodTemplateSpec(cr *cachev1.Openedx) corev1.PodTemplateSpec {
	labels := labels(cr, "job")

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      demoJobName(cr),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: newDemoPodSpec(cr),
	}
}

func (r *OpenedxReconciler) demoJob(instance *cachev1.Openedx) *batchv1.Job {
	job := newDemoJob(instance)
	job.Spec.Template = newDemoPodTemplateSpec(instance)

	controllerutil.SetControllerReference(instance, job, r.Scheme)
	return job
}

func (r *OpenedxReconciler) isDemoJobDone(instance *cachev1.Openedx) bool {

	job := &batchv1.Job{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      demoJobName(instance),
		Namespace: instance.Namespace,
	}, job)

	if err != nil {
		log.Error(err, "demojob not found")
		return false
	}

	if job.Status.Succeeded > 0 {
		return true
	}

	return false // Job not complete.

}
