+++
title = "Deploying Kubernetes with CoreDNS using kubeadm"
description = "A guide to installing CoreDNS in Kubernetes via kubeadm"
tags = ["Kubernetes", "Service", "Discovery", "Kube-DNS", "Kubeadm", "DNS", "Documentation"]
date = "2018-01-29T10:12:43-00:00"
author = "sandeep"
+++

Kubernetes 1.9 has recently been launched and it ships with CoreDNS being a part of it. 
We can now install CoreDNS as the default service discovery via Kubeadm, which is the toolkit to install Kubernetes easily in a single step.

Currently, [CoreDNS is Alpha in Kubernetes 1.9](https://github.com/kubernetes/features/issues/427). We have a roadmap which will make CoreDNS Beta in version 1.10 and eventually be the default DNS, replacing kube-dns.

>It is important to note that currently when switching from kube-dns to CoreDNS, the configurations that come with kube-dns (stubzones, federations...) will no longer exist and will switch to a default configuration in CoreDNS.
>The translation of the configurations from kube-dns shall be supported from the upcoming version of Kubernetes (v1.10) where CoreDNS will be Beta.


## Understanding CoreDNS Configuration

This is the default CoreDNS configuration shipped with kubeadm. It is saved in a kubernetes configmap named `coredns`:
~~~ text
# kubectl -n kube-system get configmap coredns -oyaml
apiVersion: v1
data:
  Corefile: |
    .:53 {
        errors
        health
        kubernetes cluster.local 10.96.0.0/12 {
           pods insecure
           upstream /etc/resolv.conf
        }
        prometheus :9153
        proxy . /etc/resolv.conf
        cache 30
    }
kind: ConfigMap
metadata:
  creationTimestamp: 2017-12-21T12:55:15Z
  name: coredns
  namespace: kube-system
  resourceVersion: "161"
  selfLink: /api/v1/namespaces/kube-system/configmaps/coredns
  uid: 30bf0882-e64e-11e7-baf6-0cc47a8055d6
~~~
The Corefile part is the configuration of CoreDNS. 
This configuration is based on the following [plugins](https://coredns.io/plugins/) of CoreDNS:

* [errors](https://coredns.io/plugins/errors/): Errors are logged to stdout.
* [health](https://coredns.io/plugins/health/): Health of CoreDNS is reported to http://localhost:8080/health.
* [kubernetes](https://coredns.io/plugins/kubernetes/): CoreDNS will reply to DNS queries based on IP of the services and pods of Kubernetes. You can find more details [here](https://coredns.io/plugins/kubernetes/). 

> The Kubernetes plugin has its options `Cluster Domain` and `Service CIDR` defined as `cluster.local` and `10.96.0.0/12` respectively by default through kubeadm. We can modify and choose the desired values through the kubeadm `--service-dns-domain` and `--service-cidr` flags.

> The `pods insecure` option is provided for backward compatibility with kube-dns.

> `Upstream` is used for resolving services that point to external hosts (External Services).

* [prometheus](https://coredns.io/plugins/prometheus/): Metrics of CoreDNS are available at http://localhost:9153/metrics in [Prometheus](https://prometheus.io/) format.
* [proxy](https://coredns.io/plugins/proxy/): Any queries that are not within the cluster domain of Kubernetes will be forwarded to predefined resolvers (/etc/resolv.conf).
* [cache](https://coredns.io/plugins/cache/): This enables a frontend cache.

We can modify the default behavior by modifying this configmap. A restart of the CoreDNS pod is required for the changes to take effect. 

## Installing CoreDNS in fresh Kubernetes cluster
In order to install CoreDNS instead of kube-dns for a fresh Kubernetes cluster, we need to use the `feature-gates` flag and set it to `CoreDNS=true`. 
Use the following command to install CoreDNS as default DNS service while installing a fresh Kubernetes cluster.
~~~ text
# kubeadm init --feature-gates CoreDNS=true
~~~
~~~ text
# kubeadm init --feature-gates CoreDNS=true
[init] Using Kubernetes version: v1.9.0
[init] Using Authorization modes: [Node RBAC]
[preflight] Running pre-flight checks.
	[WARNING SystemVerification]: docker version is greater than the most recently validated version. Docker version: 17.09.1-ce. Max validated version: 17.03
	[WARNING FileExisting-crictl]: crictl not found in system path
[preflight] Starting the kubelet service
[certificates] Generated ca certificate and key.
[certificates] Generated apiserver certificate and key.
[certificates] apiserver serving cert is signed for DNS names [sandeep2 kubernetes kubernetes.default kubernetes.default.svc kubernetes.default.svc.cluster.local] and IPs [10.96.0.1 147.75.107.43]
[certificates] Generated apiserver-kubelet-client certificate and key.
[certificates] Generated sa key and public key.
[certificates] Generated front-proxy-ca certificate and key.
[certificates] Generated front-proxy-client certificate and key.
[certificates] Valid certificates and keys now exist in "/etc/kubernetes/pki"
[kubeconfig] Wrote KubeConfig file to disk: "admin.conf"
[kubeconfig] Wrote KubeConfig file to disk: "kubelet.conf"
[kubeconfig] Wrote KubeConfig file to disk: "controller-manager.conf"
[kubeconfig] Wrote KubeConfig file to disk: "scheduler.conf"
[controlplane] Wrote Static Pod manifest for component kube-apiserver to "/etc/kubernetes/manifests/kube-apiserver.yaml"
[controlplane] Wrote Static Pod manifest for component kube-controller-manager to "/etc/kubernetes/manifests/kube-controller-manager.yaml"
[controlplane] Wrote Static Pod manifest for component kube-scheduler to "/etc/kubernetes/manifests/kube-scheduler.yaml"
[etcd] Wrote Static Pod manifest for a local etcd instance to "/etc/kubernetes/manifests/etcd.yaml"
[init] Waiting for the kubelet to boot up the control plane as Static Pods from directory "/etc/kubernetes/manifests".
[init] This might take a minute or longer if the control plane images have to be pulled.
[apiclient] All control plane components are healthy after 31.502217 seconds
[uploadconfig] Storing the configuration used in ConfigMap "kubeadm-config" in the "kube-system" Namespace
[markmaster] Will mark node sandeep2 as master by adding a label and a taint
[markmaster] Master sandeep2 tainted and labelled with key/value: node-role.kubernetes.io/master=""
[bootstraptoken] Using token: 4cd282.a826a13b3705a4ec
[bootstraptoken] Configured RBAC rules to allow Node Bootstrap tokens to post CSRs in order for nodes to get long term certificate credentials
[bootstraptoken] Configured RBAC rules to allow the csrapprover controller automatically approve CSRs from a Node Bootstrap Token
[bootstraptoken] Configured RBAC rules to allow certificate rotation for all node client certificates in the cluster
[bootstraptoken] Creating the "cluster-info" ConfigMap in the "kube-public" namespace
[addons] Applied essential addon: CoreDNS
[addons] Applied essential addon: kube-proxy

Your Kubernetes master has initialized successfully!
~~~
CoreDNS install is confirmed if we see the following output while deploying Kubernetes.
~~~
[addons] Applied essential addon: CoreDNS
~~~

## Updating your existing cluster to use CoreDNS
In case you have an existing cluster, it is also possible to upgrade your DNS service to CoreDNS, replacing kube-dns, using the `kubeadm upgrade` command. 

It is possible to check the CoreDNS version that will be installed before proceeding to apply the changes by using `kubeadm upgrade plan` and by setting `feature-gates` flag as `CoreDNS=true`.

Checking the CoreDNS version to upgrade:
~~~ text
# kubeadm upgrade plan  --feature-gates CoreDNS=true
~~~
~~~ text
# kubeadm upgrade plan  --feature-gates CoreDNS=true
...

Components that must be upgraded manually after you have upgraded the control plane with 'kubeadm upgrade apply':
COMPONENT   CURRENT      AVAILABLE
Kubelet     1 x v1.9.0   v1.10.0-alpha.1

Upgrade to the latest experimental version:

COMPONENT            CURRENT   AVAILABLE
API Server           v1.9.0    v1.10.0-alpha.1
Controller Manager   v1.9.0    v1.10.0-alpha.1
Scheduler            v1.9.0    v1.10.0-alpha.1
Kube Proxy           v1.9.0    v1.10.0-alpha.1
CoreDNS              1.0.1     1.0.1
Etcd                 3.1.10    3.1.10

You can now apply the upgrade by executing the following command:

	kubeadm upgrade apply v1.10.0-alpha.1

Note: Before you can perform this upgrade, you have to update kubeadm to v1.10.0-alpha.1.

~~~

The upgrade of the cluster with CoreDNS as the default DNS can then be performed using `kubeadm upgrade apply` and `feature-gates CoreDNS=true`:
~~~ text
# kubeadm upgrade apply <version> --feature-gates CoreDNS=true 
~~~
~~~ text
# kubeadm upgrade apply v1.10.0-alpha.1  --feature-gates CoreDNS=true 
[preflight] Running pre-flight checks.
[upgrade] Making sure the cluster is healthy:
[upgrade/config] Making sure the configuration is correct:
[upgrade/config] Reading configuration from the cluster...
[upgrade/config] FYI: You can look at this config file with 'kubectl -n kube-system get cm kubeadm-config -oyaml'
[upgrade/version] You have chosen to change the cluster version to "v1.10.0-alpha.1"
[upgrade/versions] Cluster version: v1.10.0-alpha.1
[upgrade/versions] kubeadm version: v1.9.0
[upgrade/version] Found 1 potential version compatibility errors but skipping since the --force flag is set: 

	- Specified version to upgrade to "v1.10.0-alpha.1" is at least one minor release higher than the kubeadm minor release (10 > 9). Such an upgrade is not supported
[upgrade/prepull] Will prepull images for components [kube-apiserver kube-controller-manager kube-scheduler]
[upgrade/apply] Upgrading your Static Pod-hosted control plane to version "v1.10.0-alpha.1"...
[upgrade/staticpods] Writing new Static Pod manifests to "/etc/kubernetes/tmp/kubeadm-upgraded-manifests781134294"
[controlplane] Wrote Static Pod manifest for component kube-apiserver to "/etc/kubernetes/tmp/kubeadm-upgraded-manifests781134294/kube-apiserver.yaml"
[controlplane] Wrote Static Pod manifest for component kube-controller-manager to "/etc/kubernetes/tmp/kubeadm-upgraded-manifests781134294/kube-controller-manager.yaml"
[controlplane] Wrote Static Pod manifest for component kube-scheduler to "/etc/kubernetes/tmp/kubeadm-upgraded-manifests781134294/kube-scheduler.yaml"
[upgrade/staticpods] Moved new manifest to "/etc/kubernetes/manifests/kube-apiserver.yaml" and backed up old manifest to "/etc/kubernetes/tmp/kubeadm-backup-manifests038673725/kube-apiserver.yaml"
[upgrade/staticpods] Waiting for the kubelet to restart the component
[apiclient] Found 1 Pods for label selector component=kube-apiserver
[upgrade/staticpods] Component "kube-apiserver" upgraded successfully!
[upgrade/staticpods] Moved new manifest to "/etc/kubernetes/manifests/kube-controller-manager.yaml" and backed up old manifest to "/etc/kubernetes/tmp/kubeadm-backup-manifests038673725/kube-controller-manager.yaml"
[upgrade/staticpods] Waiting for the kubelet to restart the component
[apiclient] Found 1 Pods for label selector component=kube-controller-manager
[upgrade/staticpods] Component "kube-controller-manager" upgraded successfully!
[upgrade/staticpods] Moved new manifest to "/etc/kubernetes/manifests/kube-scheduler.yaml" and backed up old manifest to "/etc/kubernetes/tmp/kubeadm-backup-manifests038673725/kube-scheduler.yaml"
[upgrade/staticpods] Waiting for the kubelet to restart the component
[apiclient] Found 1 Pods for label selector component=kube-scheduler
[upgrade/staticpods] Component "kube-scheduler" upgraded successfully!
[uploadconfig] Storing the configuration used in ConfigMap "kubeadm-config" in the "kube-system" Namespace
[bootstraptoken] Configured RBAC rules to allow Node Bootstrap tokens to post CSRs in order for nodes to get long term certificate credentials
[bootstraptoken] Configured RBAC rules to allow the csrapprover controller automatically approve CSRs from a Node Bootstrap Token
[bootstraptoken] Configured RBAC rules to allow certificate rotation for all node client certificates in the cluster
[addons] Applied essential addon: CoreDNS
[addons] Applied essential addon: kube-proxy

[upgrade/successful] SUCCESS! Your cluster was upgraded to "v1.10.0-alpha.1". Enjoy!

[upgrade/kubelet] Now that your control plane is upgraded, please proceed with upgrading your kubelets in turn.

~~~

## Verifying CoreDNS Service
To verify CoreDNS is running, we can check the pod status and deployment in the node.
Note here that the CoreDNS service will remain as "kube-dns" which ensures a smooth transition while upgrading your service discovery from kube-dns to CoreDNS.

Check `pod` status:
~~~ text
# kubectl -n kube-system get pods -o wide
NAME                               READY     STATUS    RESTARTS   AGE       IP              NODE
coredns-546545bc84-ttsh4           1/1       Running   0          5h        10.32.0.61      sandeep2
etcd-sandeep2                      1/1       Running   0          5h        147.75.107.43   sandeep2
kube-apiserver-sandeep2            1/1       Running   0          4h        147.75.107.43   sandeep2
kube-controller-manager-sandeep2   1/1       Running   0          4h        147.75.107.43   sandeep2
kube-proxy-fkxmg                   1/1       Running   0          4h        147.75.107.43   sandeep2
kube-scheduler-sandeep2            1/1       Running   4          5h        147.75.107.43   sandeep2
weave-net-jhjtc                    2/2       Running   0          5h        147.75.107.43   sandeep2
~~~

Check `Deployment`:
~~~ text
# kubectl -n kube-system get deployments
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
coredns   1            1          1           1        4h
~~~

We can check if CoreDNS is functioning normally through a few basic `dig` commands:

~~~ text
# dig @10.32.0.61 kubernetes.default.svc.cluster.local +noall +answer

; <<>> DiG 9.10.3-P4-Ubuntu <<>> @10.32.0.61 kubernetes.default.svc.cluster.local +noall +answer
; (1 server found)
;; global options: +cmd
kubernetes.default.svc.cluster.local. 23 IN A	10.96.0.1

# dig @10.32.0.61 ptr 1.0.96.10.in-addr.arpa. +noall +answer

; <<>> DiG 9.10.3-P4-Ubuntu <<>> @10.32.0.61 ptr 1.0.96.10.in-addr.arpa. +noall +answer
; (1 server found)
;; global options: +cmd
1.0.96.10.in-addr.arpa.	30	IN	PTR	kubernetes.default.svc.cluster.local.
~~~
