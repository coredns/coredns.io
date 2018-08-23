+++
date = "2017-05-08T20:01:00Z"
description = "Creating custom DNS entries inside or outside the cluster domain using CoreDNS."
tags = ["Kubernetes", "Service", "Discovery", "Kube-DNS", "Custom", "DNS", "Documentation"]
title = "Custom DNS Entries For Kubernetes"
author = "john"
+++

As described in our [previous post](/2017/03/01/coredns-for-kubernetes-service-discovery-take-2/),
CoreDNS can be used in place of Kube-DNS for service discovery in Kubernetes clusters. Because of the flexible architecture
of CoreDNS, this can enable some interesting use cases. In this blog, we'll show how to solve a common problem - creating
custom DNS entries for your services.

There are a couple of different possiblities here:

* [Making an alias](https://github.com/kubernetes/kubernetes/issues/39792) for an external name
* [Dynamically adding services to another domain](https://github.com/kubernetes/dns/issues/55), without running another server
* Adding an arbitrary entry inside the cluster domain

CoreDNS can solve all of these use cases. Let's start with the first one, which is pretty common. In this situation, you
want to be able to use the same name for a given service, whether you are accessing it inside or outside the cluster. This
is helpful, for example, when using TLS certificates that are bound to that name.

Suppose we have a service, `foo.default.svc.cluster.local` that is available to outside clients as `foo.example.com`.
That is, when looked up outside the cluster, `foo.example.com` will resolve to the load balancer VIP - the external
IP address for the service. Inside the cluster, it will resolve to the same thing, and so using this name internally
will cause traffic to hairpin - travel out of the cluster and then back in via the external IP. Instead, we want it
to resolve to the internal ClusterIP, avoiding the hairpin.

To do this in CoreDNS, we make use of the `rewrite` plugin. This plugin can modify a query before it is sent
down the chain to whatever backend is going to answer it. Recall the `Corefile` (CoreDNS configuration file) we
used in the last blog:

~~~ corefile
.:53 {
    errors
    log
    health
    kubernetes cluster.local 10.0.0.0/24
    proxy . /etc/resolv.conf
    cache 30
}
~~~

To get the behavior we want, we just need to add a rewrite rule mapping `foo.example.com` to `foo.default.svc.cluster.local`:

~~~ corefile
.:53 {
    errors
    log
    health
    rewrite name foo.example.com foo.default.svc.cluster.local
    kubernetes cluster.local 10.0.0.0/24
    proxy . /etc/resolv.conf
    cache 30
}
~~~

Once we add that to the `ConfigMap` via `kubectl edit` or `kubectl apply`, we have to let CoreDNS know that the `Corefile`
has changed. You can send it a `SIGUSR1` to tell it to reload graceful - that is, without loss of service:

~~~ txt
$ kubectl exec -n kube-system coredns-980047985-g2748 -- kill -SIGUSR1 1
~~~

Running our test pod, we can see this works:

~~~ txt
$ kubectl run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools
If you don't see a command prompt, try pressing enter.
/ # host foo
foo.default.svc.cluster.local has address 10.0.0.72
/ # host foo.example.com
foo.example.com has address 10.0.0.72
/ # host bar.example.com
Host bar.example.com not found: 3(NXDOMAIN)
/ #
~~~

That's all there is to solving the first problem.

The second problem is just as easy. Here, we just want to be able to serve DNS entries out of a different zone
than the cluster domain. Since CoreDNS is a general-purpose DNS server, there are many other ways
to serve up zones than just the `kubernetes` plugin. For simplicity, we'll use the `file` plugin along
with another `ConfigMap` entry to satisfy this use case. However, you could use the `etcd` plugin to store services
directly within an etcd instance, or the `auto` plugin to manage a set of zones (very nice when used along
with [git-sync](https://github.com/kubernetes/git-sync)).

To create the new zone, we need to modify the `coredns.yaml` we have been using to create an additional file
in the pod. To do this we have to edit the `ConfigMap` by adding a `file` line to the `Corefile`, and also
by adding another key, `example.db`, for the zone file:

~~~ yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: coredns
  namespace: kube-system
data:
  Corefile: |
    .:53 {
        errors
        log
        health
        rewrite name foo.example.com foo.default.svc.cluster.local
        kubernetes cluster.local 10.0.0.0/24
        file /etc/coredns/example.db example.org
        proxy . /etc/resolv.conf
        cache 30
    }
  example.db: |
    ; example.org test file
    example.org.            IN      SOA     sns.dns.icann.org. noc.dns.icann.org. 2015082541 7200 3600 1209600 3600
    example.org.            IN      NS      b.iana-servers.net.
    example.org.            IN      NS      a.iana-servers.net.
    example.org.            IN      A       127.0.0.1
    a.b.c.w.example.org.    IN      TXT     "Not a wildcard"
    cname.example.org.      IN      CNAME   www.example.net.

    service.example.org.    IN      SRV     8080 10 10 example.org.
~~~

and we also need to edit the `volumes` section of the `Pod` template spec:

~~~ yaml
      volumes:
        - name: config-volume
          configMap:
            name: coredns
            items:
            - key: Corefile
              path: Corefile
            - key: example.db
              path: example.db
~~~

Once we apply this using `kubectl apply -f`, a new CoreDNS pod will be built, because of the new file
in the volume. Later changes to the file won't require a new pod, just a graceful restart like we did before.
Let's take a look:

~~~ txt
$ kubectl run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools
If you don't see a command prompt, try pressing enter.
/ # host foo
foo.default.svc.cluster.local has address 10.0.0.72
/ # host foo.example.com
foo.example.com has address 10.0.0.72
/ # host example.org
example.org has address 127.0.0.1
/ #
~~~

Perfect! We can now edit that `ConfigMap` and send `SIGUSR1` any time we want to add entries to `example.org`. Of course,
as mentioned earlier, we could also use the `etcd` backend and avoid the hassle of modifying the `ConfigMap` and
sending the signal.

This brings us to the last problem. That one can be solved using the new support for `fallthrough` in the `kubernetes`
plugin. This functionality has been added in the recently released version 007 of CoreDNS - we'll come back with
another blog soon show how to use it.
