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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/replicatedhq/kubeflare/pkg/apis/crds/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDNSRecords implements DNSRecordInterface
type FakeDNSRecords struct {
	Fake *FakeCrdsV1alpha1
	ns   string
}

var dnsrecordsResource = schema.GroupVersionResource{Group: "crds.kubeflare.io", Version: "v1alpha1", Resource: "dnsrecords"}

var dnsrecordsKind = schema.GroupVersionKind{Group: "crds.kubeflare.io", Version: "v1alpha1", Kind: "DNSRecord"}

// Get takes name of the dNSRecord, and returns the corresponding dNSRecord object, and an error if there is any.
func (c *FakeDNSRecords) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.DNSRecord, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(dnsrecordsResource, c.ns, name), &v1alpha1.DNSRecord{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DNSRecord), err
}

// List takes label and field selectors, and returns the list of DNSRecords that match those selectors.
func (c *FakeDNSRecords) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.DNSRecordList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(dnsrecordsResource, dnsrecordsKind, c.ns, opts), &v1alpha1.DNSRecordList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.DNSRecordList{ListMeta: obj.(*v1alpha1.DNSRecordList).ListMeta}
	for _, item := range obj.(*v1alpha1.DNSRecordList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested dNSRecords.
func (c *FakeDNSRecords) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(dnsrecordsResource, c.ns, opts))

}

// Create takes the representation of a dNSRecord and creates it.  Returns the server's representation of the dNSRecord, and an error, if there is any.
func (c *FakeDNSRecords) Create(ctx context.Context, dNSRecord *v1alpha1.DNSRecord, opts v1.CreateOptions) (result *v1alpha1.DNSRecord, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(dnsrecordsResource, c.ns, dNSRecord), &v1alpha1.DNSRecord{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DNSRecord), err
}

// Update takes the representation of a dNSRecord and updates it. Returns the server's representation of the dNSRecord, and an error, if there is any.
func (c *FakeDNSRecords) Update(ctx context.Context, dNSRecord *v1alpha1.DNSRecord, opts v1.UpdateOptions) (result *v1alpha1.DNSRecord, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(dnsrecordsResource, c.ns, dNSRecord), &v1alpha1.DNSRecord{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DNSRecord), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeDNSRecords) UpdateStatus(ctx context.Context, dNSRecord *v1alpha1.DNSRecord, opts v1.UpdateOptions) (*v1alpha1.DNSRecord, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(dnsrecordsResource, "status", c.ns, dNSRecord), &v1alpha1.DNSRecord{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DNSRecord), err
}

// Delete takes name of the dNSRecord and deletes it. Returns an error if one occurs.
func (c *FakeDNSRecords) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(dnsrecordsResource, c.ns, name), &v1alpha1.DNSRecord{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDNSRecords) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(dnsrecordsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.DNSRecordList{})
	return err
}

// Patch applies the patch and returns the patched dNSRecord.
func (c *FakeDNSRecords) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.DNSRecord, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(dnsrecordsResource, c.ns, name, pt, data, subresources...), &v1alpha1.DNSRecord{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.DNSRecord), err
}
