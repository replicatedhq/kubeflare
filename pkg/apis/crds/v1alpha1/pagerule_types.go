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

type AutoMinifyPageRule struct {
	HTML bool `json:"html"`
	CSS  bool `json:"css"`
	JS   bool `json:"js"`
}

type AlwaysUseHTTPSPageRule struct {
}

type ForwardingURLPageRule struct {
	StatusCode  int    `json:"statusCode"`
	RedirectURL string `json:"redirectUrl"`
}

type CacheLevel string
type Toggle string

const (
	BypassCacheLevel          CacheLevel = "bypass"
	BasicCacheLevel           CacheLevel = "basic"
	SimplifiedCacheLevel      CacheLevel = "simplified"
	AggressiveCacheLevel      CacheLevel = "aggressive"
	CacheEverythingCacheLevel CacheLevel = "cache_everything"
	OnToggle                  Toggle     = "on"
	OffToggle                 Toggle     = "off"
)

type CacheLevelPageRule struct {
	// +kubebuilder:validation:Enum=bypass;basic;simplified;aggressive;cache_everything
	Level CacheLevel `json:"level"`
}

type EdgeCacheTTLPageRule struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=31536000
	Value int `json:"value"`
}

type CacheKeyCookieField struct {
	CheckPresence []string `json:"checkPresence,omitempty"`
	Include       []string `json:"include,omitempty"`
}

type CacheKeyHeaderField struct {
	CheckPresence []string `json:"checkPresence,omitempty"`
	Include       []string `json:"include,omitempty"`
	Exclude       []string `json:"exclude,omitempty"`
}

type CacheKeyHostField struct {
	Resolved bool `json:"resolved"`
}

type CacheKeyQueryStringField struct {
	Ignore  bool     `json:"ignore,omitempty"`
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
}

type CacheKeyUserField struct {
	DeviceType bool `json:"deviceType,omitempty"`
	Geo        bool `json:"geo,omitempty"`
	Lang       bool `json:"lang,omitempty"`
}

type CacheKeyFields struct {
	Cookie      CacheKeyCookieField      `json:"cookie,omitempty"`
	Header      CacheKeyHeaderField      `json:"header,omitempty"`
	Host        CacheKeyHostField        `json:"host,omitempty"`
	QueryString CacheKeyQueryStringField `json:"queryString,omitempty"`
	User        CacheKeyUserField        `json:"user,omitempty"`
}

type ExplicitCacheControlPageRule struct {
	// +kubebuilder:validation:Enum=on;off
	Value Toggle `json:"value"`
}

type SortQueryStringsPageRule struct {
	// +kubebuilder:validation:Enum=on;off
	Value Toggle `json:"value"`
}

type Rule struct {
	RequestURL string `json:"requestUrl"`

	ForwardingURL        *ForwardingURLPageRule        `json:"forwardingUrl,omitempty"`
	AlwaysUseHTTPS       *AlwaysUseHTTPSPageRule       `json:"alwaysUseHttps,omitempty"`
	AutoMinify           *AutoMinifyPageRule           `json:"autoMinify,omitempty"`
	CacheLevel           *CacheLevelPageRule           `json:"cacheLevel,omitempty"`
	EdgeCacheTTL         *EdgeCacheTTLPageRule         `json:"edgeCacheTTL,omitempty"`
	CacheKeyFields       *CacheKeyFields               `json:"cacheKeyFields,omitempty"`
	ExplicitCacheControl *ExplicitCacheControlPageRule `json:"explicitCacheControl,omitempty"`
	SortQueryStrings     *SortQueryStringsPageRule     `json:"sortQueryStrings,omitempty"`

	Priority *int `json:"priority,omitempty"`
	Enabled  bool `json:"enabled,omitempty"`
}

// PageRuleSpec defines the desired state of PageRule
type PageRuleSpec struct {
	Zone string `json:"zone"`
	Rule *Rule  `json:"pageRule,omitempty"`
}

// PageRuleStatus defines the observed state of PageRule
// We are storing the requested priority here because the priority is different on cloudflare side
// and hence we cannot depend on the one from its API to detect changes to the spec
type PageRuleStatus struct {
	ID                  string `json:"id,omitempty"`
	LastAppliedPriority int    `json:"lastAppliedPriority,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PageRule is the Schema for the pagerules API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type PageRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PageRuleSpec   `json:"spec,omitempty"`
	Status PageRuleStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PageRuleList contains a list of PageRule
type PageRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PageRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PageRule{}, &PageRuleList{})
}
