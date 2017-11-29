+++
date = "2017-03-01T20:33:08Z"
description = "CoreDNS is now fully compliant with the Kubernetes DNS Service Discovery Specification."
tags = ["Kubernetes", "Service", "Discovery", "Kube-DNS", "Documentation"]
title = "CoreDNS for Kubernetes Service Discovery, Take 2"
author = "john"
+++

A couple months ago we published a blog post on [how to use CoreDNS instead of
kube-dns](https://community.infoblox.com/t5/Community-Blog/CoreDNS-for-Kubernetes-Service-Discovery/ba-p/8187)
in Kubernetes. Since then, we have made a lot of progress. We worked with the
community to define a specification for [Kubernetes DNS-Based Service Discovery]
(https://github.com/kubernetes/dns/blob/master/docs/specification.md), enabling us to ensure
compatibility across the existing Kube-DNS implementation
and our new one in CoreDNS. Version 1.0.0 of this specification mostly follows
the current behavior of Kube-DNS. Versions 005 and higher of CoreDNS implement the full
specification and more.

At the time of that blog, the CoreDNS Kubernetes plugin only supported
serving `A` records with the cluster IP of an ordinary service. Now it includes:

* `A`, `SRV`, and `PTR` records for regular and headless services
* `A` records for named endpoints that are part of a service (i.e., records for
  'pets')
* `A` records for pods (optional) as described in the spec
* `TXT` record for discovering the DNS schema version in use

The pod `A` record support is not needed in all clusters, and is disabled by
default. Additionally, CoreDNS support for this use case goes beyond the
standard behavior you'll find in Kube-DNS. In Kube-DNS, these records do not
reflect the state of the cluster. Any query to `w-x-y-z.namespace.pod.cluster.local`
will return an `A` record with `w.x.y.z`, even if that IP does not belong to
specified namespace or even to the cluster address space. The original idea
behind this was to enable the use of wildcard SSL certificates for domains like
`*.namespace.pod.cluster.local`. CoreDNS can duplicate this behavior with the
configuration option `pods insecure`. We deem it "insecure" because this lack of
validation breaks the identity guarantee of the wildcard certificate.

The CoreDNS integration offers the option `pods verified`, which will verify
that the IP address `w.x.y.z` returned is in fact the IP of a pod in the
specified namespace. This prevents the spoofing of a DNS name in the namespace.
It does however, potentially increase the memory footprint of the CoreDNS instances
substantially, since now it needs to watch all pods, not just service endpoints.

Ok, enough intro - let's see how to deploy this. Like in the previous blog, we
use a ConfigMap and a Deployment. To make it even easier, we have created a little
[utility script](https://github.com/coredns/deployment/blob/master/kubernetes/deploy.sh) and
[deployment manifest template](https://github.com/coredns/deployment/blob/master/kubernetes/coredns.yaml.sed)
that you can use to deploy. To use it, simply put them both in the same directory
and then run the `deploy.sh` script, passing it the CIDR for your services. This will
generate the ConfigMap with the necessary Corefile. It will also lookup the cluster IP
of the existing kube-dns service. For example, running

    $ ./deploy.sh 10.3.0.0/24 cluster.local

with `kubectl` pointing to a cluster that has kube-dns service with cluster IP 10.3.0.10
results in the manifest file below.[^1]

[^1]: *Important:* If you are using Google Container Engine, there are additional processes that will not allow you to replace the kube-dns deployment (or replication controller). If you try to apply the above, it will successfully alter the Service but will soon kill your CoreDNS deployment and restart the Kube DNS replication controller. Since the service is updated, you actually will lose DNS in your cluster, as the service selector for the service no longer points to the pods created by the Kube DNS replication controller. There is probably a way around this but I haven't had a chance to find the fix yet.


~~~yaml
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
        kubernetes cluster.local 10.3.0.0/24
        proxy . /etc/resolv.conf
        cache 30
    }
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: coredns
  namespace: kube-system
  labels:
    k8s-app: coredns
    kubernetes.io/cluster-service: "true"
    kubernetes.io/name: "CoreDNS"
spec:
  replicas: 1
  selector:
    matchLabels:
      k8s-app: coredns
  template:
    metadata:
      labels:
        k8s-app: coredns
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ''
        scheduler.alpha.kubernetes.io/tolerations: '[{"key":"CriticalAddonsOnly", "operator":"Exists"}]'
    spec:
      containers:
      - name: coredns
        image: coredns/coredns:latest
        imagePullPolicy: Always
        args: [ "-conf", "/etc/coredns/Corefile" ]
        volumeMounts:
        - name: config-volume
          mountPath: /etc/coredns
        ports:
        - containerPort: 53
          name: dns
          protocol: UDP
        - containerPort: 53
          name: dns-tcp
          protocol: TCP
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 60
          timeoutSeconds: 5
          successThreshold: 1
          failureThreshold: 5
      dnsPolicy: Default
      volumes:
        - name: config-volume
          configMap:
            name: coredns
            items:
            - key: Corefile
              path: Corefile
---
apiVersion: v1
kind: Service
metadata:
  name: kube-dns
  namespace: kube-system
  labels:
    k8s-app: coredns
    kubernetes.io/cluster-service: "true"
    kubernetes.io/name: "CoreDNS"
spec:
  selector:
    k8s-app: coredns
  clusterIP: 10.3.0.10
  ports:
  - name: dns
    port: 53
    protocol: UDP
  - name: dns-tcp
    port: 53
    protocol: TCP
~~~

Let's take a closer look at that Corefile. It's really very similar to the one in the
previous blog post.

~~~ corefile
.:53 {
    errors
    log
    health
    kubernetes cluster.local 10.3.0.0/24
    proxy . /etc/resolv.conf
    cache 30
}
~~~

The one difference though, is the additions of `10.3.0.0/24` zone. This tells the
Kubernetes plugin that it is responsible for serving `PTR` requests for the reverse
zone `0.3.10.in-addr.arpa.`. In other words, this is what allows reverse DNS
resolution of services. Let's give that a try:

~~~txt
$ ./deploy.sh 10.3.0.0/24 | kubectl apply -f -
configmap "coredns" created
deployment "coredns" created
service "kube-dns" configured
$ kubectl run -it --rm --restart=Never --image=infoblox/dnstools:latest dnstools
Waiting for pod default/dnstools to be running, status is Pending, pod ready: false
If you don't see a command prompt, try pressing enter.
/ # host kubernetes
kubernetes.default.svc.cluster.local has address 10.3.0.1
/ # host kube-dns.kube-system
kube-dns.kube-system.svc.cluster.local has address 10.3.0.10
/ # host 10.3.0.1
1.0.3.10.in-addr.arpa domain name pointer kubernetes.default.svc.cluster.local.
/ # host 10.3.0.10
10.0.3.10.in-addr.arpa domain name pointer kube-dns.kube-system.svc.cluster.local.
/ #
~~~

Great, so it's working. Now, let's take a look at our CoreDNS logs:

~~~txt
$ kubectl get --namespace kube-system pods
NAME                                    READY     STATUS    RESTARTS   AGE
coredns-3558181428-0zhnh                1/1       Running   0          2m
coredns-3558181428-xri9i                1/1       Running   0          2m
heapster-v1.2.0-4088228293-a8gkc        2/2       Running   0          126d
kube-apiserver-10.222.243.77            1/1       Running   2          126d
kube-controller-manager-10.222.243.77   1/1       Running   2          126d
kube-proxy-10.222.243.77                1/1       Running   2          126d
kube-proxy-10.222.243.78                1/1       Running   0          126d
kube-scheduler-10.222.243.77            1/1       Running   2          126d
kubernetes-dashboard-v1.4.1-gi2xr       1/1       Running   0          24d
tiller-deploy-3299276078-e8phb          1/1       Running   0          24d
$ kubectl logs --namespace kube-system coredns-3558181428-0zhnh
2017/02/23 14:48:29 [INFO] Kubernetes plugin configured without a label selector. No label-based filtering will be performed.
.:53
2017/02/23 14:48:29 [INFO] CoreDNS-005
CoreDNS-005
10.2.6.127 - [23/Feb/2017:14:49:44 +0000] "AAAA IN kubernetes.default.svc.cluster.local. udp 54 false 512" NOERROR 107 544.128µs
10.2.6.127 - [23/Feb/2017:14:49:44 +0000] "MX IN kubernetes.default.svc.cluster.local. udp 54 false 512" NOERROR 107 7.576897ms
10.2.6.127 - [23/Feb/2017:14:49:52 +0000] "A IN kube-dns.kube-system.default.svc.cluster.local. udp 64 false 512" NXDOMAIN 117 471.176µs
23/Feb/2017:14:49:52 +0000 [ERROR 0 kube-dns.kube-system.default.svc.cluster.local. A] no items found
10.2.6.127 - [23/Feb/2017:14:50:00 +0000] "PTR IN 10.0.3.10.in-addr.arpa. udp 40 false 512" NOERROR 92 752.956µs
$ kubectl logs --namespace kube-system coredns-3558181428-xri9i
2017/02/23 14:48:29 [INFO] Kubernetes plugin configured without a label selector. No label-based filtering will be performed.
.:53
2017/02/23 14:48:29 [INFO] CoreDNS-005
CoreDNS-005
10.2.6.127 - [23/Feb/2017:14:49:44 +0000] "A IN kubernetes.default.svc.cluster.local. udp 54 false 512" NOERROR 70 1.10732ms
10.2.6.127 - [23/Feb/2017:14:49:52 +0000] "A IN kube-dns.kube-system.svc.cluster.local. udp 56 false 512" NOERROR 72 409.74µs
10.2.6.127 - [23/Feb/2017:14:49:52 +0000] "AAAA IN kube-dns.kube-system.svc.cluster.local. udp 56 false 512" NOERROR 109 210.817µs
10.2.6.127 - [23/Feb/2017:14:49:52 +0000] "MX IN kube-dns.kube-system.svc.cluster.local. udp 56 false 512" NOERROR 109 796.703µs
10.2.6.127 - [23/Feb/2017:14:49:56 +0000] "PTR IN 1.0.3.10.in-addr.arpa. udp 39 false 512" NOERROR 89 694.649µs
$
~~~

Here, we can see that the queries were load-balanced across the two CoreDNS replicas. In production, it's a good idea
to disable the logging of all the queries, as logging slows down the query dramatically (often by an order of magnitude).
To do that, simply remove the `log` line from the Corefile.

That's all there is to it. Please feel free to submit any problems, questions, or feature ideas as issues on the
[CoreDNS GitHub](https://github.com/coredns/coredns).
