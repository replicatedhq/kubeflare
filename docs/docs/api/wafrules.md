# WAF Rule

A `kind: WebApplicationFirewallRule` resource will manage one or more WAF rules in a managed [zone](../zone).
Each resource can contain one or more WAF rules when specified under `spec.rules`.

## Attributes

### Zone

Each `kind: WebApplicationFirewallRule` must contain a `spec.zone` string attribute.
The value of this attribute must match a [zone](../zone) managed by Kubeflare.
The API token to manage the WAF rule(s) will be read from the associated Zone kind resource.

### Rule(s)

For more information on this type, see the [Cloudflare documentation](https://api.cloudflare.com/#waf-rules-edit-rule).

The following attributes are supported in the `rules` object:

| Name | Type | Description |
|------|------|-------------|
| id | string | The WAF rule ID
| mode | string | The WAF rule mode
| packageid | string | The WAF rule package (optional)

## Examples

### Single WAF Rule

The following example will set the mode for WAF Rule `PHP10001` to `simulate`:

```yaml
apiVersion: crds.kubeflare.io/v1alpha1
kind: WebApplicationFirewallRule
metadata:
  name: php-100001
spec:
  zone: domainname.io
  rules:
    - id: "PHP100001"
      mode: "simulate"
```

### Multiple WAF Rules

The following example will configure multiple PHP WAF Rules for a domain:

```yaml
apiVersion: crds.kubeflare.io/v1alpha1
kind: WebApplicationFirewallRule
metadata:
  name: php-rules
spec:
  zone: domainname.io
  rules:
    - id: "PHP100001"
      mode: "simulate"
    - id: "PHP100011"
      mode: "challenge"
```
