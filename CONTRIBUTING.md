# Contribuing Guide


### Integration tests

These require a domain and API token from Cloudflare.
There's a test domain in the GitHub repo and GitHub Actions uses it when running.

To execute locally:

```
export CF_API_EMAIL=you@you.com
export CF_API_KEY=1234
export CF_ZONE_ID=12345
export CF_ZONE_NAME=domain.com

make integration
```

To run a single integration test (useful when adding a new api call):

```
export CF_API_EMAIL=you@you.com
export CF_API_KEY=1234
export CF_ZONE_ID=12345
export CF_ZONE_NAME=domain.com

make -C integration/tests/cache-level run
```

## Setup

- You should create an A record (proxied) for "mobile".
