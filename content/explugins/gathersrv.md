+++
title = "gathersrv"
description = "*gathersrv* plugin allows to gather DNS responses with SRV records from several domains (for example k8s clusters) and hide them behind a single common/distributed domain"
weight = 10
tags = [  "plugin" , "gathersrv" ]
categories = [ "plugin", "external" ]
date = "2023-07-02T12:00:00+00:00"
repo = "https://github.com/ziollek/gathersrv"
home = "https://github.com/ziollek/gathersrv#readme"
+++

## Description

This plugin could be helpful for services that are logically distributed over several k8s clusters and use [headless service](https://kubernetes.io/docs/concepts/services-networking/service/#headless-services) to expose themselves.
The aim of this plugin is to provide a method to discover all service instances through a single service domain. The result of querying distributed service domain contains
masqueraded results gathered from multiple clusters.
In contrast to [multicluster plugin](https://github.com/coredns/multicluster) it does not require that k8s clusters have to share *the same cluster zone*.
This plugin works as a proxy that generates requests to configured clusters and later on, rewrites them before returning them to the client. So it is a bit more flexible (it can be used / run outside k8s).

## Syntax

~~~ txt
gathersrv DISTRIBUTED_ZONE {
    CLUSTER_DOMAIN_ONE HOSTNAME_PREFIX_ONE
    ...
    CLUSTER_DOMAIN_N HOSTNAME_PREFIX_N
}
~~~

## Example

Let's assume there are:
* two k8s clusters - cluster-a, cluster-b
* two zones for above clusters cluster-a.local, cluster-b.local
* possibility to query k8s dns outside cluster
* headless service - demo-service deployed on above clusters in the same namespace (default)

If we ask about

```
dig -t SRV _demo._tcp.demo-service.default.svc.cluster-a.local
```

we will see a result as below:

```
;; ANSWER SECTION:
_demo._tcp.demo-service.default.svc.cluster-a.local. 30 IN SRV 0 50 8080 demo-service-0.default.svc.cluster-a.local.
_demo._tcp.demo-service.default.svc.cluster-a.local. 30 IN SRV 0 50 8080 demo-service-1.default.svc.cluster-a.local.

;; ADDITIONAL SECTION:
demo-service-0.default.svc.cluster-a.local. 30 IN A 10.8.1.2
demo-service-1.default.svc.cluster-a.local. 30 IN A 10.8.1.2
```

Respectively for the second cluster

```
dig -t SRV _demo._tcp.demo-service.default.svc.cluster-a.local

...

;; ANSWER SECTION:
_demo._tcp.demo-service.default.svc.cluster-b.local. 30 IN SRV 0 50 8080 demo-service-0.default.svc.cluster-b.local.
_demo._tcp.demo-service.default.svc.cluster-b.local. 30 IN SRV 0 50 8080 demo-service-1.default.svc.cluster-b.local.

;; ADDITIONAL SECTION:
demo-service-0.default.svc.cluster-b.local. 30 IN A 10.9.1.2
demo-service-1.default.svc.cluster-b.local. 30 IN A 10.9.1.2
```

Using gathersrv plugin with coredns we can configure it to provide merged information behind single domain - in this case distributed.local



```
dig -t SRV _demo._tcp.demo-service.default.svc.distributed.local

...

;; ANSWER SECTION:
_demo._tcp.demo-service.default.svc.distributed.local. 30 IN SRV 0 50 8080 a-demo-service-0.default.svc.distributed.local.
_demo._tcp.demo-service.default.svc.distributed.local. 30 IN SRV 0 50 8080 a-demo-service-1.default.svc.distributed.local.
_demo._tcp.demo-service.default.svc.distributed.local. 30 IN SRV 0 50 8080 b-demo-service-0.default.svc.distributed.local.
_demo._tcp.demo-service.default.svc.distributed.local. 30 IN SRV 0 50 8080 b-demo-service-1.default.svc.distributed.local.

;; ADDITIONAL SECTION:
a-demo-service-0.default.svc.distributed.local. 30 IN A 10.8.1.2
a-demo-service-1.default.svc.distributed.local. 30 IN A 10.8.1.2
b-demo-service-0.default.svc.distributed.local. 30 IN A 10.9.1.2
b-demo-service-1.default.svc.distributed.local. 30 IN A 10.9.1.2
```


As shown above - the result response not only contains proper ip addresses but also translated hostnames.
This translation adds some prefix that indicates the original cluster and replaces the cluster domain (.cluster-a.local., .cluster-b.local.) with a distributed domain.
In effect service hostnames share their parent domain with service - a-demo-service-0.**default.svc.distributed.local.**.
Thanks to that the result could be consumed by restricted service drivers for example [mongodb+srv](https://docs.mongodb.com/manual/reference/connection-string/#dns-seed-list-connection-format).

It is worth mentioning that the POD's ip addresses will need to be routable outside cluster-a and cluster-b if you want to connect to them.

Below configuration reflects example from the use case. Addresses of dns service for cluster-a and cluster-b are 10.8.0.1, 10.9.0.1 respectively.

```
distributed.local. {
  gathersrv distribiuted.local. {
	cluster-a.local. a-
	cluster-b.local. b-
  }
  forward . 127.0.0.1:5300
}

cluster-a.local.:5300 {
  forward . 10.8.0.1:53
}

cluster-b.local.:5300 {
  forward . 10.9.0.1:53
}
```
