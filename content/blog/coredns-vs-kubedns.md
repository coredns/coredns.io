
+++ title = "Cluster DNS Resource Usage: CoreDNS vs Kube-DNS" description = "Comparing CoreDNS and Kube-DNS resource requirements" tags = ["Kubernetes", "Service", "Discovery", "DNS", "Documentation", "Deployment"] date = "2018-11-27T00:00:00-00:00" author = "chris" +++


When compiling data for a [resource deployment guide for CoreDNS](https://coredns.io/2018/11/15/scaling-coredns-in-kubernetes-clusters/) a few weeks ago, I also collected the same data for kube-dns using the same test environments.  Although CoreDNS and Kube-dns ultimately perform the same task, there are some key differences in implementation that affect resource consumption and performance.  At a high level, some of these differences are:

* CoreDNS is a single container per instance, vs kube-dns which uses three.
* Kube-dns uses dnsmasq for caching, which is single threaded C.  CoreDNS is multi-threaded Go.
* CoreDNS enables negative caching in the default deployment. Kube-dns does not.  

These differences affect performance in various ways.  The larger number of containers per instance in kube-dns increases base memory requirements, and also adds some performance overhead (as requests/responses need to be passed back and forth between containers).  For kube-dns, dnsmasq may be highly optimized in C, but it's also single threaded so it can only use one core per instance. CoreDNS enables negative caching, which aids in handling external names searches.


## Memory

Both CoreDNS and kube-dns maintain a local cache of all Services and Endpoints in the cluster.  So as the number of Services and Endpoints scale up, so do the memory requirements for each DNS Pod. At default settings, CoreDNS should be expected to use less memory than kube-dns. This is in part due to the overhead of the three containers used by kube-dns, vs only one container in CoreDNS.

The chart below shows the estimated memory required to run a single instance of CoreDNS or Kube-dns based on the number of Services and Endpoints.

![CoreDNS vs Kube-DNS estimated memory at scale](https://docs.google.com/spreadsheets/d/e/2PACX-1vS7d2MlgN1gMrrOHXa7Zn6S3VqujST5L-4PHX7jr4IUhVcTi0guXVRCgtIYrtLm3qxZWFlMHT-Xt9n3/pubchart?oid=1145811570&format=image)

The sources of the above data are from Kubernetes e2e scale tests, in conjunction with small cluster QPS load tests.  The Kubernetes e2e scale tests provide testing on very large clusters, but do not apply any QPS load.  To account for additional memory needed while handling a QPS load, the chart adds in the memory deltas observed when applying maximal QPS load during the CPU tests (below).  This was about 58Mi for kube-dns, and 5Mi for CoreDNS.


## CPU

In terms of CPU performance, CoreDNS performs much better for external names (e.g. `infoblox.com`), and slightly worse for internal names (e.g. `kubernetes`). 

| DNS Server  | Query Type  | QPS    | Avg Latency (ms) |
|-------------|-------------|--------|------------------|
| CoreDNS     | external    | 6733   | 12.02            |
| CoreDNS     | internal    | 33669  | 2.608            |
| Kube-dns    | external    | 2227   | 41.585           |
| Kube-dns    | internal    | 36648  | 2.639            |


Take aways:

* Kube-dns performed about 10% better for internal names.  This is probably due to dnsmasq being more optimized than CoreDNS's built-in caching.
* CoreDNS performed about 3X better for external names. This is partly caused by negative responses not being cached in kube-dns deployments. However enabling negative cache in the kube-dns deployment did not significantly change the outcome, so the bulk of performance gain is elsewhere.

| DNS Server           | Query Type  | QPS    | Avg Latency (ms) |
|----------------------|-------------|--------|------------------|
| Kube-dns + neg-cache | external    | 2552   | 36.665           |
| Kube-dns + neg-cache | internal    | 28971  | 3.385            |

## More

The version of kube-dns and default configuration used in these tests were those released with Kubernetes 1.12.

For more details about the test environments see: [Scaling CoreDNS in Kubernetes Clusters] (https://github.com/coredns/deployment/blob/master/kubernetes/Scaling_CoreDNS.md).
