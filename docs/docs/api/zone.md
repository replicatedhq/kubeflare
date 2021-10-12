# Zone

A `kind: Zone` resource should be created per domain name that's managed in Cloudflare. Here you can specify the API
Token and any Zone settings.

A `kind: Zone` resource is required in order to configure any additional api types on a Cloudflare domain. Each domain
managed by Kubeflare should have exactly 1 resource of this type deployed to the cluster.

All keys under `settings` default to the same default values as documented in the Cloudflare API.

If the `settings` key is not included in the manifest, no settings will be changed in the Cloudflare zone. If
the `settings` key is specified, only keys that are specified will be applied to the Cloudflare zone.

```yaml
apiVersion: crds.kubeflare.io/v1alpha1
kind: Zone
metadata:
  name: domainname.io
spec:
  apiToken: api-token-name
  settings:
    alwaysUseHttps: true
    alwaysOnline: true
    minify:
      css: true
```

## `apiToken`

Each zone should have an API Token specified. The value of this field should be the name of a `kind: APIToken` resource.

## `settings`

All [Cloudflare Zone Settings](https://api.cloudflare.com/#zone-settings-properties) can be specified in the `settings`
key of this resource. Kubeflare uses a lowerCamelCase standard to specify all fields in the Cloudflare Zone.

Note that the Cloudflare API and docs use string types with values of "off" and "on" for boolean settings. Kubeflare
uses boolean objects (true, false) and will map those to the string types accepted by Cloudflare.

| Kubeflare Setting | Cloudflare Setting | Data Type |
|-------------------|--------------------|-----------|
| advancedDDOS | advanced_ddos | boolean
| alwaysOnline | [always_online](https://api.cloudflare.com/#zone-settings-change-always-online-setting) | boolean
| alwaysUseHttps | [always_use_https](https://api.cloudflare.com/#zone-settings-change-always-use-https-setting) | boolean
| opportunisticOnion | [opportunistic_onion](https://api.cloudflare.com/#zone-settings-change-opportunistic-onion-setting) | boolean
| automaticHTTPSRewrites | [automatic_https_rewrites](https://api.cloudflare.com/#zone-settings-change-automatic-https-rewrites-setting) | boolean
| browserCacheTTL | [browser_cache_ttl](https://api.cloudflare.com/#zone-settings-change-browser-cache-ttl-setting) | int
| browserCheck | [browser_check](https://api.cloudflare.com/#zone-settings-change-browser-check-setting) | boolean
| cacheLevel | [cache_level](https://api.cloudflare.com/#zone-settings-change-cache-level-setting) | string
| challengeTTL | [challenge_ttl](https://api.cloudflare.com/#zone-settings-change-challenge-ttl-setting) | int
| developmentMode |[development_mode](https://api.cloudflare.com/#zone-settings-change-development-mode-setting) | boolean
| emailObfuscation | [email_obfuscation](https://api.cloudflare.com/#zone-settings-change-email-obfuscation-setting) | boolean
| hotlinkProtection | [hotlink_protection](https://api.cloudflare.com/#zone-settings-change-hotlink-protection-setting) | boolean
| ipGeolocation | [ip_geoloation](https://api.cloudflare.com/#zone-settings-change-ip-geolocation-setting) | boolean
| ipv6 | [ipv6](https://api.cloudflare.com/#zone-settings-change-ipv6-setting)| boolean
| minify | [minify](https://api.cloudflare.com/#zone-settings-change-minify-setting) |
| mobileRedirect | [mobile_redirect](https://api.cloudflare.com/#zone-settings-change-mobile-redirect-setting) |
| mirage | [mirage](https://api.cloudflare.com/#zone-settings-change-mirage-setting) | boolean
| originErrorPagePassThru | [origin_error_page_pass_thru](https://api.cloudflare.com/#zone-settings-change-enable-error-pages-on-setting) | boolean
| opportunisticEncryption | [opportunistic_encryption](https://api.cloudflare.com/#zone-settings-change-opportunistic-encryption-setting) | boolean
| polish | [polish](https://api.cloudflare.com/#zone-settings-change-polish-setting) | boolean
| webp | [webp](https://api.cloudflare.com/#zone-settings-change-webp-setting) | boolean
| brotli | [brotli](https://api.cloudflare.com/#zone-settings-change-brotli-setting) | boolean
| prefetchPreload | [prefetch_preload](https://api.cloudflare.com/#zone-settings-change-prefetch-preload-setting) | boolean
| privacyPass | [privacy_pass](https://api.cloudflare.com/#zone-settings-change-privacy-pass-setting) | boolean
| responseBuffering | [response_buffering](https://api.cloudflare.com/#zone-settings-change-response-buffering-setting) | boolean
| rocketLoader | [rocket_loader](https://api.cloudflare.com/#zone-settings-change-rocket-loader-setting) | boolean
| securityHeader | [security_header](https://api.cloudflare.com/#zone-settings-change-security-header-hsts-setting) |
| securityLevel | [security_level](https://api.cloudflare.com/#zone-settings-change-security-level-setting) | string
| serverSideExclude | [server_side_exclude](https://api.cloudflare.com/#zone-settings-change-server-side-exclude-setting) | boolean
| sortQueryStringForCache | [sort_query_string_for_cache](https://api.cloudflare.com/#zone-settings-change-enable-query-string-sort-setting) | boolean
| ssl | [ssl](https://api.cloudflare.com/#zone-settings-change-ssl-setting) | boolean
| minTLSVersion | [min_tls_version](https://api.cloudflare.com/#zone-settings-change-minimum-tls-version-setting) | string
| ciphers |[ciphers](https://api.cloudflare.com/#zone-settings-change-ciphers-setting) | []string
| tls13 | [tls_1_3](https://api.cloudflare.com/#zone-settings-change-tls-1.3-setting) | boolean
| tlsClientAuth | [tls_client_auth](https://api.cloudflare.com/#zone-settings-change-tls-client-auth-setting) | boolean
| trueClientIPHeader | [true_client_ip_header](https://api.cloudflare.com/#zone-settings-change-true-client-ip-setting) | boolean
| waf | [waf](https://api.cloudflare.com/#zone-settings-change-web-application-firewall-waf-setting) | boolean
| http2 | [http2](https://api.cloudflare.com/#zone-settings-change-http2-setting) | boolean
| http3 | [http3](https://api.cloudflare.com/#zone-settings-change-http3-setting) | boolean
| 0rtt | [0rtt](https://api.cloudflare.com/#zone-settings-change-0-rtt-session-resumption-setting) | boolean
| pseudoIPV4 | [pseudo_ipv4](https://api.cloudflare.com/#zone-settings-change-pseudo-ipv4-setting) | boolean
| websockets | [websockets](https://api.cloudflare.com/#zone-settings-change-websockets-setting) | boolean
| imageResizing | [image_resizing](https://api.cloudflare.com/#zone-settings-change-image-resizing-setting) | boolean
| http2Prioritization | [h2_prioritization](https://api.cloudflare.com/#zone-settings-change-http/2-edge-prioritization-setting) | boolean
