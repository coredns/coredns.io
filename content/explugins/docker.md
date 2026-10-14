+++
title = "docker"
description = "*example* - enables reading zone data from the Docker Daemon"
weight = 10
tags = [  "plugin" , "docker" ]
categories = [ "plugin", "external" ]
date = "2025-12-18T18:39:00+04:00"
repo = "https://github.com/dokku/coredns-docker"
home = "https://github.com/dokku/coredns-docker/blob/main/README.md"
+++

## Description

The *docker* plugin provides a DNS interface to Docker containers. It connects to the local Docker daemon, subscribes to the container event stream, and maintains a live view of container names, network aliases, Docker DNS names, Docker Compose `project.service` pairs, and custom names attached via Docker labels. The plugin answers A, AAAA, SRV, TXT, CNAME, and PTR queries for the containers it observes, along with a synthetic SOA and NS record at each configured zone apex.

This is an out-of-tree plugin maintained at <https://github.com/dokku/coredns-docker>. It must be compiled into CoreDNS by adding `docker:github.com/dokku/coredns-docker` to `plugin.cfg` before running `go generate` and `go build`, or by downloading a pre-built binary from the [releases page](https://github.com/dokku/coredns-docker/releases).

## Syntax

~~~ txt
docker
~~~

With only the plugin name specified and no options, the *docker* plugin answers queries in the `docker.` zone using records derived from the local Docker daemon with default settings.

~~~ txt
docker {
    zone ZONE [ZONE...]
    ttl SECONDS
    label_prefix PREFIX
    max_backoff DURATION
    networks NETWORK [NETWORK...]
    fallthrough [ZONES...]
    host_mode [ptr]
}
~~~

* `zone` **ZONE [ZONE...]** sets the DNS zones the plugin is authoritative for. Multiple zones may be listed. Defaults to `docker.`. Trailing dots are optional; an empty zone is rejected.
* `ttl` **SECONDS** sets the TTL on all answers. Valid range is `0` to `3600`. Defaults to `30`. A TTL of `0` disables downstream caching.
* `label_prefix` **PREFIX** sets the Docker label namespace the plugin reads when looking for custom records. Defaults to `com.dokku.coredns-docker`.
* `max_backoff` **DURATION** caps the exponential backoff used when reconnecting to a Docker daemon that has become unreachable. Defaults to `60s`.
* `networks` **NETWORK [NETWORK...]** whitelists the Docker networks the plugin serves. Containers attached only to non-whitelisted networks are ignored. If omitted, every network is served.
* `fallthrough` **[ZONES...]** If a query for a record in the zones for which the plugin is authoritative results in NXDOMAIN, normally that is what the response will be. However, if this option is specified, the query will instead be passed on down the plugin chain. If **[ZONES...]** is omitted, fallthrough happens for all zones for which the plugin is authoritative.
* `host_mode` **[ptr]** resolves container names to the host IP and host port of each container's port bindings instead of the container's internal network IP. With the optional `ptr` flag, PTR records are also generated for host IPs (off by default to reduce reverse-lookup noise).

## Name Sources

Without any labels, every container automatically receives A/AAAA records for:

* the container name set via `docker run --name` or `container_name:` in Compose,
* Docker network aliases,
* Docker-internal DNS names,
* the `project.service` pair for containers managed by Docker Compose.

PTR records pointing back to these names are generated for both IPv4 (`in-addr.arpa.`) and IPv6 (`ip6.arpa.`) reverse zones. To serve reverse lookups the reverse zones must be listed in the CoreDNS server block so that CoreDNS delivers PTR queries to the plugin:

~~~ txt
docker.:1053 in-addr.arpa:1053 ip6.arpa:1053 {
    docker
}
~~~

## Docker Labels

Additional records are configured by attaching Docker labels to containers. Every label uses the prefix set via `label_prefix` (default `com.dokku.coredns-docker`).

* `hostname=NAME[,NAME...]` registers one or more extra A/AAAA names for the container alongside its default names. Whitespace around each name is trimmed; empty values are ignored.
* `cname=TARGET` makes every name the container would receive a CNAME to **TARGET**. Per RFC 1034 §3.6.2, A/AAAA/SRV/PTR/TXT records are suppressed for containers with a `cname` label. Targets without a trailing dot have one appended automatically.
* `txt=VALUE` attaches a TXT record to the container's FQDN.
* `txt.KEY=VALUE` attaches a TXT record to `KEY.<container>.<zone>`. Multiple `txt.*` labels on the same container accumulate as separate TXT resource records. Values that start with a double quote are parsed as RFC 1035 master-file TXT rdata, supporting multi-string values and standard escapes (`\"`, `\\`, `\DDD`). Values longer than 255 bytes are automatically split into multiple character-strings on the wire.
* `srv._PROTO._SERVICE=PORT` advertises an SRV record at `_SERVICE._PROTO.<container>.<zone>`. If no `srv` labels are set, the plugin derives SRV records from the container's exposed ports (`NetworkSettings.Ports`).
* `wildcard=true` generates wildcard records (`*.<container>.<zone>`) alongside the exact records. Per RFC 4592 a wildcard matches exactly one label, and exact matches always take precedence.

## Host Mode

With `host_mode` enabled, the plugin reports host-bound IP addresses and host ports from each container's port bindings instead of the container's internal Docker network IP. This is the appropriate configuration when CoreDNS runs **outside** Docker (for example, directly on a developer laptop) and needs to return addresses the host can actually reach.

In host mode:

* A/AAAA records use the host IP from each port binding.
* SRV records report the host port, not the container port.
* Wildcard bind addresses are normalized: `0.0.0.0` becomes `127.0.0.1`, `::` becomes `::1`.
* Containers without published ports produce no records.
* PTR records are off by default. Pass the `ptr` flag to `host_mode` to opt back in.

## Stale Records

When the Docker daemon becomes unreachable, the plugin continues to serve the last set of records it synchronized. During this window answers carry a TTL of `5` seconds (or the configured `ttl` if it is already lower) so clients re-query quickly once the daemon returns. The plugin reconnects with exponential backoff up to `max_backoff`, then re-synchronizes and resumes normal TTLs.

## Ready

This plugin reports readiness to the *ready* plugin. It is ready once it has successfully connected to the Docker daemon, or when it has previously synced records at least once (stale mode).

## Metrics

If monitoring is enabled (via the *prometheus* plugin) then the following metrics are exported:

* `coredns_docker_success_requests_total{server}` - count of DNS requests handled successfully.
* `coredns_docker_failed_requests_total{server}` - count of DNS requests that failed.
* `coredns_docker_request_duration_seconds{server, type}` - histogram of DNS request latency, labeled by CoreDNS server and query type.
* `coredns_docker_stale_requests_total{server}` - count of requests served from stale data while the Docker daemon was disconnected.
* `coredns_docker_fallthrough_requests_total{server}` - count of requests passed to the next plugin via fallthrough.
* `coredns_docker_last_sync_timestamp_seconds` - Unix timestamp of the last successful record sync.
* `coredns_docker_records_total` - current number of A/AAAA record names tracked.
* `coredns_docker_srv_records_total` - current number of SRV record names tracked.
* `coredns_docker_ptr_records_total` - current number of PTR record names tracked.
* `coredns_docker_cname_records_total` - current number of CNAME record names tracked.
* `coredns_docker_txt_records_total` - current number of TXT record names tracked.
* `coredns_docker_connected` - `1` when connected to the Docker daemon, `0` otherwise.
* `coredns_docker_containers_total` - number of Docker containers currently tracked.
* `coredns_docker_sync_duration_seconds` - histogram of record sync latency in seconds.
* `coredns_docker_sync_errors_total` - count of failed record sync attempts.

## Examples

Serve the `docker.` zone and answer reverse lookups for containers on IPv4 and IPv6 networks.

~~~ txt
docker.:1053 in-addr.arpa:1053 ip6.arpa:1053 {
    docker {
        zone docker.
    }
    cache 30
}
~~~

Serve multiple zones from a single Corefile. Both the server block and the plugin block must list every zone.

~~~ txt
docker.:1053 internal.:1053 in-addr.arpa:1053 ip6.arpa:1053 {
    docker {
        zone docker. internal.
    }
}
~~~

Run CoreDNS outside Docker and rely on each container's host port bindings instead of internal Docker network IPs:

~~~ txt
docker.:1053 {
    docker {
        zone docker.
        host_mode
    }
}
~~~

Serve container names for a limited set of networks and forward everything else upstream through the *forward* plugin:

~~~ txt
. :1053 {
    docker {
        zone docker.
        networks bridge app
        fallthrough
    }
    forward . 1.1.1.1 1.0.0.1
    cache 30
}
~~~

## Bugs

Reverse (PTR) lookups for host-bound IP addresses are off by default when `host_mode` is enabled, because multiple containers often share the same host IP (especially when binding to loopback). Use `host_mode ptr` to opt back in and accept the noisier reverse-lookup output.

## See Also

See the *ready* plugin for readiness endpoints, the *prometheus* plugin to scrape the metrics listed above, the *cache* plugin to cache answers across the CoreDNS chain, and the *forward* plugin to resolve non-Docker queries upstream when `fallthrough` is enabled.

The plugin source, pre-built releases, and additional documentation live at <https://github.com/dokku/coredns-docker>.
