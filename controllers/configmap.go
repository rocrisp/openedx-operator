package controllers

import (
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConfigMapName(instance *cachev1.Openedx) string {
	return instance.Name + "-configmap"
}

// newConfigMap returns a new ConfigMap instance for the given OpenedX.
func newConfigMap(instance *cachev1.Openedx) *corev1.ConfigMap {
	labels := labels(instance, "configmap")

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ConfigMapName(instance),
			Namespace: instance.Namespace,
			Labels:    labels,
		},
	}
}

func (r *OpenedxReconciler) ConfigMap(instance *cachev1.Openedx) *corev1.ConfigMap {
	cm := newConfigMap(instance)
	cm.ObjectMeta.Name = "openedx-configxx"

	lbls := cm.ObjectMeta.Labels
	lbls["newname"] = "OpenedXConfigmap"
	cm.ObjectMeta.Labels = lbls

	return cm
}
