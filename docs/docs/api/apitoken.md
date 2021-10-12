# APIToken

A `kind: APIToken` resource represents a credential that can be used to call Cloudflare API. This token can be provided
inline (not recommended) or referenced from a Kubernetes secret.

### Examples

#### Inline

```yaml
apiVersion: crds.kubeflare.io/v1alpha1
kind: APIToken
metadata:
  name: api-token-name
spec:
  name: token-name
  email: email-address
  value: the-secret-api-token
```

#### Referencing a Kubernetes Secret

```yaml
apiVersion: crds.kubeflare.io/v1alpha1
kind: APIToken
metadata:
  name: api-token-name
spec:
  name: token-name
  email: email-address
  valueFrom:
    secretKeyRef:
      name: secret-name
      key: key-in-secret
```
