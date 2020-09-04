#!/bin/bash

curl -X PATCH "https://api.cloudflare.com/client/v4/zones/${CF_ZONE_ID}/settings/min_tls_version" \
     -H "X-Auth-Email: ${CF_API_EMAIL}" \
     -H "X-Auth-Key: ${CF_API_KEY}" \
     -H "Content-Type: application/json" \
     --data '{"value":"1.0"}'
