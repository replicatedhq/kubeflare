# Install Kubeflare using Kubectl

Kubeflare requires Kubernetes 1.16 or later to install.

To install the current version of Kubeflare:

```shell
git clone git@github.com:replicatedhq/kubeflare.git
```

```shell
kubectl apply -f kubeflare/config/crds/v1
```

```shell
cat <<EOF | kubectl apply -f -
---
apiVersion: v1
kind: Namespace
metadata:
  name: kubeflare-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kubeflare
rules:
- apiGroups: ['']
  resources:
  - namespaces
  - secrets
  verbs: [get, list, watch]
- apiGroups: [crds.kubeflare.io]
  resources: ['*']
  verbs: [get, list, watch, update, patch, create, delete]
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kubeflare
  namespace: kubeflare-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kubeflare
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kubeflare
subjects:
- name: kubeflare
  namespace: kubeflare-system
  kind: ServiceAccount
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kubeflare
  namespace: kubeflare-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kubeflare
  template:
    metadata:
      labels:
        app: kubeflare
    spec:
      serviceAccountName: kubeflare
      containers:
      - name: kubeflare
        image: replicated/kubeflare-manager:0.3.0
        imagePullPolicy: IfNotPresent
        args:
        - --metrics-addr=:8088
        - --poll-interval=5m
EOF
```

## Upgrading

To upgrade when there's a new release available, run the same commands from the install section again, but with the updated image tag.

## Uninstalling

To uninstall Kueflare completely from your cluster:

```shell
kubectl delete crd apitokens.crds.kubeflare.io
kubectl delete crd dnsrecords.crds.kubeflare.io
kubectl delete crd tokens.crds.kubeflare.io
kubectl delete crd zones.crds.kubeflare.io
kubectl delete ns kubeflare-system
```
