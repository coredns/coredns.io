+++ title = "Scaling CoreDNS in Kubernetes Clusters" description = "A guide for tuning CoreDNS resources/requirements in Kubernetes clusters" tags = ["Kubernetes", "Service", "Discovery", "DNS", "Documentation", "Deployment"] date = "2018-11-15T00:00:00-00:00" author = "chris" +++

I'm sharing the results of some tests I ran with CoreDNS (1.2.5) in Kubernetes (1.12) to provide some reference points for tuning CoreDNS to your cluster.
In addition to testing CoreDNS in its default configuration, I tested CoreDNS with the optional *autopath* plugin enabled.
The *autopath* plugin is an optimization that helps transparently mitigate the DNS performance penalties Pods incur due
to Kubernetes' infamous ndots:5 issue. These tests quantify the memory/performance trade when enabling *autopath*.

The guides and fomulas in this post are based on a set of tests of clusters in GCE, your mileage may vary.
This blog post is a excerpt of the complete results, you can see more detail [here](https://github.com/coredns/deployment/blob/master/kubernetes/Scaling_CoreDNS.md).

## Memory and Pods

In large scale Kubernetes clusters, CoreDNS's memory usage is predominantly affected by the number of Pods and Services in the cluster. 

![CoreDNS in Kubernetes Memory Use](https://docs.google.com/spreadsheets/d/e/2PACX-1vS7d2MlgN1gMrrOHXa7Zn6S3VqujST5L-4PHX7jr4IUhVcTi0guXVRCgtIYrtLm3qxZWFlMHT-Xt9n3/pubchart?oid=191775389&format=image)

### With default CoreDNS settings

To estimate the amount of memory required for a CoreDNS instance (using default settings), you can use the following formula:

>  MB required (default settings) = (Number of Pods + Services) / 1000 + 54

### With the *autopath* plugin

The *autopath* plugin is an optional optimization that improves performance for queries of names external to the cluster (e.g. `infoblox.com`). 
Enabling the *autopath* plugin requires CoreDNS to use significantly more memory to store information about Pods.  
Enabling the *autopath* plugin also puts additional load on the Kubernetes API, since it must monitor all changes to Pods.

To estimate the amount of memory required for a CoreDNS instance (using the *autopath* plugin), you can use the following formula:

>  MB required (w/ autopath) = (Number of Pods + Services) / 250 + 56

## CPU and QPS

Max QPS was tested by using the `kubernetes/perf-tests/dns` tool, on a cluster using CoreDNS. 
The two types of queries used were *internal queries* (e.g. `kubernetes`), and *external queries* (e.g. `infoblox.com`).  

### With default CoreDNS settings

Single instance of CoreDNS (default settings) on a GCE n1-standard-2 node:


| Query Type  | QPS              | Avg Latency (ms)   |
|-------------|------------------|--------------------|
| external    | 6733<sup>1</sup> | 12.02<sup>1</sup>  |
| internal    | 33669            | 2.608              |


<sup>1</sup> From the server perspective it is processing 33667 QPS with 2.404 ms latency, but from the client perspective,
each single name lookup actually comprised 5 serial lookups.

### With the *autopath* plugin

The *autopath* plugin in CoreDNS is an option that mitigates the ClusterFirst search list penalty. When enabled, it reduces the number of DNS queries a client makes when looking up an external name.  

Single instance of CoreDNS (with the *autopath* plugin enabled) on a GCE n1-standard-2 node:


| Query Type  | QPS   | Avg Latency (ms) |
|-------------|-------|------------------|
| external    | 31428 | 2.605            |
| internal    | 33918 | 2.62             |


Note that the numbers for external queries are much improved here.  This is due to the *autopath* plugin optimization.

The server perspective latency for external queries goes up slightly when *autopath* is enabled (+8%).  
This is because it's doing the extra work of checking each search domain on the server side.  
But since it can answer in one round trip instead of five, the overall client perspective performance is much improved.


### More...

For more information about the test environments and how the data was collected, see the full results [here](https://github.com/coredns/deployment/blob/master/kubernetes/Scaling_CoreDNS.md).

