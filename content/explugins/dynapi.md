+++
title = "dynapi"
description = "*dynapi* manages CoreDNS A and AAAA record sets through an authenticated HTTP API."
weight = 10
tags = [ "plugin", "dynapi" ]
categories = [ "plugin", "external" ]
date = "2026-10-02T00:00:00+03:00"
repo = "https://github.com/coredns/dynapi"
home = "https://github.com/coredns/dynapi#readme"
import_path = "github.com/coredns/dynapi/plugins/dynapi"
+++

## Description

The *dynapi* plugin lets applications read, replace, and delete A and AAAA record
sets through an authenticated JSON HTTP API. It sends signed DNS requests to
*dynupdate* over loopback TCP. The *dynupdate* plugin serves the records, enforces
write permissions, and persists changes. The *tsig* plugin authenticates DNS
requests.

Reads use a signed zone transfer (AXFR) to retrieve the exact stored record set.
Writes use DNS UPDATE transactions to replace or delete a complete record set
atomically. Ordinary DNS queries continue through the CoreDNS plugin chain.

HTTP clients authenticate with a bearer token. DNS requests use a separate TSIG
key. The bearer token permits reads throughout the configured zone, while writes
must match the key's *dynupdate* permission rules.

## Syntax

```corefile
dynapi [ADDRESS] {
    token_env VARIABLE
    upstream ADDRESS
    identity KEY
    secret_env VARIABLE
    max_requests NUMBER
}
```

* **ADDRESS** is the HTTP listener address. It defaults to `127.0.0.1:8080`.
* `token_env` reads the HTTP bearer token from an environment variable. Use
  `token TOKEN` instead to configure a literal token.
* `upstream` specifies the DNS listener in the same server block. Both the HTTP
  and DNS addresses must be literal loopback IP addresses with ports.
* `identity` names the TSIG key used for DNS requests and *dynupdate* write permissions.
* `secret_env` reads the base64-encoded TSIG secret from an environment variable.
  Use `secret BASE64_SECRET` instead to configure a literal secret.
* `max_requests` limits concurrent HTTP requests and backend DNS connections.
  It defaults to 32. Excess requests receive HTTP 503.

Configure exactly one form of each credential. The bearer token must contain at
least 32 characters without whitespace. The TSIG secret must encode at least
16 bytes. The key uses HMAC-SHA256.

The server block must serve one zone and include *dynupdate*, *tsig*, and
*transfer*. DNS over TCP must be enabled.

## HTTP API

The [OpenAPI specification](https://github.com/coredns/dynapi/blob/v0.1.0/openapi.yaml)
documents the request and response schemas, bearer authentication, and error
responses for v0.1.0.

Requests use the following path, with `type` set to `A` or `AAAA`:

```text
/v1/zones/{zone}/records/{name}/{type}
```

Names must be full names within the configured zone. A trailing dot is optional.

* `GET` returns the exact stored record set as a JSON object containing `ttl`
  and `addresses`. A missing set returns HTTP 404.
* `PUT` atomically replaces the complete set and returns its normalized JSON
  payload. It requires `Content-Type: application/json`, an explicit TTL, and
  1–256 addresses of the requested family.
* `DELETE` removes the complete set for that name and type. It returns HTTP 204,
  including when the set is already absent.

For example, a successful GET returns:

```json
{"ttl":60,"addresses":["192.0.2.10","192.0.2.11"]}
```

## Examples

Create `example.org.zone` with the initial zone records:

```dns
$ORIGIN example.org.
@ 60 IN SOA ns.example.org. hostmaster.example.org. 1 3600 600 86400 60
@ 60 IN NS ns.example.org.
ns 60 IN A 127.0.0.1
```

This Corefile serves the zone on port 1053 and exposes the HTTP API on port 8080.
The TSIG key can change only A and AAAA records for `host.example.org`.

```corefile
example.org:1053 {
    bind 127.0.0.1

    dynapi 127.0.0.1:8080 {
        token_env DYNAPI_TOKEN
        upstream 127.0.0.1:1053
        identity update-key.example.org.
        secret_env DYNAPI_TSIG_SECRET
    }

    tsig {
        secret update-key.example.org. {$DYNAPI_TSIG_SECRET}
        require_opcode UPDATE
        require AXFR
    }

    transfer {
        to 127.0.0.1
    }

    dynupdate {
        file example.org.zone
        database example.org.db
        allow update-key.example.org. host.example.org. A AAAA
    }
}
```

Generate the two credentials and start CoreDNS:

```sh
export DYNAPI_TOKEN="$(openssl rand -hex 32)"
export DYNAPI_TSIG_SECRET="$(openssl rand -base64 32)"
./coredns -conf Corefile
```

In another terminal with the same `DYNAPI_TOKEN`, replace the host's A records:

```sh
curl -X PUT http://127.0.0.1:8080/v1/zones/example.org/records/host.example.org/A \
  -H "Authorization: Bearer $DYNAPI_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"ttl":60,"addresses":["192.0.2.10","192.0.2.11"]}'
```

Read the stored records or delete them:

```sh
curl http://127.0.0.1:8080/v1/zones/example.org/records/host.example.org/A \
  -H "Authorization: Bearer $DYNAPI_TOKEN"

curl -X DELETE http://127.0.0.1:8080/v1/zones/example.org/records/host.example.org/A \
  -H "Authorization: Bearer $DYNAPI_TOKEN"
```

## Limitations

* GET reads a full zone transfer and is intended for small zones. It does not
  expand wildcards or follow CNAMEs.
* TTL controls DNS caching. It does not expire stored records. DNS caches can
  retain older answers until their TTL expires.
* A failed connection or lost HTTP response can follow a committed write.
  There are no automatic retries or conditional writes.
* Keep both credentials if clients must continue using them after a restart.
  This version requires a restart to change configuration and rejects Corefile
  reloads while keeping the existing service running.

## See Also

See the [plugin repository](https://github.com/coredns/dynapi#readme) for further
details, architecture documentation, and Go client examples.
