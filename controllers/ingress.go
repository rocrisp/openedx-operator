package controllers

import (
	cachev1 "github.com/rocrisp/openedx-operator/api/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const ingressPort = 8000

// newIngress returns a new Ingress instance for the given Openedx.
func newIngress(name string, cr *cachev1.Openedx) *extv1beta1.Ingress {
	labels := labels(cr, "ingress")
	return &extv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
	}
}

// newIngressWithName returns a new Ingress with the given name and ArgoCD.
func (r *OpenedxReconciler) lmsIngress(name string, cr *cachev1.Openedx) *extv1beta1.Ingress {
	ingress := newIngress(name, cr)

	// Add rules
	ingress.Spec.Rules = []extv1beta1.IngressRule{
		{
			Hosts: []extv1beta1.Paths {
				Host:  "www.lms-openedx.apps.demo.coreostrain.me",
				IngressRuleValue: extv1beta1.IngressRuleValue{
					HTTP: &extv1beta1.HTTPIngressRuleValue{
						Paths: []extv1beta1.HTTPIngressPath{
							{
								Backend: extv1beta1.IngressBackend{
									ServiceName: "nginx",
									ServicePort: intstr.FromInt(ingressPort),
								},
							},
						},
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(cr, ingress, r.Scheme)
	return ingress
}

// newIngressWithName returns a new Ingress with the given name and ArgoCD.
func (r *OpenedxReconciler) previewIngress(name string, cr *cachev1.Openedx) *extv1beta1.Ingress {
	ingress := newIngress(name, cr)

	// Add rules
	ingress.Spec.Rules = []extv1beta1.IngressRule{
		{
			Host: "preview.www.lms-openedx.apps.demo.coreostrain.me",
			IngressRuleValue: extv1beta1.IngressRuleValue{
				HTTP: &extv1beta1.HTTPIngressRuleValue{
					Paths: []extv1beta1.HTTPIngressPath{
						{
							Backend: extv1beta1.IngressBackend{
								ServiceName: "nginx",
								ServicePort: intstr.FromInt(ingressPort),
							},
						},
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(cr, ingress, r.Scheme)
	return ingress
}

// newIngressWithName returns a new Ingress with the given name and ArgoCD.
func (r *OpenedxReconciler) cmsIngress(name string, cr *cachev1.Openedx) *extv1beta1.Ingress {
	ingress := newIngress(name, cr)

	// Add rules
	ingress.Spec.Rules = []extv1beta1.IngressRule{
		{
			Host: "studio.www.lms-openedx.apps.demo.coreostrain.me",
			IngressRuleValue: extv1beta1.IngressRuleValue{
				HTTP: &extv1beta1.HTTPIngressRuleValue{
					Paths: []extv1beta1.HTTPIngressPath{
						{
							Backend: extv1beta1.IngressBackend{
								ServiceName: "nginx",
								ServicePort: intstr.FromInt(ingressPort),
							},
						},
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(cr, ingress, r.Scheme)
	return ingress
}
