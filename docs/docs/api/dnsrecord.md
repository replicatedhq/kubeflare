# DNS Record

A `kind: DNSRecord` resource will manage one or more DNS records in a managed [zone](../zone).
Each resource can contain a single DNS record when specified under `spec.record` vs `spec.records`.

## Attributes

### Zone

Each `kind: DNSRecord` must contain a `spec.zone` string attribute.
The value of this attribute must match a [zone](../zone) managed by Kubeflare.
The API token to manage the DNS record(s) will be read from the associated Zone kind resource.

### Record(s)

For more information on this type, see the [Cloudflare documentation](https://api.cloudflare.com/#dns-records-for-a-zone-create-dns-record).

The following attributes are supported in the `record` or `records` object:

| Name | Type | Description |
|------|------|-------------|
| type | string | The DNS record type
| name | string | The DNS record name
| content | string | The DNS record content
| ttl | int | Time to live for the record. (set to 1 for auto)
| proxied (optional)| boolean | When true, proxy the record through Cloudflare 

## Examples

### Single A Record

The following example will ensure that a single A record for `www` exists, proxied, and pointing to 1.1.1.1:

```yaml
apiVersion: crds.kubeflare.io/v1alpha1
kind: DNSRecord
metadata:
  name: www.domainname.io
spec:
  zone: domainname.io
  record:
    type: "A"
    name: "www"
    content: "1.1.1.1"
    proxied: true
    ttl: 3600
```

### Multiple MX Records

The following example will configure GSuite MX records for a domain:

```yaml
apiVersion: crds.kubeflare.io/v1alpha1
kind: DNSRecord
metadata:
  name: mx-records
spec:
  zone: domainname.io
  records:
    - type: "MX"
      name: "domainname.io"
      content: "aspmx.l.google.com"
      priority: 1
    - type: "MX"
      name: "domainname.io"
      content: "alt1.aspmx.l.google.com"
      priority: 5
    - type: "MX"
      name: "domainname.io"
      content: "alt2.aspmx.l.google.com"
      priority: 5
    - type: "MX"
      name: "domainname.io"
      content: "alt3.aspmx.l.google.com"
      proxied: false      
      priority: 10
    - type: "MX"
      name: "domainname.io"
      content: "alt4.aspmx.l.google.com"
      priority: 10
```