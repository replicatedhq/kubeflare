# Kubeflare

Kuebflare is a Kubernetes cluster add-on (Operator) that allows you to manage your Cloudflare settings using a Kubernetes declarative API.

After installing the Kubeflare Operator to your Kubernetes cluster, some new custom types will be created in the cluster. These types allow you to define a Cloudflare Zone (domain) and specify the settings and DNS records to create. When this YAML is deployed to the cluster, the Kubeflare Operator will reconcile this with the Cloudflare API to deploy the settings requested in the YAML to the Cloudflare account.

## Motivation

This project was created at [Replicated](https://www.replicated.com) to manage Cloudflare settings using our GitOps workflow. We wanted a way for a devleoper to commit their DNS records and other Cloudflare settings to be reviewed and deployed with their code, as a single deployment. This tightly couples the infrastructure changes with the application changes, and makes deploying new services easier and more transparant.

## Examples

Below is an example of a Kubernetes manifest that we deploy for a domain (with some information redacted):

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
---
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
---
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

## Settings supported

The Cloudflare API is large and supports many settings. Kubeflare doesn't support all yet, but the current release of Kubeflare supports the following settings:

- All [zone settings](https://api.cloudflare.com/#zone-settings-get-all-zone-settings) (all settings listed under "Zone Settings")
- DNS Records 


This project is independent of Cloudflare and built using their public APIs. This is not a Cloudflare project.
