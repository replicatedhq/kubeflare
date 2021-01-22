# DNS Record

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