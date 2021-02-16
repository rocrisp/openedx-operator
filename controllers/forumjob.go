package controllers

import (
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func forumjobName(instance *cachev1.Openedx) string {
	return instance.Name + "-forumjob"
}

// getArgoExportCommand will return the command for the ArgoCD export process.
func getArgoExportCommand(cr *cachev1.Openedx) []string {
	cmd := make([]string, 0)
	cmd = append(cmd, "sh")
	cmd = append(cmd, "-e")
	cmd = append(cmd, "-c")
	cmd = append(cmd, "bundle exec rake search:initialize")
	cmd = append(cmd, "bundle exec rake search:rebuild_index")
	return cmd
}

func getArgoExportContainerEnv(cr *cachev1.Openedx) []corev1.EnvVar {
	env := make([]corev1.EnvVar, 0)

	env = append(env, corev1.EnvVar{
		Name:  "MONGOHQ_URL",
		Value: "mongodb://mogodb:27017/cs_comments_service",
	})

	env = append(env, corev1.EnvVar{
		Name:  "SEARCH_SERVER",
		Value: "http://elasticsearch:9200",
	})

	env = append(env, corev1.EnvVar{
		Name:  "MONGODB_AUTH",
		Value: "",
	})

	env = append(env, corev1.EnvVar{
		Name:  "MONGODB_HOST",
		Value: "mongodb",
	})

	env = append(env, corev1.EnvVar{
		Name:  "MONGODB_PORT",
		Value: "27017",
	})

	return env
}

// newJob returns a new Job instance.
func newJob(instance *cachev1.Openedx) *batchv1.Job {
	labels := labels(instance, "job")
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      forumjobName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
	}
}

func newPodSpec(cr *cachev1.Openedx) corev1.PodSpec {
	pod := corev1.PodSpec{}

	pod.Containers = []corev1.Container{{
		Command:         getArgoExportCommand(cr),
		Env:             getArgoExportContainerEnv(cr),
		Image:           "docker.io/overhangio/openedx-forum:10.4.0",
		ImagePullPolicy: corev1.PullAlways,
		Name:            "forum",
	}}

	pod.RestartPolicy = corev1.RestartPolicyOnFailure
	return pod
}

func newPodTemplateSpec(cr *cachev1.Openedx) corev1.PodTemplateSpec {
	labels := labels(cr, "job")

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      forumjobName(cr),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: newPodSpec(cr),
	}
}

func (r *OpenedxReconciler) job1(instance *cachev1.Openedx) *batchv1.Job {
	job := newJob(instance)
	job.Spec.Template = newPodTemplateSpec(instance)

	controllerutil.SetControllerReference(instance, job, r.Scheme)
	return job
}
