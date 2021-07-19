# Kubeflare Tutorial

This tutorial describes how to setup Kubeflare to manage your Cloudflare resources.

## Cloudflare Account Setup

If you don't have a Cloudflare account and site already, create one following the instructions [here] (https://support.cloudflare.com/hc/en-us/articles/201720164-Creating-a-Cloudflare-account-and-adding-a-website)

## Install

Install Kubeflare via [kubectl](/install/kubectl) or [Helm](/install/helm).

## Create a DNS Record

### Create a secret with your API Token

kubectl -n kubeflare-system create secret generic cf-api-secret --from-literal cf-api-token=<your-api-token>

### Define a Zone

```yaml
apiVersion: crds.kubeflare.io/v1alpha1
kind: Zone
metadata:
  name: domainname.io
spec:
  apiToken: 
    valueFrom:
      secretKeyRef:
        name: cf-api-secret
        key: cf-api-token
```

Full API Documentation for Zones in Kubeflare [here] (/api/zone).

### Create an A Record

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
