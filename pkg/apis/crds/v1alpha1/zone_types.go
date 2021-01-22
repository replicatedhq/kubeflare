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

type SecurityHeader struct {
	Enabled           *bool `json:"enabled,omitempty"`
	MaxAge            *int  `json:"maxAge,omitempty"`
	IncludeSubdomains *bool `json:"includeSubdomains,omitempty"`
	NoSniff           *bool `json:"noSniff,omitempty"`
}

type MobileRedirect struct {
	Status          *bool   `json:"status,omi2tempty"`
	MobileSubdomain *string `json:"mobileSubdomain,omitempty"`
	StripURI        *bool   `json:"stripURI,omitempty"`
}

type MinifySetting struct {
	CSS  *bool `json:"css,omitempty"`
	HTML *bool `json:"html,omitempty"`
	JS   *bool `json:"js,omitempty"`
}

type ZoneSettings struct {
	AdvancedDDOS            *bool           `json:"advancedDDOS,omitempty"`
	AlwaysOnline            *bool           `json:"alwaysOnline,omitempty"`
	AlwaysUseHTTPS          *bool           `json:"alwaysUseHttps,omitempty"`
	OpportunisticOnion      *bool           `json:"opportunisticOnion,omitempty"`
	AutomaticHTTPSRewrites  *bool           `json:"automaticHTTPSRewrites,omitempty"`
	BrowserCacheTTL         *int            `json:"browserCacheTTL,omitempty"`
	BrowserCheck            *bool           `json:"browserCheck,omitempty"`
	CacheLevel              *string         `json:"cacheLevel,omitempty"`
	ChallengeTTL            *int            `json:"challengeTTL,omitempty"`
	DevelopmentMode         *bool           `json:"developmentMode,omitempty"`
	EmailObfuscation        *bool           `json:"emailObfuscation,omitempty"`
	HotlinkProtection       *bool           `json:"hotlinkProtection,omitempty"`
	IPGeolocation           *bool           `json:"ipGeolocation,omitempty"`
	IPV6                    *bool           `json:"ipv6,omitempty"`
	Minify                  *MinifySetting  `json:"minify,omitempty"`
	MobileRedirect          *MobileRedirect `json:"mobileRedirect,omitempty"`
	Mirage                  *bool           `json:"mirage,omitempty"`
	OriginErrorPagePassThru *bool           `json:"originErrorPagePassThru,omitempty"`
	OpportunisticEncryption *bool           `json:"opportunisticEncryption,omitempty"`
	Polish                  *bool           `json:"polish,omitempty"`
	WebP                    *bool           `json:"webp,omitempty"`
	Brotli                  *bool           `json:"brotli,omitempty"`
	PrefetchPreload         *bool           `json:"prefetchPreload,omitempty"`
	PrivacyPass             *bool           `json:"privacyPass,omitempty"`
	ResponseBuffering       *bool           `json:"responseBuffering,omitempty"`
	RocketLoader            *bool           `json:"rocketLoader,omitempty"`
	SecurityHeader          *SecurityHeader `json:"securityHeader,omitempty"`
	SecurityLevel           *string         `json:"securityLevel,omitempty"`
	ServerSideExclude       *bool           `json:"serverSideExclude,omitempty"`
	SortQueryStringForCache *bool           `json:"sortQueryStringForCache,omitempty"`
	SSL                     *bool           `json:"ssl,omitempty"`
	MinTLSVersion           *string         `json:"minTLSVersion,omitempty"`
	Ciphers                 []*string       `json:"ciphers,omitempty"`
	TLS13                   *bool           `json:"tls13,omitempty"`
	TLSClientAuth           *bool           `json:"tlsClientAuth,omitempty"`
	TrueClientIPHeader      *bool           `json:"trueClientIPHeader,omitempty"`
	WAF                     *bool           `json:"waf,omitempty"`
	HTTP2                   *bool           `json:"http2,omitempty"`
	HTTP3                   *bool           `json:"http3,omitempty"`
	ZeroRTT                 *bool           `json:"0rtt,omitempty"`
	PseudoIPV4              *bool           `json:"pseudoIPV4,omitempty"`
	Websockets              *bool           `json:"websockets,omitempty"`
	ImageResizing           *bool           `json:"imageResizing,omitempty"`
	HTTP2Prioritization     *bool           `json:"http2Prioritization,omitempty"`
}

// ZoneSpec defines the desired state of Zone
type ZoneSpec struct {
	APIToken string        `json:"apiToken"`
	Settings *ZoneSettings `json:"settings,omitempty"`
}

// ZoneStatus defines the observed state of Zone
type ZoneStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Zone is the Schema for the zones API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Zone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZoneSpec   `json:"spec,omitempty"`
	Status ZoneStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ZoneList contains a list of Zone
type ZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Zone `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Zone{}, &ZoneList{})
}
