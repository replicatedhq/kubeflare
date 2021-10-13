// +build !ignore_autogenerated

/*
Copyright 2020 Replicated, Inc.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIToken) DeepCopyInto(out *APIToken) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIToken.
func (in *APIToken) DeepCopy() *APIToken {
	if in == nil {
		return nil
	}
	out := new(APIToken)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *APIToken) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APITokenList) DeepCopyInto(out *APITokenList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]APIToken, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APITokenList.
func (in *APITokenList) DeepCopy() *APITokenList {
	if in == nil {
		return nil
	}
	out := new(APITokenList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *APITokenList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APITokenSpec) DeepCopyInto(out *APITokenSpec) {
	*out = *in
	if in.ValueFrom != nil {
		in, out := &in.ValueFrom, &out.ValueFrom
		*out = new(ValueFrom)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APITokenSpec.
func (in *APITokenSpec) DeepCopy() *APITokenSpec {
	if in == nil {
		return nil
	}
	out := new(APITokenSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APITokenStatus) DeepCopyInto(out *APITokenStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APITokenStatus.
func (in *APITokenStatus) DeepCopy() *APITokenStatus {
	if in == nil {
		return nil
	}
	out := new(APITokenStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccessApplication) DeepCopyInto(out *AccessApplication) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccessApplication.
func (in *AccessApplication) DeepCopy() *AccessApplication {
	if in == nil {
		return nil
	}
	out := new(AccessApplication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AccessApplication) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccessApplicationList) DeepCopyInto(out *AccessApplicationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DNSRecord, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccessApplicationList.
func (in *AccessApplicationList) DeepCopy() *AccessApplicationList {
	if in == nil {
		return nil
	}
	out := new(AccessApplicationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AccessApplicationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccessApplicationSpec) DeepCopyInto(out *AccessApplicationSpec) {
	*out = *in
	if in.AllowedIdPs != nil {
		in, out := &in.AllowedIdPs, &out.AllowedIdPs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AutoRedirectToIdentity != nil {
		in, out := &in.AutoRedirectToIdentity, &out.AutoRedirectToIdentity
		*out = new(bool)
		**out = **in
	}
	if in.CORSHeaders != nil {
		in, out := &in.CORSHeaders, &out.CORSHeaders
		*out = new(CORSHeader)
		(*in).DeepCopyInto(*out)
	}
	if in.AccessPolicies != nil {
		in, out := &in.AccessPolicies, &out.AccessPolicies
		*out = make([]AccessPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccessApplicationSpec.
func (in *AccessApplicationSpec) DeepCopy() *AccessApplicationSpec {
	if in == nil {
		return nil
	}
	out := new(AccessApplicationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccessApplicationStatus) DeepCopyInto(out *AccessApplicationStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccessApplicationStatus.
func (in *AccessApplicationStatus) DeepCopy() *AccessApplicationStatus {
	if in == nil {
		return nil
	}
	out := new(AccessApplicationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AccessPolicy) DeepCopyInto(out *AccessPolicy) {
	*out = *in
	if in.Include != nil {
		in, out := &in.Include, &out.Include
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Precedence != nil {
		in, out := &in.Precedence, &out.Precedence
		*out = new(int)
		**out = **in
	}
	if in.Exclude != nil {
		in, out := &in.Exclude, &out.Exclude
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Require != nil {
		in, out := &in.Require, &out.Require
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AccessPolicy.
func (in *AccessPolicy) DeepCopy() *AccessPolicy {
	if in == nil {
		return nil
	}
	out := new(AccessPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AlwaysUseHTTPSPageRule) DeepCopyInto(out *AlwaysUseHTTPSPageRule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AlwaysUseHTTPSPageRule.
func (in *AlwaysUseHTTPSPageRule) DeepCopy() *AlwaysUseHTTPSPageRule {
	if in == nil {
		return nil
	}
	out := new(AlwaysUseHTTPSPageRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoMinifyPageRule) DeepCopyInto(out *AutoMinifyPageRule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoMinifyPageRule.
func (in *AutoMinifyPageRule) DeepCopy() *AutoMinifyPageRule {
	if in == nil {
		return nil
	}
	out := new(AutoMinifyPageRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CORSHeader) DeepCopyInto(out *CORSHeader) {
	*out = *in
	if in.AllowedMethods != nil {
		in, out := &in.AllowedMethods, &out.AllowedMethods
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AllowedOrigins != nil {
		in, out := &in.AllowedOrigins, &out.AllowedOrigins
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AllowedHeaders != nil {
		in, out := &in.AllowedHeaders, &out.AllowedHeaders
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CORSHeader.
func (in *CORSHeader) DeepCopy() *CORSHeader {
	if in == nil {
		return nil
	}
	out := new(CORSHeader)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSRecord) DeepCopyInto(out *DNSRecord) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSRecord.
func (in *DNSRecord) DeepCopy() *DNSRecord {
	if in == nil {
		return nil
	}
	out := new(DNSRecord)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DNSRecord) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSRecordList) DeepCopyInto(out *DNSRecordList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DNSRecord, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSRecordList.
func (in *DNSRecordList) DeepCopy() *DNSRecordList {
	if in == nil {
		return nil
	}
	out := new(DNSRecordList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DNSRecordList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSRecordSpec) DeepCopyInto(out *DNSRecordSpec) {
	*out = *in
	if in.Record != nil {
		in, out := &in.Record, &out.Record
		*out = new(Record)
		(*in).DeepCopyInto(*out)
	}
	if in.Records != nil {
		in, out := &in.Records, &out.Records
		*out = make([]*Record, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Record)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSRecordSpec.
func (in *DNSRecordSpec) DeepCopy() *DNSRecordSpec {
	if in == nil {
		return nil
	}
	out := new(DNSRecordSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSRecordStatus) DeepCopyInto(out *DNSRecordStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSRecordStatus.
func (in *DNSRecordStatus) DeepCopy() *DNSRecordStatus {
	if in == nil {
		return nil
	}
	out := new(DNSRecordStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EdgeCacheTTLPageRule) DeepCopyInto(out *EdgeCacheTTLPageRule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EdgeCacheTTLPageRule.
func (in *EdgeCacheTTLPageRule) DeepCopy() *EdgeCacheTTLPageRule {
	if in == nil {
		return nil
	}
	out := new(EdgeCacheTTLPageRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ForwardingURLPageRule) DeepCopyInto(out *ForwardingURLPageRule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ForwardingURLPageRule.
func (in *ForwardingURLPageRule) DeepCopy() *ForwardingURLPageRule {
	if in == nil {
		return nil
	}
	out := new(ForwardingURLPageRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MinifySetting) DeepCopyInto(out *MinifySetting) {
	*out = *in
	if in.CSS != nil {
		in, out := &in.CSS, &out.CSS
		*out = new(bool)
		**out = **in
	}
	if in.HTML != nil {
		in, out := &in.HTML, &out.HTML
		*out = new(bool)
		**out = **in
	}
	if in.JS != nil {
		in, out := &in.JS, &out.JS
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MinifySetting.
func (in *MinifySetting) DeepCopy() *MinifySetting {
	if in == nil {
		return nil
	}
	out := new(MinifySetting)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MobileRedirect) DeepCopyInto(out *MobileRedirect) {
	*out = *in
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(bool)
		**out = **in
	}
	if in.MobileSubdomain != nil {
		in, out := &in.MobileSubdomain, &out.MobileSubdomain
		*out = new(string)
		**out = **in
	}
	if in.StripURI != nil {
		in, out := &in.StripURI, &out.StripURI
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MobileRedirect.
func (in *MobileRedirect) DeepCopy() *MobileRedirect {
	if in == nil {
		return nil
	}
	out := new(MobileRedirect)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PageRule) DeepCopyInto(out *PageRule) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PageRule.
func (in *PageRule) DeepCopy() *PageRule {
	if in == nil {
		return nil
	}
	out := new(PageRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PageRule) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PageRuleList) DeepCopyInto(out *PageRuleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PageRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PageRuleList.
func (in *PageRuleList) DeepCopy() *PageRuleList {
	if in == nil {
		return nil
	}
	out := new(PageRuleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PageRuleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PageRuleSpec) DeepCopyInto(out *PageRuleSpec) {
	*out = *in
	if in.Rule != nil {
		in, out := &in.Rule, &out.Rule
		*out = new(Rule)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PageRuleSpec.
func (in *PageRuleSpec) DeepCopy() *PageRuleSpec {
	if in == nil {
		return nil
	}
	out := new(PageRuleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PageRuleStatus) DeepCopyInto(out *PageRuleStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PageRuleStatus.
func (in *PageRuleStatus) DeepCopy() *PageRuleStatus {
	if in == nil {
		return nil
	}
	out := new(PageRuleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Record) DeepCopyInto(out *Record) {
	*out = *in
	if in.TTL != nil {
		in, out := &in.TTL, &out.TTL
		*out = new(int)
		**out = **in
	}
	if in.Priority != nil {
		in, out := &in.Priority, &out.Priority
		*out = new(int)
		**out = **in
	}
	if in.Proxied != nil {
		in, out := &in.Proxied, &out.Proxied
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Record.
func (in *Record) DeepCopy() *Record {
	if in == nil {
		return nil
	}
	out := new(Record)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rule) DeepCopyInto(out *Rule) {
	*out = *in
	if in.ForwardingURL != nil {
		in, out := &in.ForwardingURL, &out.ForwardingURL
		*out = new(ForwardingURLPageRule)
		**out = **in
	}
	if in.AlwaysUseHTTPS != nil {
		in, out := &in.AlwaysUseHTTPS, &out.AlwaysUseHTTPS
		*out = new(AlwaysUseHTTPSPageRule)
		**out = **in
	}
	if in.AutoMinify != nil {
		in, out := &in.AutoMinify, &out.AutoMinify
		*out = new(AutoMinifyPageRule)
		**out = **in
	}
	if in.EdgeCacheTTL != nil {
		in, out := &in.EdgeCacheTTL, &out.EdgeCacheTTL
		*out = new(EdgeCacheTTLPageRule)
		**out = **in
	}
	if in.Priority != nil {
		in, out := &in.Priority, &out.Priority
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rule.
func (in *Rule) DeepCopy() *Rule {
	if in == nil {
		return nil
	}
	out := new(Rule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecurityHeader) DeepCopyInto(out *SecurityHeader) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
	if in.MaxAge != nil {
		in, out := &in.MaxAge, &out.MaxAge
		*out = new(int)
		**out = **in
	}
	if in.IncludeSubdomains != nil {
		in, out := &in.IncludeSubdomains, &out.IncludeSubdomains
		*out = new(bool)
		**out = **in
	}
	if in.NoSniff != nil {
		in, out := &in.NoSniff, &out.NoSniff
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecurityHeader.
func (in *SecurityHeader) DeepCopy() *SecurityHeader {
	if in == nil {
		return nil
	}
	out := new(SecurityHeader)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ValueFrom) DeepCopyInto(out *ValueFrom) {
	*out = *in
	if in.SecretKeyRef != nil {
		in, out := &in.SecretKeyRef, &out.SecretKeyRef
		*out = new(v1.SecretKeySelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ValueFrom.
func (in *ValueFrom) DeepCopy() *ValueFrom {
	if in == nil {
		return nil
	}
	out := new(ValueFrom)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WAFRule) DeepCopyInto(out *WAFRule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WAFRule.
func (in *WAFRule) DeepCopy() *WAFRule {
	if in == nil {
		return nil
	}
	out := new(WAFRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebApplicationFirewallRule) DeepCopyInto(out *WebApplicationFirewallRule) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebApplicationFirewallRule.
func (in *WebApplicationFirewallRule) DeepCopy() *WebApplicationFirewallRule {
	if in == nil {
		return nil
	}
	out := new(WebApplicationFirewallRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WebApplicationFirewallRule) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebApplicationFirewallRuleList) DeepCopyInto(out *WebApplicationFirewallRuleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]WebApplicationFirewallRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebApplicationFirewallRuleList.
func (in *WebApplicationFirewallRuleList) DeepCopy() *WebApplicationFirewallRuleList {
	if in == nil {
		return nil
	}
	out := new(WebApplicationFirewallRuleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WebApplicationFirewallRuleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebApplicationFirewallRuleSpec) DeepCopyInto(out *WebApplicationFirewallRuleSpec) {
	*out = *in
	if in.Rules != nil {
		in, out := &in.Rules, &out.Rules
		*out = make([]*WAFRule, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(WAFRule)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebApplicationFirewallRuleSpec.
func (in *WebApplicationFirewallRuleSpec) DeepCopy() *WebApplicationFirewallRuleSpec {
	if in == nil {
		return nil
	}
	out := new(WebApplicationFirewallRuleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WebApplicationFirewallRuleStatus) DeepCopyInto(out *WebApplicationFirewallRuleStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WebApplicationFirewallRuleStatus.
func (in *WebApplicationFirewallRuleStatus) DeepCopy() *WebApplicationFirewallRuleStatus {
	if in == nil {
		return nil
	}
	out := new(WebApplicationFirewallRuleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Zone) DeepCopyInto(out *Zone) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Zone.
func (in *Zone) DeepCopy() *Zone {
	if in == nil {
		return nil
	}
	out := new(Zone)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Zone) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZoneList) DeepCopyInto(out *ZoneList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Zone, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZoneList.
func (in *ZoneList) DeepCopy() *ZoneList {
	if in == nil {
		return nil
	}
	out := new(ZoneList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ZoneList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZoneSettings) DeepCopyInto(out *ZoneSettings) {
	*out = *in
	if in.AdvancedDDOS != nil {
		in, out := &in.AdvancedDDOS, &out.AdvancedDDOS
		*out = new(bool)
		**out = **in
	}
	if in.AlwaysOnline != nil {
		in, out := &in.AlwaysOnline, &out.AlwaysOnline
		*out = new(bool)
		**out = **in
	}
	if in.AlwaysUseHTTPS != nil {
		in, out := &in.AlwaysUseHTTPS, &out.AlwaysUseHTTPS
		*out = new(bool)
		**out = **in
	}
	if in.OpportunisticOnion != nil {
		in, out := &in.OpportunisticOnion, &out.OpportunisticOnion
		*out = new(bool)
		**out = **in
	}
	if in.AutomaticHTTPSRewrites != nil {
		in, out := &in.AutomaticHTTPSRewrites, &out.AutomaticHTTPSRewrites
		*out = new(bool)
		**out = **in
	}
	if in.BrowserCacheTTL != nil {
		in, out := &in.BrowserCacheTTL, &out.BrowserCacheTTL
		*out = new(int)
		**out = **in
	}
	if in.BrowserCheck != nil {
		in, out := &in.BrowserCheck, &out.BrowserCheck
		*out = new(bool)
		**out = **in
	}
	if in.CacheLevel != nil {
		in, out := &in.CacheLevel, &out.CacheLevel
		*out = new(string)
		**out = **in
	}
	if in.ChallengeTTL != nil {
		in, out := &in.ChallengeTTL, &out.ChallengeTTL
		*out = new(int)
		**out = **in
	}
	if in.DevelopmentMode != nil {
		in, out := &in.DevelopmentMode, &out.DevelopmentMode
		*out = new(bool)
		**out = **in
	}
	if in.EmailObfuscation != nil {
		in, out := &in.EmailObfuscation, &out.EmailObfuscation
		*out = new(bool)
		**out = **in
	}
	if in.HotlinkProtection != nil {
		in, out := &in.HotlinkProtection, &out.HotlinkProtection
		*out = new(bool)
		**out = **in
	}
	if in.IPGeolocation != nil {
		in, out := &in.IPGeolocation, &out.IPGeolocation
		*out = new(bool)
		**out = **in
	}
	if in.IPV6 != nil {
		in, out := &in.IPV6, &out.IPV6
		*out = new(bool)
		**out = **in
	}
	if in.Minify != nil {
		in, out := &in.Minify, &out.Minify
		*out = new(MinifySetting)
		(*in).DeepCopyInto(*out)
	}
	if in.MobileRedirect != nil {
		in, out := &in.MobileRedirect, &out.MobileRedirect
		*out = new(MobileRedirect)
		(*in).DeepCopyInto(*out)
	}
	if in.Mirage != nil {
		in, out := &in.Mirage, &out.Mirage
		*out = new(bool)
		**out = **in
	}
	if in.OriginErrorPagePassThru != nil {
		in, out := &in.OriginErrorPagePassThru, &out.OriginErrorPagePassThru
		*out = new(bool)
		**out = **in
	}
	if in.OpportunisticEncryption != nil {
		in, out := &in.OpportunisticEncryption, &out.OpportunisticEncryption
		*out = new(bool)
		**out = **in
	}
	if in.Polish != nil {
		in, out := &in.Polish, &out.Polish
		*out = new(bool)
		**out = **in
	}
	if in.WebP != nil {
		in, out := &in.WebP, &out.WebP
		*out = new(bool)
		**out = **in
	}
	if in.Brotli != nil {
		in, out := &in.Brotli, &out.Brotli
		*out = new(bool)
		**out = **in
	}
	if in.PrefetchPreload != nil {
		in, out := &in.PrefetchPreload, &out.PrefetchPreload
		*out = new(bool)
		**out = **in
	}
	if in.PrivacyPass != nil {
		in, out := &in.PrivacyPass, &out.PrivacyPass
		*out = new(bool)
		**out = **in
	}
	if in.ResponseBuffering != nil {
		in, out := &in.ResponseBuffering, &out.ResponseBuffering
		*out = new(bool)
		**out = **in
	}
	if in.RocketLoader != nil {
		in, out := &in.RocketLoader, &out.RocketLoader
		*out = new(bool)
		**out = **in
	}
	if in.SecurityHeader != nil {
		in, out := &in.SecurityHeader, &out.SecurityHeader
		*out = new(SecurityHeader)
		(*in).DeepCopyInto(*out)
	}
	if in.SecurityLevel != nil {
		in, out := &in.SecurityLevel, &out.SecurityLevel
		*out = new(string)
		**out = **in
	}
	if in.ServerSideExclude != nil {
		in, out := &in.ServerSideExclude, &out.ServerSideExclude
		*out = new(bool)
		**out = **in
	}
	if in.SortQueryStringForCache != nil {
		in, out := &in.SortQueryStringForCache, &out.SortQueryStringForCache
		*out = new(bool)
		**out = **in
	}
	if in.SSL != nil {
		in, out := &in.SSL, &out.SSL
		*out = new(bool)
		**out = **in
	}
	if in.MinTLSVersion != nil {
		in, out := &in.MinTLSVersion, &out.MinTLSVersion
		*out = new(string)
		**out = **in
	}
	if in.Ciphers != nil {
		in, out := &in.Ciphers, &out.Ciphers
		*out = make([]*string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(string)
				**out = **in
			}
		}
	}
	if in.TLS13 != nil {
		in, out := &in.TLS13, &out.TLS13
		*out = new(bool)
		**out = **in
	}
	if in.TLSClientAuth != nil {
		in, out := &in.TLSClientAuth, &out.TLSClientAuth
		*out = new(bool)
		**out = **in
	}
	if in.TrueClientIPHeader != nil {
		in, out := &in.TrueClientIPHeader, &out.TrueClientIPHeader
		*out = new(bool)
		**out = **in
	}
	if in.WAF != nil {
		in, out := &in.WAF, &out.WAF
		*out = new(bool)
		**out = **in
	}
	if in.HTTP2 != nil {
		in, out := &in.HTTP2, &out.HTTP2
		*out = new(bool)
		**out = **in
	}
	if in.HTTP3 != nil {
		in, out := &in.HTTP3, &out.HTTP3
		*out = new(bool)
		**out = **in
	}
	if in.ZeroRTT != nil {
		in, out := &in.ZeroRTT, &out.ZeroRTT
		*out = new(bool)
		**out = **in
	}
	if in.PseudoIPV4 != nil {
		in, out := &in.PseudoIPV4, &out.PseudoIPV4
		*out = new(bool)
		**out = **in
	}
	if in.Websockets != nil {
		in, out := &in.Websockets, &out.Websockets
		*out = new(bool)
		**out = **in
	}
	if in.ImageResizing != nil {
		in, out := &in.ImageResizing, &out.ImageResizing
		*out = new(bool)
		**out = **in
	}
	if in.HTTP2Prioritization != nil {
		in, out := &in.HTTP2Prioritization, &out.HTTP2Prioritization
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZoneSettings.
func (in *ZoneSettings) DeepCopy() *ZoneSettings {
	if in == nil {
		return nil
	}
	out := new(ZoneSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZoneSpec) DeepCopyInto(out *ZoneSpec) {
	*out = *in
	if in.Settings != nil {
		in, out := &in.Settings, &out.Settings
		*out = new(ZoneSettings)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZoneSpec.
func (in *ZoneSpec) DeepCopy() *ZoneSpec {
	if in == nil {
		return nil
	}
	out := new(ZoneSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ZoneStatus) DeepCopyInto(out *ZoneStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ZoneStatus.
func (in *ZoneStatus) DeepCopy() *ZoneStatus {
	if in == nil {
		return nil
	}
	out := new(ZoneStatus)
	in.DeepCopyInto(out)
	return out
}
