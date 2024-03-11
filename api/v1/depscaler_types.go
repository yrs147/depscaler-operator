/*
Copyright 2024.

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

const (
	SUCCESS = "Success"
	FAILED  = "Failed"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DepScalerSpec defines the desired state of DepScaler
type DepScalerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Maximum=23
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	Begin int `json:"begin"`

	// +kubebuilder:validation:Maximum=23
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	End int `json:"end"`

	// +kubebuilder:validation:Required
	Replicas    int32            `json:"replicas"`
	Deployments []NamespacedName `json:"deployments"`
}

type NamespacedName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// DepScalerStatus defines the observed state of DepScaler
type DepScalerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status string `json:"status,,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DepScaler is the Schema for the depscalers API
type DepScaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DepScalerSpec   `json:"spec,omitempty"`
	Status DepScalerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DepScalerList contains a list of DepScaler
type DepScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DepScaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DepScaler{}, &DepScalerList{})
}
