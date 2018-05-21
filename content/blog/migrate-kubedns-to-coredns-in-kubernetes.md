+++
title = "Migration from kube-dns to CoreDNS"
description = "A guide to migration from kube-dns to CoreDNS in Kubernetes"
tags = ["Kubernetes", "Service", "Discovery", "Kube-DNS", "DNS", "Migration", "Documentation"]
date = "2018-05-21T20:23:43-00:00"
author = "sandeep"
+++

CoreDNS is currently a Beta feature in Kubernetes and on course to being graduated to [General Availability (GA) for Kubernetes 1.11](https://github.com/kubernetes/community/pull/1956).
This means that CoreDNS will be available as a standard in Kubernetes via the installation toolkits such as kubeadm, kube-up, minikube and kops.

This document will guide you to migrating the DNS service from CoreDNS to kube-dns when using the various tools available to spin up a Kubernetes cluster.

## Installing CoreDNS via Kubeadm

There is an extensive guide on how to install CoreDNS instead of kube-dns via Kubeadm available [here](https://coredns.io/2018/01/29/deploying-kubernetes-with-coredns-using-kubeadm).
From Kubernetes v1.10, CoreDNS supports the translation of the kube-dns ConfigMap to CoreDNS ConfigMap.
That is, if you had configured `stubdomains`, `upstreamnameservers` and `federation` via the kube-dns ConfigMap, it will now be translated automatically to the equivalent CoreDNS ConfigMap during when choosing to install CoreDNS using `kubeadm upgrade`.

`Stubdomain` and `upstreamnameserver` in kube-dns translates to the [`proxy`](https://coredns.io/plugins/proxy/) in CoreDNS.
The `federation` in kube-dns has an equivalent [`federation`](https://coredns.io/plugins/federation/) in CoreDNS.


Example ConfigMap of kube-dns.
~~~text
apiVersion: v1
data:
  federations: |
    {"foo" : "foo.feddomain.com"}
  stubDomains: |
    {"abc.com" : ["1.2.3.4"], "my.cluster.local" : ["2.3.4.5"]}
  upstreamNameservers: |
    ["8.8.8.8", "8.8.4.4"]
kind: ConfigMap
metadata:
  creationTimestamp: 2018-01-22T20:21:56Z
  name: kube-dns
  namespace: kube-system
~~~

CoreDNS Corefile after translation.

~~~text
   .:53 {
        errors
        health
        kubernetes cluster.local  in-addr.arpa ip6.arpa {
           upstream  8.8.8.8 8.8.4.4
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
        }
        federation cluster.local {
           foo foo.feddomain.com
        }
        prometheus :9153
        proxy .  8.8.8.8 8.8.4.4
        cache 30
    }
    abc.com:53 {
        errors
        cache 30
        proxy . 1.2.3.4
    }
    my.cluster.local:53 {
        errors
        cache 30
        proxy . 2.3.4.5
    }
~~~



## Installing CoreDNS via Minikube.

CoreDNS is available in the `addon manager` and is disabled by default.

~~~text
$ minikube addons list
- kube-dns: enabled
- registry: disabled
- registry-creds: disabled
- freshpod: disabled
- addon-manager: enabled
- dashboard: enabled
- coredns: disabled
- heapster: disabled
- efk: disabled
- ingress: disabled
- default-storageclass: enabled
- storage-provisioner: enabled
~~~

To enable CoreDNS, run the following command:
> NOTE: Be sure to disable kube-dns after enabling CoreDNS. Otherwise, if both CoreDNS and kube-dns are running, queries may randomly hit either CoreDNS or kube-dns.

~~~
$ minikube addons enable coredns
coredns was successfully enabled
~~~


## CoreDNS in kube-up

[Kube-up](https://kubernetes.io/docs/getting-started-guides/scratch/) is another way to start a Kubernetes cluster, now mostly used for deploying Kubernetes in GCE for end-to-end (e2e) testing purposes.
The environment variable `ENABLE_CLUSTER_DNS` (default=true) is required to install DNS service.
For CoreDNS can be installed as the default DNS service, the environment variable `CLUSTER_DNS_CORE_DNS` needs to be set to `true`.

## CoreDNS in Kops

Currently, Kops v1.10 is set to include CoreDNS as an option to be installed instead of kube-dns.
In order to install CoreDNS in place of kube-dns, we need to specify the `provider` as `CoreDNS` in the [cluster yaml configuration for Kops](https://github.com/kubernetes/kops/blob/master/docs/cluster_spec.md).

~~~text
spec:
  kubeDNS:
    provider: CoreDNS
~~~
This will install CoreDNS instead of kube-dns.

## Installing CoreDNS via other methods

For users keen to install CoreDNS in place of kube-dns but who are not using `kubeadm`, `minikube`, `kube-up`, or `kops`, there are instructions in the [CoreDNS deployment repository](https://github.com/coredns/deployment/tree/master/kubernetes), which will help you to migrate from kube-dns to CoreDNS.
Users should delete the kube-dns deployment after deploying CoreDNS. Otherwise, if both CoreDNS and kube-dns are running, queries may randomly hit either CoreDNS or kube-dns.
