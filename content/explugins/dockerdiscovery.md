+++
title = "dockerdiscovery"
description = "*dockerdiscovery* - add/remove DNS records for Docker containers."
weight = 50
tags = [  "plugin" , "docker", "discovery", "containers" ]
categories = [ "plugin", "external" ]
date = "2023-11-01T09:10:00-03:00"
repo = "https://github.com/kevinjqiu/coredns-dockerdiscovery"
home = "https://github.com/kevinjqiu/coredns-dockerdiscovery/blob/master/README.md"
+++

## Description

Docker discovery plugin for CoreDNS.

## Syntax

    docker [DOCKER_ENDPOINT] {
        domain DOMAIN_NAME
        hostname_domain HOSTNAME_DOMAIN_NAME
        network_aliases DOCKER_NETWORK
        label LABEL
        compose_domain COMPOSE_DOMAIN_NAME
    }

* `DOCKER_ENDPOINT`: the path to the docker socket. If unspecified, defaults to `unix:///var/run/docker.sock`. It can also be TCP socket, such as `tcp://127.0.0.1:999`.
* `DOMAIN_NAME`: the name of the domain for [container name](https://docs.docker.com/engine/reference/run/#name---name), e.g. when `DOMAIN_NAME` is `docker.loc`, your container with `my-nginx` (as subdomain) [name](https://docs.docker.com/engine/reference/run/#name---name) will be assigned the domain name: `my-nginx.docker.loc`
* `HOSTNAME_DOMAIN_NAME`: the name of the domain for [hostname](https://docs.docker.com/config/containers/container-networking/#ip-address-and-hostname). Work same as `DOMAIN_NAME` for hostname.
* `COMPOSE_DOMAIN_NAME`: the name of the domain when it is determined the
    container is managed by docker-compose.  e.g. for a compose project of
    "internal" and service of "nginx", if `COMPOSE_DOMAIN_NAME` is
    `compose.loc` the fqdn will be `nginx.internal.compose.loc`
* `DOCKER_NETWORK`: the name of the docker network. Resolve directly by [network aliases](https://docs.docker.com/v17.09/engine/userguide/networking/configure-dns) (like internal docker dns resolve host by aliases whole network)
* `LABEL`: container label of resolving host (by default enable and equals ```coredns.dockerdiscovery.host```)

## Examples

~~~ corefile
. {
    docker unix:///var/run/docker.sock {
        domain {$COREDNS_DOCKER_DOMAIN}
        hostname_domain {$COREDNS_DOCKER_HOSTNAME_DOMAIN}
        compose_domain {$COREDNS_DOCKER_COMPOSE_DOMAIN}

        domain docker
        hostname_domain docker
        compose_domain docker

        domain local
        hostname_domain local
        compose_domain local

        ttl 30

        network_aliases net1
        network_aliases net2
    }
}
~~~
