+++
title = "redis_cache"
description = "*redis_cache* - shared L2 DNS cache backed by a Redis-compatible key-value store."
weight = 10
tags = [  "plugin" , "redis_cache" ]
categories = [ "plugin", "external" ]
date = "2026-05-11T00:00:00+00:00"
repo = "https://github.com/dragoangel/coredns-redis-cache-plugin"
home = "https://github.com/dragoangel/coredns-redis-cache-plugin/blob/master/README.md"
+++

## Description

*redis_cache* stores DNS responses in a shared Redis-compatible backend (Redis, Valkey, or
any RESP-protocol server) so that multiple CoreDNS instances can amortize upstream lookups
across the fleet — for example several pods in a Kubernetes cluster, or a fleet of
node-local-dns daemons. It is intended to sit *behind* the built-in *cache* plugin, which
stays as the L1 (in-process) cache; *redis_cache* is the L2 (networked) cache.

If the Redis backend is unreachable the plugin becomes a noop and lookups continue to flow
through the rest of the chain. Writes never block the DNS reply (they run in a fire-and-forget
goroutine on a detached context). Reads are bounded by the configured `timeout read` budget
(default `500ms`) — the GET + TTL pipeline, pool wait and any retries all share that single
budget — so a stalled Redis adds at most one `read` timeout to a single DNS reply before the
plugin falls through. Read and write errors are surfaced via the `get_errors_total` and
`set_errors_total` metrics so a broken cache is distinguishable from a cold one.

Each response is cached for the duration of its record TTL, clamped into a configurable range:
`max(min, min(record_TTL, max))`. Defaults are `1h` max for positive responses and `30m` max
for denials, both with no minimum floor; raise or lower either bound via the `success` and
`denial` directives.

## Syntax

~~~ txt
redis_cache [ZONES...] {
    success MAX_TTL [MIN_TTL]
    denial MAX_TTL [MIN_TTL]
    endpoint ENDPOINT
    read_endpoint ENDPOINT [ENDPOINT...]
    key_prefix STRING
    db NUMBER
    sentinel MASTER_NAME SENTINEL_ADDR [SENTINEL_ADDR...]
    cluster SEED_ADDR [SEED_ADDR...]
    read_from latency|random|primary
    username USERNAME
    password PASSWORD
    sentinel_username USERNAME
    sentinel_password PASSWORD
    timeout {
        connect DURATION
        read DURATION
        write DURATION
    }
    pool {
        size N
        min_idle N
        max_idle N
        max_active N
        max_idle_time DURATION
        max_lifetime DURATION
        wait_timeout DURATION
    }
    retries {
        max N
        min_backoff DURATION
        max_backoff DURATION
    }
    tcp_keepalive DURATION
    tls
    tls_cert PATH
    tls_key PATH
    tls_ca PATH
    tls_verify_chain BOOL
    tls_verify_hostname BOOL
    resolver ADDRESS
}
~~~

Each sub-directive can be omitted; when present, its own arguments are required. Bare
`redis_cache` with no block attempts to connect to `127.0.0.1:6379` with default TTL
bounds — useful only against a sidecar Redis on localhost; production deployments must
specify at least one of `endpoint`, `sentinel`, or `cluster`. The chosen topology mode
determines which other directives are valid; the parser errors at load time on conflicting
combinations:

* `cluster` mode rejects `endpoint`, `read_endpoint`, `sentinel`, and any `db` other than
  `0` (Redis Cluster only supports DB 0). Seed addresses come from `cluster`; the rest of
  the topology is discovered via `CLUSTER SLOTS`.
* `sentinel` mode rejects `endpoint` and `read_endpoint` — the master and replicas are
  discovered via Sentinel.
* Default mode (neither `cluster` nor `sentinel`): writes go to `endpoint`. With no
  `read_endpoint`, the same client serves reads. With one, that client serves reads.
  With ≥2, each GET picks a replica at random. Rejects `read_from` and
  `sentinel_username` / `sentinel_password`.

* **ZONES** (positional) — zones to cache for. Defaults to the surrounding server-block
  zones.
* `success MAX_TTL [MIN_TTL]` — override TTL bounds for positive responses. MAX_TTL caps
  the cache duration (default `1h`). MIN_TTL sets a floor (default `0`) — when the upstream
  record TTL is shorter than this value, the cache duration is raised to this floor. Each
  value accepts a Go duration (`30s`, `1h`) or a bare integer (seconds); sub-second values
  like `500ms` are rejected.
* `denial MAX_TTL [MIN_TTL]` — same as `success` but for negative responses
  (NXDOMAIN/NODATA). Defaults: MAX_TTL `30m`, MIN_TTL `0`.
* `endpoint` — write endpoint address (default `127.0.0.1:6379`). Accepts IPs or hostnames.
  If a port is omitted, 6379 is assumed.
* `read_endpoint` — one or more read-only replica addresses. GETs route here, SETs go to
  `endpoint`. With ≥2 replicas, each GET picks one at random.
* `key_prefix STRING` — namespace prefix for cache keys (default `cdrc`). Keys are stored
  as `<key_prefix>:<hex>`; the `:` separator is appended automatically. Set to `""` to
  disable the prefix entirely (bare hex keys on a dedicated instance). A trailing `:` in
  the configured value is trimmed so `key_prefix mycache` and `key_prefix mycache:` are
  equivalent.
* `db NUMBER` — Redis logical database index for the data plane. Default `0`. Not allowed
  in `cluster` mode (Redis Cluster supports only DB 0).
* `sentinel` — enable Sentinel mode. **Master Group Name is mandatory** and must be
  followed by one or more sentinel addresses. The plugin discovers the current master and
  replicas via Sentinel (single quorum subscription); writes go to the master, reads pick
  a replica at random per GET.
* `cluster` — enable Cluster mode. Takes one or more seed node addresses; the smart client
  discovers the full topology via `CLUSTER SLOTS`.
* `read_from` — replica routing strategy in cluster mode. Only valid when `cluster` is set.
    * `latency` (default) — pick the replica with the lowest measured RTT.
    * `random` — pick a random replica.
    * `primary` — read only from primaries (no replica reads).
* `username` — ACL username for the data plane (primary, replicas, or cluster nodes).
  Optional.
* `password` — AUTH password for the data plane. Optional.
* `sentinel_username` — ACL username for the Sentinel API. Optional; only used in
  `sentinel` mode.
* `sentinel_password` — AUTH password for the Sentinel API. Optional; only used in
  `sentinel` mode.
* `timeout` — Redis connection and operation timeouts:
    * `connect` — TCP dial timeout (default: `1s`).
    * `read` — per-command read timeout (default: `500ms`).
    * `write` — per-command write timeout (default: `2s`).
* `pool` — connection-pool tuning. Values are non-negative integers.
    * `size N` — maximum sockets per client (default `10 × runtime.GOMAXPROCS()`).
    * `min_idle N` — minimum idle sockets to keep warm (default `0`).
    * `max_idle N` — maximum idle sockets (default `0` = unlimited).
    * `max_active N` — hard cap on total open sockets including in-use (default `0` =
      unlimited).
    * `max_idle_time DURATION` — close a connection that has been idle for this long
      (default `30m`). Set to less than your load balancer / NAT idle drop window.
    * `max_lifetime DURATION` — force-recycle any connection older than this regardless
      of activity (default `0` = no limit).
    * `wait_timeout DURATION` — how long a query waits for a free pool connection
      before erroring (default `500ms`).
* `retries` — retry behavior for transient network errors:
    * `max N` — number of retries per operation (default `1`), `0` disables retries.
    * `min_backoff DURATION` — initial backoff between retries (default `8ms` — go-redis).
    * `max_backoff DURATION` — cap on backoff between retries (default `512ms` — go-redis).
      Constraint: `min_backoff` must not exceed `max_backoff` when both are set.
* `tcp_keepalive DURATION` — TCP keepalive probe interval (default Go's built-in).
  Set below your NAT / firewall / mesh idle-drop window to prevent silent kills.
* `tls` — enable TLS. **No args.** Verifies the server cert against the OS trust store.
  Use `tls_ca` to override the trust store, `tls_cert` / `tls_key` for mTLS. Implicitly
  enabled by any other `tls_*` directive — bare `tls` is only needed when no other TLS knob
  is set. The TLS config applies to every connection the plugin opens (Sentinel API,
  master, replicas, cluster nodes); bundle CAs if planes use different roots.
* `tls_cert PATH` — PEM client certificate for mTLS. Must be paired with `tls_key`.
* `tls_key PATH` — PEM private key matching `tls_cert`.
* `tls_ca PATH` — PEM CA file used to verify the server certificate. **Replaces** the OS
  trust store when set; use only when your server's cert chains to a CA the OS doesn't
  ship.
* `tls_verify_chain BOOL` — verify the server certificate chains to a trusted root. Default
  `on`. Set to `off` to disable all server-cert verification (chain *and* hostname); use
  only for development or fully-trusted networks. Accepts `on`/`off`, `true`/`false`,
  `yes`/`no`, `1`/`0`.
* `tls_verify_hostname BOOL` — verify the server cert's SAN/CN matches the dialed
  hostname. Default `on`. Workaround for topologies where the dialed name cannot match
  the cert SAN (per-pod certs, Cluster MOVED redirects, Sentinel master/replica discovery,
  VIP fronting); chain verification still runs. Properly-issued certs should not require
  this. Has no effect when `tls_verify_chain` is `off`. See the example below.
* `resolver ADDRESS` — DNS server to use for resolving Redis endpoint hostnames instead of
  the system resolver. Useful in deployments where CoreDNS itself intercepts the system
  resolver (e.g. node-local-dns) and resolving the Redis service name through it would
  create a circular dependency. Set this to an upstream DNS service IP. Port defaults to
  53.

### Authentication

The data plane (Redis nodes) and the Sentinel API authenticate independently — credentials
across the two planes may be the same or different. In each plane the auth mode follows the
standard Redis convention:

* neither set → unauthenticated.
* password only → legacy `AUTH <password>` (matches `requirepass` on any version, or
  authenticates as the `default` user on ACL-enabled servers).
* username + password → full ACL auth (Redis 6+ for the data plane, Sentinel 6.2+ for the
  Sentinel API).

### Cache key isolation

The cache key is `xxhash64(qclass || qtype || DO || CD || lowercase(qname))`, namespaced
by `key_prefix`. All five components are mixed into the hash *and* re-verified after each
GET — a mismatch is treated as a miss, self-healed via async eviction, and reported via
`coredns_redis_cache_collisions_total`.

Practical guarantees this gives operators running mixed-client traffic:

* IN and CHAOS lookups (e.g. `version.bind.`) never share a slot with normal Internet
  queries for the same qname.
* DNSSEC-aware (`DO=1`) and non-DNSSEC clients keep separate entries — neither receives
  the other's response with extra or missing `RRSIG` / `NSEC` records.
* DNSSEC-validating (`CD=0`) and validation-bypassing (`CD=1`) queries are isolated. A
  CD=1 query for a DNSSEC-bogus name cannot poison the cache against a CD=0 client that
  would have received SERVFAIL from a validating upstream.

## Known Compatibility

The plugin speaks only standard RESP commands (`AUTH`, `GET`, `SET … EX`, `TTL`, `EXPIRE`,
`PING`, plus `CLUSTER SLOTS` in cluster mode and `SENTINEL get-master-addr-by-name` in
Sentinel mode), so it is expected to work with any reasonably complete Redis-protocol
implementation.

## Metrics

If monitoring is enabled (via the *prometheus* directive) the following metrics are exported:

* `coredns_redis_cache_hits_total{server}` - The count of cache hits from Redis.
* `coredns_redis_cache_misses_total{server}` - The count of cache misses from Redis.
* `coredns_redis_cache_get_errors_total{server,reason}` - The count of errors when reading entries from Redis. See *Error reasons* below for the `reason` buckets.
* `coredns_redis_cache_set_errors_total{server,reason}` - The count of errors when adding entries to Redis. Same `reason` buckets as `get_errors_total`.
* `coredns_redis_cache_encode_errors_total{server}` - The count of DNS messages that could not be serialized to wire format and so were not cached.
* `coredns_redis_cache_response_mismatches_total{server}` - The count of upstream replies whose question did not match the original request and were therefore refused for caching (the reply itself is still passed to the client). Non-zero suggests a misbehaving forwarder upstream or an attempted cache-poisoning probe.
* `coredns_redis_cache_collisions_total{server}` - The count of cache hits whose stored entry did not match the request (qname/qtype/qclass/DO/CD all re-verified after GET; mismatched entries are treated as a miss and asynchronously evicted). Should be zero in normal operation. The only innocent trigger is a statistical xxhash64 collision, which is ≈2⁻⁶⁴ per pair and effectively never fires at any plausible cache size. A non-zero value therefore points to a bug to investigate — Redis returning the wrong key's value, in-process mutation of cached bytes, or a coding error in this plugin — rather than something to ignore.

The `server` label indicates which server handled the request, see the *metrics* plugin for details.

### Error reasons

`get_errors_total` and `set_errors_total` are bucketed by `reason`:

* `timeout` - context deadline / cancellation, a network timeout, or a connection-pool wait
  timeout. Look at Redis latency / CPU, pool sizing, and the configured `timeout read` /
  `pool wait_timeout` budgets.
* `connection` - non-timeout network failures: dial refused, connection reset, EOF mid-op.
  Look at connectivity (DNS, firewall, route), and whether Redis is up and accepting
  connections.
* `other` - RESP-level errors (`NOAUTH`, `WRONGPASS`, parse failures, unhandled `MOVED`,
  etc.) or anything that isn't a network error. Typically a configuration or code issue
  rather than a transient outage.

## Examples

Examples after the first show only the `redis_cache { ... }` block; wrap it in the same
`. { cache {...} … forward . … }` shape from the Standalone example. They also omit
`success` / `denial` — reuse the values from Standalone or rely on the defaults documented
in the directive list.

Local L1 plus a shared Redis L2:

~~~ corefile
. {
    cache {
        success 9984 30
        denial 9984 5
    }
    redis_cache {
        endpoint redis.cache.svc.cluster.local:6379
        success 1h 1m
        denial 30m 30s
    }
    forward . 8.8.8.8:53
}
~~~

Writes to a known master, reads random-balanced across explicit replicas:

~~~ corefile
redis_cache {
    endpoint 10.0.0.1:6379
    read_endpoint 10.0.0.2:6379 10.0.0.3:6379
    password secretPass
}
~~~

Sentinel with separate data-plane and Sentinel-API passwords:

~~~ corefile
redis_cache {
    sentinel mymaster 10.0.0.1:26379 10.0.0.2:26379 10.0.0.3:26379
    password masterReplicaPass
    sentinel_password sentinelPass
}
~~~

Redis 6+ ACL (username + password):

~~~ corefile
redis_cache {
    endpoint redis.cache.svc.cluster.local:6379
    username dns-cache
    password s3cret
}
~~~

Cluster mode for capacity scaling beyond a single node's RAM:

~~~ corefile
redis_cache {
    cluster valkey-cluster-0:6379 valkey-cluster-1:6379 valkey-cluster-2:6379
    password secretPass
    read_from latency
}
~~~

> **Kubernetes note:** the smart client connects directly to every primary and replica the
> seeds advertise via `CLUSTER SLOTS`. If nodes advertise pod IPs (chart default), ensure
> they're routable from CoreDNS pods, or set `cluster-announce-hostname` on each node so
> the announced addresses match what `resolver` resolves.

TLS — server-only, OS trust store, no client cert:

~~~ corefile
redis_cache {
    endpoint redis.example.com:6380
    tls
    password s3cret
}
~~~

TLS — server-only, internal CA:

~~~ corefile
redis_cache {
    endpoint redis.example.com:6380
    tls_ca /etc/ssl/certs/redis-ca.pem
    password s3cret
}
~~~

TLS — mTLS:

~~~ corefile
redis_cache {
    endpoint redis.cache.svc.cluster.local:6379
    username dns-cache
    password s3cret
    tls_cert /etc/redis/tls/client.crt
    tls_key  /etc/redis/tls/client.key
    tls_ca   /etc/redis/tls/ca.pem
}
~~~

TLS — Kubernetes Redis Cluster with per-pod certs. Workaround for setups where issuing
certs whose SAN matches the dialed name is not practical: a StatefulSet-deployed
Redis/Valkey cluster typically presents per-pod certs (SAN =
`<pod>.<headless-svc>.<ns>.svc.cluster.local`), the client dials a service name, and
Cluster MOVED redirects further route to peers whose SANs won't match anything
pre-declared. Chain verification still applies to every peer:

~~~ corefile
redis_cache {
    cluster redis-cluster-0.redis-cluster-headless.cache.svc.cluster.local:6379 \
            redis-cluster-1.redis-cluster-headless.cache.svc.cluster.local:6379 \
            redis-cluster-2.redis-cluster-headless.cache.svc.cluster.local:6379
    tls_ca              /etc/redis/tls/ca.pem
    tls_verify_hostname off
    password            s3cret
}
~~~

Same workaround applies to Sentinel-discovered masters/replicas and HA-proxy/VIP fronting
a fleet of per-pod certs. Prefer issuing certs whose SAN covers the dialed name where you
control the PKI.

Kubernetes node-local-dns. When CoreDNS itself intercepts the cluster DNS VIP, resolving
the Redis service name through it would loop. Use `resolver` to point at the upstream
kube-dns; `__PILLAR__CLUSTER__DNS__` is substituted by node-local-dns at runtime:

~~~ corefile
.:53 {
    errors
    cache {
        success 9984 30
        denial 9984 5
    }
    redis_cache {
        endpoint k8s-dns-cache-redis-master.k8s-dns-cache.svc.cluster.local:6379
        read_endpoint k8s-dns-cache-redis-replicas.k8s-dns-cache.svc.cluster.local:6379
        password secretPass
        success 1h 1m
        denial 30m 30s
        resolver __PILLAR__CLUSTER__DNS__
    }
    forward . __PILLAR__UPSTREAM__SERVERS__
}
~~~

## Building

Add this line to CoreDNS's `plugin.cfg`. **It must appear after the `cache:cache` line** so
the in-process *cache* runs as L1 and *redis_cache* as L2:

~~~ text
cache:cache
redis_cache:github.com/dragoangel/coredns-redis-cache-plugin
~~~

Then `go get "github.com/dragoangel/coredns-redis-cache-plugin@latest" go generate coredns.go && go build` in the CoreDNS source tree.

## See Also

See the [Redis](https://redis.io) and [Valkey](https://valkey.io) sites for backend details.

Spiritual successor to [miekg/redis](https://github.com/miekg/redis) (directive `redisc`, archived November 2025).
