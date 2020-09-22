+++
title = "k8s_dns_chaos"
description = "*k8s_dns_chaos* - enables inject DNS chaos in a Kubernetes cluster for Chaos Engineering."
weight = 10
tags = [  "plugin" , "example" ]
categories = [ "plugin", "external" ]
date = "2020-09-22T21:57:00+08:00"
repo = "https://github.com/chaos-mesh/k8s_dns_chaos"
home = "https://github.com/chaos-mesh/k8s_dns_chaos/blob/master/README.md"
+++

## Description

This plugin implements the [Kubernetes DNS-Based Service Discovery
Specification](https://github.com/kubernetes/dns/blob/master/docs/specification.md).

CoreDNS running with the k8s_dns_chaos plugin can be used to do chaos tests on DNS.

This plugin can only be used once per Server Block.

## Syntax

~~~
k8s_dns_chaos [ZONES...]
~~~

The *k8s_dns_chaos* supports all options in plugin *[kubernetes](https://coredns.io/plugins/kubernetes/)*, besides, it also supports other configuration items for chaos.

```
kubernetes [ZONES...] {
    endpoint URL
    tls CERT KEY CACERT
    kubeconfig KUBECONFIG CONTEXT
    namespaces NAMESPACE...
    labels EXPRESSION
    pods POD-MODE
    endpoint_pod_names
    ttl TTL
    noendpoints
    transfer to ADDRESS...
    fallthrough [ZONES...]
    ignore empty_service

    chaos ACTION SCOPE [PODS...]
    grpcport PORT
}
```

Only `[ZONES...]`, `chaos` and `grpcport` is different with plugin with *[kubernetes](https://coredns.io/plugins/kubernetes/)*:

* `[ZONES...]` defines which zones of the host will be treated as internal hosts in the Kubernetes cluster.

* `chaos` **ACTION** **SCOPE** **[PODS...]** set the behavior and scope of chaos. 

  Valid value for **Action**:

  * `random`: return random IP for DNS request.
  * `error`:  return error for DNS request.

  Valid value for **SCOPE**:
    
  * `inner`: chaos only works on the inner host of the Kubernetes cluster.
  * `outer`: chaos only works on the outer host of the Kubernetes cluster.
  * `all`:   chaos works on all the hosts.

  **[PODS...]** defines which Pods will take effect, the format is `Namespace`.`PodName`.

* `grpcport` **PORT** sets the port of GRPC service, which is used for the hot update of the chaos rules. The default value is `9288`. The interface of the GRPC service is defined in [dns.proto](pb/dns.proto).

## Examples

All DNS requests in Pod `busybox.busybox-0` will get error:

```yaml
    k8s_dns_chaos cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
        ttl 30
        chaos error all busybox.busybox-0
    }
```

The shell command below will execute failed:
  
```shell
    kubectl exec busybox-0 -it -n busybox -- ping -c 1 google.com
    ping: bad address 'google.com'
```