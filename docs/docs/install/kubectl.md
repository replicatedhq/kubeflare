# Install Kubeflare using Kubectl

Kubeflare requires Kubernetes 1.16 or later to install.

To install the current version of Kubeflare:

```shell
kubectl apply -f https://git.io/kubeflare
```

## Upgrading

To upgrade when there's a new release available, run the same command again:

```shell
kubectl apply -f https://git.io/kubeflare
```

## Uninstalling

To uninstall Kueflare completely from your cluster:

```shell
kubectl delete -f https://git.io/kubeflare
```

Or, if you need to uninstall manually:

```shell
kubectl delete crd apitokens.crds.kubeflare.io
kubectl delete crd dnsrecords.crds.kubeflare.io
kubectl delete crd tokens.crds.kubeflare.io
kubectl delete crd zones.crds.kubeflare.io
kubectl delete ns kubeflare-system
```