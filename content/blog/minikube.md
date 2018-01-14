+++
date = "2017-04-28T12:30:00Z"
description = "Getting CoreDNS to work with Minikube."
tags = ["Kubernetes", "Service", "Discovery", "Kube-DNS", "Minikube", "DNS", "Documentation"]
title = "CoreDNS for Minikube"
author = "john"
+++

In our [previous post](/2017/03/01/coredns-for-kubernetes-service-discovery-take-2/),
we showed how CoreDNS can be used in place of Kube-DNS for service discovery in Kubernetes clusters.
In that blog, there is a footnote about issues trying to replace Kube-DNS when using Google Container Engine (GKE).
As it so happens, there is a similar issue with [minikube](https://github.com/kubernetes/minikube), which is a local
Kubernetes environment that is very useful for developers.

When you try to replace Kube-DNS, you will find that shortly after you modify its service to point to CoreDNS, your
changes will be reverted. This is because Minikube has an _addon manager_ that periodically verifies the configuration
state of any installed addons, and Kube-DNS is one of those addons.

Luckily, this is really easy to solve for Minikube. The `minikube` command allows you to modify the installed
plugins for Minikube. So, we simply need to disable the `kube-dns` addon before running our `kubectl apply -f`
that is described in the previous blog:

~~~ txt
$ minikube addons list
- dashboard: enabled
- default-storageclass: enabled
- kube-dns: enabled
- heapster: disabled
- ingress: disabled
- registry-creds: disabled
- addon-manager: enabled
$ minikube addons disable kube-dns
kube-dns was successfully disabled
$ minikube addons list
- heapster: disabled
- ingress: disabled
- registry-creds: disabled
- addon-manager: enabled
- dashboard: enabled
- default-storageclass: enabled
- kube-dns: disabled
~~~

Now we can apply the `coredns.yaml`, and also delete the `kube-dns` `ReplicationController` which
will not be automatically deleted by disabling the addon.

~~~ txt
$ kubectl apply -f coredns.yaml
serviceaccount "coredns" configured
clusterrole "system:coredns" configured
clusterrolebinding "system:coredns" configured
configmap "coredns" configured
deployment "coredns" configured
service "kube-dns" configured
$ kubectl get -n kube-system pods
NAME                          READY     STATUS    RESTARTS   AGE
coredns-980047985-g2748       1/1       Running   1          36m
kube-addon-manager-minikube   1/1       Running   0          9d
kube-dns-v20-qzvr2            3/3       Running   0          1m
kubernetes-dashboard-ks1jp    1/1       Running   0          9d
$ kubectl delete -n kube-system rc kube-dns-v20
replicationcontroller "kube-dns-v20" deleted
$
~~~

And there we have it, CoreDNS is up and running and won't be overwritten by the addon manager.
