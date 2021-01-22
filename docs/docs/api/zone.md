# Zone

A `zone` record is created per domain name that's managed in Cloudflare.
Here you can specify the API Token and any Zone settings.

A `zone` record is required in order to configure any additional api types on a Cloudflare domain.

```yaml
apiVersion: crds.kubeflare.io/v1alpha1
kind: Zone
metadata:
  name: domainname.io
spec:
  apiToken: redacted
  settings:
    alwaysUseHttps: true
    alwaysOnline: true
    minify:
      css: true
```