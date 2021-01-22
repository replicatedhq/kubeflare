/*
Copyright 2019 Replicated, Inc.

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

type AccessPolicy struct {
	Decision   string   `json:"descision"`
	Name       string   `json:"name"`
	Include    []string `json:"include"` // TODO
	Precedence *int     `json:"precendence,omitempty"`
	Exclude    []string `json:"exclude,omitempty"` // TODO
	Require    []string `json:"require,omitempty"` // TODO
}

type CORSHeader struct {
	AllowedMethods   []string `json:"allowedMethods"`
	AllowedOrigins   []string `json:"allowedOrigins"`
	AllowedHeaders   []string `json:"allowedHeader"`
	AllowAllMethods  bool     `json:"allowAllMethods"`
	AllowAllOrigins  bool     `json:"allowAllOrigins"`
	AllowAllHeaders  bool     `json:"allowAllHeaders"`
	AllowCredentials bool     `json:"allowCredentials"`
	MaxAge           int      `json:"maxAge"`
}

// AccessApplicationSpec defines the desired state of AccessApplication
type AccessApplicationSpec struct {
	Zone                   string         `json:"zone"`
	Name                   string         `json:"name"`
	Domain                 string         `json:"domain"`
	SessionDuration        string         `json:"sessionDuration,omitempty"`
	AllowedIdPs            []string       `json:"allowedIdPs,omitempty"`
	AutoRedirectToIdentity *bool          `json:"autoRedirectToIndentiy,omitempty"`
	CORSHeaders            *CORSHeader    `json:"corsHeaders,omitempty"`
	AccessPolicies         []AccessPolicy `json:"accessPolicies,omitempty"`
}

// AccessApplicationStatus defines the observed state of AccessApplicationS
type AccessApplicationStatus struct {
	ApplicationID string `json:"applicationID,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DNSRecord is the Schema for the accessapplication API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type AccessApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AccessApplicationSpec   `json:"spec,omitempty"`
	Status AccessApplicationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AccessApplicationList contains a list of DNSRecord
type AccessApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DNSRecord `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AccessApplication{}, &AccessApplicationList{})
}
