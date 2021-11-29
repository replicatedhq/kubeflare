# WorkerRoute

Routes allow users to map a URL pattern to a Worker script to enable Workers to run on custom domains.
A `kind: WorkerRoute` resource will manage one or more [Worker Routes](https://developers.cloudflare.com/workers/platform/routes)
in a managed [zone](../zone).
Each resource contains a single WorkerRoute record.
For more information on this type, see the [Cloudflare documentation](https://api.cloudflare.com/#worker-routes-properties).

## Attributes

### Zone

Each `kind: WorkerRoute` must contain a `spec.zone` string attribute.
The value of this attribute must match a [Zone](../zone) name managed by Kubeflare.
The API token to call Cloudflare API will be read from the associated Zone resource.

### Pattern

Route pattern for URL matching. (string)

Pattern rules:

- Route patterns must include your zone. If your zone is `example.com`, then the simplest possible route pattern you can have is `example.com`, which would match `http://example.com/` and `https://example.com/`, and nothing else. As with a URL, there is an implied path of / if you do not specify one.
- Route patterns may not contain any query parameters. For example, https://example.com/?anything is not a valid route pattern.
- If you omit a scheme in your route pattern, it will match both `http://` and `https://` URLs. If you include `http://` or `https://`, it will only match HTTP or HTTPS requests, respectively.
- If a route pattern hostname begins with `*`, then it matches the host and all subhosts. If a route pattern hostname begins with `*.`, then it matches only all subhosts.
- If a route pattern path ends with `*`, then it matches all suffixes of that path.

### Script

The name of a worker script to execute on matched requests. (string)

## Importing

To import existing worker routes run check the [Importing](../getting-started/importing.md) guide.

YAML spec created by `kubeflare import` command contains a special annotation `crds.kubeflare.io/imported-id`.
When kubeflare operator detects a new resource with this annotation, it checks if a worker route with given ID exists.
If the worker route found, operator updates `status.id`, otherwise operator creates a new worker route.

## Examples

### Execute a `worker-script` on http(s)://example.com/images/*

The following `WorkerRoute` manifest will execute a `images-script` on URLs matching `http(s)://example.com/images/*`

```yaml
apiVersion: crds.kubeflare.io/v1alpha1
kind: WorkerRoute
metadata:
  name: images-script
spec:
  pattern: example.com/images/*
  script: "images-script"
  zone: example-com
```
