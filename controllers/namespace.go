package controllers

import (
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//newConfigMap returns a new ConfigMap instance for the given OpenedX.
func newNamespace(instance *cachev1.Openedx) *corev1.Namespace {
	labels := labels(instance, "namespace")

	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "openedx",
			Namespace: "openedx",
			Labels:    labels,
		},
	}
}

func (r *OpenedxReconciler) namespace(instance *cachev1.Openedx) *corev1.Namespace {
	cm := newNamespace(instance)
	return cm
}
