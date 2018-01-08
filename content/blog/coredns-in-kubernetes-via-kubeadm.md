+++
title = "Deploying Kubernetes with CoreDNS using kubeadm"
description = "A guide to installing CoreDNS in Kubernetes via kubeadm"
tags = ["Kubernetes", "Service", "Discovery", "Kube-DNS", "Kubeadm", "DNS", "Documentation"]
author = "sandeep"
+++

Kubernetes 1.9 has recently been launched and it ships with CoreDNS being a part of it. 
We can now install CoreDNS as the default service discovery via Kubeadm, which is the toolkit to install Kubernetes easily in a single step.

Currently, [CoreDNS is Alpha in Kubernetes 1.9](https://github.com/kubernetes/features/issues/427). We have a roadmap which will make CoreDNS Beta in version 1.10 and eventually be the default DNS, replacing kube-dns.

## Installing CoreDNS in fresh Kubernetes cluster
In order to install CoreDNS instead of kube-dns for a fresh Kubernetes cluster, we need to use the `feature-gates` flag and set it to `true`. 
Use the following command to install CoreDNS as default while installing a fresh Kubernetes cluster.
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

To start using your cluster, you need to run the following as a regular user:

  mkdir -p $HOME/.kube
  sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
  sudo chown $(id -u):$(id -g) $HOME/.kube/config

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at:
  https://kubernetes.io/docs/concepts/cluster-administration/addons/

You can now join any number of machines by running the following on each node
as root:

  kubeadm join --token 4cd282.a826a13b3705a4ec 147.75.107.43:6443 --discovery-token-ca-cert-hash sha256:9d98fd8463915998b3795f6ba53ae3db1fdc93ccbba6427bca1946a172ea1eb8

~~~

## Updating your existing cluster to use CoreDNS
In case you have an existing cluster, it is also possible to upgrade your DNS to CoreDNS, replacing kube-dns using the `kubeadm upgrade` command. 

Using `kubeadm upgrade plan` and by setting `feature-gates` flag as `true`, it is possible to check the CoreDNS version that will be installed before proceeding to apply the changes.

Checking the CoreDNS version to upgrade:

~~~ text
# kubeadm upgrade plan --allow-experimental-upgrades --feature-gates CoreDNS=true
[preflight] Running pre-flight checks.
[upgrade] Making sure the cluster is healthy:
[upgrade/health] FATAL: [preflight] Some fatal errors occurred:
	[ERROR MasterNodesReady]: there are NotReady masters in the cluster: [sandeep2]
[preflight] If you know what you are doing, you can make a check non-fatal with `--ignore-preflight-errors=...`
root@sandeep2:~# kubeadm upgrade plan --allow-experimental-upgrades --feature-gates CoreDNS=true --ignore-preflight-errors=all
[preflight] Running pre-flight checks.
[upgrade] Making sure the cluster is healthy:
	[WARNING MasterNodesReady]: there are NotReady masters in the cluster: [sandeep2]
[upgrade/config] Making sure the configuration is correct:
[upgrade/config] Reading configuration from the cluster...
[upgrade/config] FYI: You can look at this config file with 'kubectl -n kube-system get cm kubeadm-config -oyaml'
[upgrade] Fetching available versions to upgrade to
[upgrade/versions] Cluster version: v1.9.0
[upgrade/versions] kubeadm version: v1.9.0
[upgrade/versions] Latest stable version: v1.9.0
[upgrade/versions] Latest version in the v1.9 series: v1.9.0
[upgrade/versions] Latest experimental version: v1.10.0-alpha.1

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

The upgrade with CoreDNS as the default DNS can then be performed using `kubeadm upgrade apply` and `feature-gates`:

~~~ text
# kubeadm upgrade apply v1.10.0-alpha.1  --feature-gates CoreDNS=true  --allow-experimental-upgrades
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

To verify CoreDNS is running, we can check the pod status, deployment and the configmap in the node.

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
coredns   1                     1                       1                         1           4h
~~~

Check the `Configmap`:
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

>It is important to note that when switching from kube-dns to CoreDNS, the configurations that come with kube-dns (stubzones, federations...) will no longer exist and will switch to a default configuration in CoreDNS.
>The translation of the configurations from kube-dns shall be supported from the upcoming version of Kubernetes (v1.10) where CoreDNS will be Beta.

We can check if CoreDNS is functioning normally through a few basic `dig` commands:

~~~ text
# dig @10.96.0.10 kubernetes.default.svc.cluster.local

; <<>> DiG 9.10.3-P4-Ubuntu <<>> @10.96.0.10 kubernetes.default.svc.cluster.local
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 58006
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;kubernetes.default.svc.cluster.local. IN A

;; ANSWER SECTION:
kubernetes.default.svc.cluster.local. 5	IN A	10.96.0.1

;; Query time: 0 msec
;; SERVER: 10.96.0.10#53(10.96.0.10)
;; WHEN: Thu Dec 21 18:10:49 UTC 2017
;; MSG SIZE  rcvd: 81


# dig @10.96.0.10 ptr 1.0.96.10.in-addr.arpa.

; <<>> DiG 9.10.3-P4-Ubuntu <<>> @10.96.0.10 ptr 1.0.96.10.in-addr.arpa.
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 30847
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;1.0.96.10.in-addr.arpa.		IN	PTR

;; ANSWER SECTION:
1.0.96.10.in-addr.arpa.	5	IN	PTR	kubernetes.default.svc.cluster.local.

;; Query time: 0 msec
;; SERVER: 10.96.0.10#53(10.96.0.10)
;; WHEN: Thu Dec 21 18:09:10 UTC 2017
;; MSG SIZE  rcvd: 101

~~~

Below is the `log` for some of the basic queries.
~~~ text
# kubectl -n kube-system logs coredns-546545bc84-p4x7k
.:53
CoreDNS-1.0.1
linux/amd64, go1.9.1, a04eeb9c
2017/12/21 18:01:35 [INFO] CoreDNS-1.0.1
2017/12/21 18:01:35 [INFO] linux/amd64, go1.9.1, a04eeb9c
10.32.0.1 - [21/Dec/2017:18:06:34 +0000] "A IN kube-dns.kube-system.svc.cluster.local. udp 68 false 4096" NOERROR qr,aa,rd,ra 84 104.117µs
10.32.0.1 - [21/Dec/2017:18:07:48 +0000] "A IN kubernetes.default.svc.cluster.local. udp 66 false 4096" NOERROR qr,aa,rd,ra 82 116.92µs
10.32.0.1 - [21/Dec/2017:18:08:53 +0000] "PTR IN 10.0.96.10.in-addr.arpa. udp 53 false 4096" NOERROR qr,aa,rd,ra 105 97.984µs
10.32.0.1 - [21/Dec/2017:18:09:10 +0000] "PTR IN 1.0.96.10.in-addr.arpa. udp 52 false 4096" NOERROR qr,aa,rd,ra 102 94.992µs
10.32.0.1 - [21/Dec/2017:18:10:49 +0000] "A IN kubernetes.default.svc.cluster.local. udp 66 false 4096" NOERROR qr,aa,rd,ra 82 99.664µs
~~~

## Plugins
CoreDNS in Kubernetes ships with the following `plugins` enabled:
- *Error*: This enables error logging.
- *Health*: Health enables a simple health check endpoint.
- *Kubernetes*: The kubernetes plugin enables the reading zone data from a Kubernetes cluster. You can find more details [here](https://coredns.io/plugins/kubernetes/). 

> The `pods insecure` option always return an A record with IP from request (without checking k8s). This option is provided for backward compatibility with kube-dns.
> Also by default, the Kubernetes plugin has the `Cluster Domain` and the `Service CIDR` defined. The `Pod CIDR` must be added to the config file after CoreDNS deployment.
> `Upstream` option in the Kubernetes plugin defines upstream resolvers to be used resolve external names found (CNAMEs) pointing to external names.

- *Prometheus*: This enables [Prometheus](https://prometheus.io/) metrics.
- *Proxy*: Proxy facilitates both a basic reverse proxy and a robust load balancer.
- *Cache*: This enables a frontend cache. It will cache all records except zone transfers and metadata records.
