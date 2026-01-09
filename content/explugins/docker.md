+++
title = "docker"
description = "*example* - enables reading zone data from the Docker Daemon"
weight = 10
tags = [  "plugin" , "docker" ]
categories = [ "plugin", "external" ]
date = "2025-12-18T18:39:00+04:00"
repo = "https://github.com/dokku/coredns-docker"
home = "https://github.com/dokku/coredns-docker/blob/master/README.md"
+++

## Description

The docker plugin serves DNS records for containers running on the local Docker daemon. It follows the Docker event stream, picking up changes whenever something happens to a container - whether it gets created, started, deleted, or restarted.

The plugin resolves container names, network aliases, DNS names, and SRV records to their respective container IP addresses within a specified network.

SRV records can be defined using container labels with the prefix `[LABEL_PREFIX].srv.`, followed by the protocol and service name. For example, with the default prefix, a label `com.dokku.coredns-docker.srv._tcp._http=80` will create an SRV record for `_http._tcp.container-name.domain` pointing to the container's IP on port 80.

If no labels with the specified prefix are found, the plugin falls back to using the container's exposed ports (`NetworkSettings.Ports`).

- For a port mapping like `80/tcp`, it generates an SRV record for `_tcp._tcp.container-name.domain`.
- For a port mapping without a protocol like `80`, it generates SRV records for both `_tcp._tcp` and `_udp._udp`.

## Compilation

It will require you to use `go get` or as a dependency on [plugin.cfg](https://github.com/coredns/coredns/blob/master/plugin.cfg).

A simple way to consume this plugin, is by adding the following on [plugin.cfg](https://github.com/coredns/coredns/blob/master/plugin.cfg), and recompile it as [detailed on coredns.io](https://coredns.io/2017/07/25/compile-time-enabling-or-disabling-plugins/#build-with-compile-time-configuration-file).

```text
docker:github.com/dokku/coredns-docker
```

After this you can compile coredns by running:

```bash
make
```

## Syntax

```text
docker [DOMAIN] {
    ttl DURATION
    label_prefix PREFIX
    max_backoff DURATION
    networks NETWORK...
}
```

- `DOMAIN` is the domain for which the plugin will respond. Defaults to `docker.`.

- `ttl` allows you to set a custom TTL for responses. **DURATION** defaults to `30 seconds`. The minimum TTL allowed is `0` seconds, and the maximum is capped at `3600` seconds. Setting TTL to 0 will prevent records from being cached. The unit for the value is seconds.

- `label_prefix` allows you to set a custom prefix for SRV record labels. **PREFIX** defaults to `com.dokku.coredns-docker`.

- `max_backoff` allows you to set a maximum backoff duration for the Docker event loop reconnection logic. **DURATION** defaults to `60s`.

- `networks` allows you to specify a list of Docker networks to monitor. If specified, containers not on one of these networks will be ignored.

## Metrics

If monitoring is enabled (via the *prometheus* directive) the following metric is exported:

- `coredns_docker_success_requests_total{server}` - Counter of DNS requests handled successfully.
- `coredns_docker_failed_requests_total{server}` - Counter of DNS requests failed.

The `server` label indicated which server handled the request.

## Ready

This plugin reports readiness to the ready plugin. It will be ready only when it has successfully connected to the Docker daemon.

## Examples

Enable docker with and resolve all containers with `.docker.` as the suffix.

```text
docker:1053 {
    docker docker.
    cache 30
}
```

You can see the [Corefile.example](./Corefile.example) for a full Corefile example.

## Usage Example

### A record

```shell
dig web.docker @127.0.0.1 -p 1053    

; <<>> DiG 9.18.1-1ubuntu1.2-Ubuntu <<>> web.docker @127.0.0.1 -p 1053
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 54986
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;web.docker.  IN A

;; ANSWER SECTION:
web.docker. 30 IN A 172.17.0.2

;; Query time: 4 msec
;; SERVER: 127.0.0.1#1053(127.0.0.1) (UDP)
```

### SRV record

```shell
dig _http._tcp.web.docker @127.0.0.1 -p 1053 SRV

; <<>> DiG 9.18.1-1ubuntu1.2-Ubuntu <<>> _http._tcp.web.docker @127.0.0.1 -p 1053 SRV
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 49945
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;_http._tcp.web.docker.  IN SRV

;; ANSWER SECTION:
_http._tcp.web.docker. 30 IN SRV 10 10 80 web.docker.

;; Query time: 0 msec
;; SERVER: 127.0.0.1#1053(127.0.0.1) (UDP)
```
