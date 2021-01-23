# PageRule

A `kind: PageRule` resource will manage one or more PageRules in a manage [zone](../zone).
Each resource contains a single PageRule record.

## Attributes

### Zone

Each `kind: PageRule` must contain a `spec.zone` string attribute.
The value of this attribute must match a [zone](../zone) managed by Kubeflare.
The API token to manage the DNS record(s) will be read from the associated Zone kind resource.

### PageRule

For more information on this type, see the [Cloudflare documentation](https://api.cloudflare.com/#page-rules-for-a-zone-create-page-rule).

The following attributes are supported in the `pagerule` or `pagerules` object.
Priority and status are optional, and exactly one of the other fields (rules) should be present on the object.

| Name | Type | Description |
|------|------|-------------|
| requestUrl | string | The incoming (original) request url |
| priority | int | | 
| status | string | |
| forwardingUrl | [ForwardingURL](#forwardingURL) | When present, the forwarding url page rule
| alwaysUseHttps | [AlwaysUseHTTPS](#alwaysUseHTTPS) | When present, the always use https page rule

### ForwardingURL

The ForwardingURL object describes a forwarding url page rule.

| Name | Type | Description |
|------|------|-------------|
| statusCode | int | 301 or 302, the status code to send |
| redirectUrl | string | The redirect/forwarded url |

### AlwaysUseHTTPS

The AlwaysUseHTTPS object is an empty object that enables the always use https pagerule.

## Examples

### Redirect (forward) www to apex domain

The following PageRule manifest will forward requests made to www.example.com to example.com:

```yaml
apiVersion: crds.kubeflare.io/v1alpha1
kind: PageRule
metadata:
  name: www.example.com
spec:
  zone: example.com
  rule:
    requestUrl: "www.example.com/*"
    forwardingUrl:
      statusCode: 302
      redirectUrl: "https://example.com/$1
```

### Always Use HTTPS

The following PageRule manifest will enable always use https on a specific path:

```yaml
apiVersion: crds.kubeflare.io/v1alpha1
kind: PageRule
metadata:
  name: www.example.com
spec:
  zone: example.com
  rule:
    requestUrl: "www.example.com/*"
    alwaysUseHttps: {}
```