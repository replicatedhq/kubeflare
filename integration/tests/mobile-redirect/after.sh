#!/bin/bash

curl -X GET "https://api.cloudflare.com/client/v4/zones/${CF_ZONE_ID}/settings/mobile_redirect" \
     -H "X-Auth-Email: ${CF_API_EMAIL}" \
     -H "X-Auth-Key: ${CF_API_KEY}" \
     -H "Content-Type: application/json"