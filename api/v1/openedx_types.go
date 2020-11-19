/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OpenedxSpec defines the desired state of Openedx
type OpenedxSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Size int32 `json:"size"`
}

// OpenedxStatus defines the observed state of Openedx
type OpenedxStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	PodStatus string `json:"podstatus"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Openedx is the Schema for the openedxes API
type Openedx struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenedxSpec   `json:"spec,omitempty"`
	Status OpenedxStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OpenedxList contains a list of Openedx
type OpenedxList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Openedx `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Openedx{}, &OpenedxList{})
}
