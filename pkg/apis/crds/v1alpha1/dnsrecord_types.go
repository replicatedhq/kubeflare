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

type Record struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	TTL      *int   `json:"ttl,omitempty"`
	Priority *int   `json:"priority,omitempty"`
	Proxied  *bool  `json:"proxied,omitempty"`
}

// DNSRecordSpec defines the desired state of DNSRecord
type DNSRecordSpec struct {
	Zone    string    `json:"zone"`
	Record  *Record   `json:"record,omitempty"`
	Records []*Record `json:"records,omitempty"`
}

// DNSRecordStatus defines the observed state of DNSRecord
type DNSRecordStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DNSRecord is the Schema for the dnsrecords API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type DNSRecord struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DNSRecordSpec   `json:"spec,omitempty"`
	Status DNSRecordStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DNSRecordList contains a list of DNSRecord
type DNSRecordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DNSRecord `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DNSRecord{}, &DNSRecordList{})
}
