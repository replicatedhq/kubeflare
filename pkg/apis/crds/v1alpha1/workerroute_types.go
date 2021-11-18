/*
Copyright 2021 The Kubernetes authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WorkerRouteSpec defines the desired state of WorkerRoute
type WorkerRouteSpec struct {
	Zone    string `json:"zone"`
	Pattern string `json:"pattern"`
	Script  string `json:"script,omitempty"`
}

// WorkerRouteStatus defines the observed state of WorkerRoute
type WorkerRouteStatus struct {
	ID        string `json:"id,omitempty"`
	LastError string `json:"lastError,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WorkerRoute is the Schema for the workerroutes API
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type WorkerRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkerRouteSpec   `json:"spec,omitempty"`
	Status WorkerRouteStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WorkerRouteList contains a list of WorkerRoute
// +kubebuilder:object:root=true
type WorkerRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WorkerRoute `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WorkerRoute{}, &WorkerRouteList{})
}
