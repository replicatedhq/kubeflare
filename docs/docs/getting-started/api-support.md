# Cloudflare API Support Matrix

Cloudflare has a lot of objects and supports configuring them through their API.
The Kubeflare controller is early and doesn't (yet) have support for all Cloudflare objects.
Our intention is to support them all, including enterprise-only functionality.

The table below is the current state of support for the Cloudflare API.
If an object is not included on the table, it's not supported yet.
Feel free to open an [issue](https://github.com/replicatedhq/kubeflare/issues/new) if there's a specific API you need or would like to help with.

| Cloudflare API | Status | Kubeflare Version |
|----------------|--------|-------------------|
| Zone Settings | Completed | 0.1.0 |
| DNS Records | Completed | 0.1.0 |
| PageRules | In Progress | 0.1.0 |
| Access Applications | In Progress | 0.1.0 |
| Web Application Firewall | In Progress | 0.1.0 |
| WorkerRoutes | Completed | (not released) |
