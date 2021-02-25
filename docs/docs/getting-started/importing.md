# Importing

If you already have a Cloudflare Zone, the `kubeflare` CLI can be used to download the supported resources and configuration and write these as Kubernetes manifests. 
Once you've reviewed the generated YAML documents, deploy them to a cluster that has Kubeflare running, and the Cloudflare zone will be managed by the operator now.

## Installing the CLI

In order to import a zone, you'll need access to the CLI.
We recommend installing this on your laptop, but it's possible to run it in a Docker container (docs needed).

To install the CLI, download the `kubeflare` binary from the [releases](https://github.com/replicatedhq/kubeflare/releases) page, extract it, and put the binary on your path.

## Importing

The `import` command doesn't communicate with your cluster. 
It requires a read only [API Token](https://support.cloudflare.com/hc/en-us/articles/200167836-Managing-API-Tokens-and-Keys#12345680) to the Cloudflare account containing the zone. 

!!! note
    Be sure to create an API Token, not an API Key for this process. 
    The `import` command will fail when given an API Key.

The `import` command requires an API Token and a Zone name.
All other fields are optional. 
By default, all record types are extracted and written as YAML manifests in `./imported/`. 
Any existing files in this directory will be overwritten if there is a name conflict.

```shell
./kubeflare import --api-token 000000-0000-REDACTED_000000_REDACTED_00 --zone myzone.com
```

## Limiting the types of resources to import

It's possible to limit the types of resources to import.
Because the CLI imports them all, you can disable one or more types of resources in the CLI.

For example, to download everything except DNS records:

```shell
./kubeflare import --api-token 000000-0000-REDACTED_000000_REDACTED_00 --zone myzone.com --dns-records=false
```

The full list of optional types to exclude is in the table below:

| Flag | Record Type(s) 
|------|---------------
| `--dns-records` | All DNS Records


