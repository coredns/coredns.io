+++
title = "idetcd"
description = "*idetcd* - etcd-based CoreDNS plugin used for identifying nodes in a cluster without domain name collision."
weight = 10
tags = [  "plugin" , "idetcd" ]
categories = [ "plugin", "external" ]
date = "2018-08-13T18:54:20+01:00"
repo = "https://github.com/jiachengxu/idetcd"
home = "https://github.com/jiachengxu/idetcd/blob/master/README.md"
+++

## Description

*idetcd* is used for identifying nodes in a cluster without domain name collision.The basic idea is quite simple: Set up CoreDNS server on every node when you going to start a cluster, and node exposes itself by taking the free domain name in etcd.

## Syntax

~~~
idetcd {
	endpoint ENDPOINT...
	limit LIMIT
	pattern PATTERN
}
~~~

* **endpoint** defines the etcd endpoints. Defaults to "http://localhost:2379".
* **limit** defines the maximum limit of the node number in the cluster, if some nodes is going to expose itself after the node number in the cluster hits this limit, it will fail.
* **pattern** defines the domain name pattern that every node follows in the cluster. And here we use golang template for the pattern.


## Examples
In the following example, we are going to start up a cluster which contains 5 nodes, on every node we can get this project by:

```
$ go get -u github.com/jiachengxu/idetcd
```

Before you move to the next step, make sure that you've **already set up a etcd instance**, and don't forget to write down the endpoints.

Then you need to add a Corefile which specifys the configuration of the CoreDNS server in the same directory of `main.go`, a simple Corefile example is as follows, please go to [CoreDNS Github repo](https://github.com/coredns/coredns) for more details.

 ~~~ corefile
 . {
     idetcd {
         endpoint ETCDENDPOINTS
         limit 5
         pattern worker{{.ID}}.tf.local.
     }
 }
 ~~~

And then you can generate binary file by:
```sh
$ go build -v -o coredns
```

Alternatively, if you have docker installed, you could also execute the following to build:
```sh
$ docker run --rm -i -t -v $PWD:/go/src/github.com/jiachengxu/idetcd \
      -w /go/src/github.com/jiachengxu/idetcd golang:1.10 go build -v -o coredns
```

Then run it by:
```sh
$ ./coredns
```

After that, all nodes in the cluster are trying to find free slots in the etcd to expose themselves, once they succeed, you can get the domain name of every node on every node in the same cluster by:
```
$ dig +short worker4.tf.local @localhost
```
Also ipv6 is supported:
```
$ dig +short worker4.tf.local AAAA @localhost
```

## Integration with AWS
Using CoreDNS with idetcd plugin to config the cluster is a one-time process which is different with the general config process. For example, if you want to set up a cluster which contains several instances on AWS, you can use the same configuration for every instance and let all the instances to expose themselves in the `init` process. This can be achieved by using [`cloud-init`](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-cloud-init) in [`user data`](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html). Here is a bash script example for AWS instances to execute at launch:

```bash
#!/bin/bash
set -x
## Install docker.
yum install -y docker
echo
chkconfig docker on
service docker start
echo
## Install git.
yum install -y git
git clone https://github.com/jiachengxu/idetcd.git /home/ec2-user/idetcd
cd /home/ec2-user/idetcd
## Using docker to build the binary file of CoreDns with idetcd plugin specified.
docker run --rm -v $PWD:/go/src/github.com/jiachengxu/idetcd -w /go/src/github.com/jiachengxu/idetcd golang:1.10 go build -v -o coredns
## Create a Corefile for specifying the configuration of CoreDNS.(Don't forget to replace the ETCDENDPOINTS and NUMBER with your own etcd endpoints and limit of node in the cluster!)
cat > Corefile << EOF
. {
    idetcd {
        endpoint ETCDENDPOINTS
        limit NUMBER
        pattern worker{{.ID}}.tf.local.
    }
}
EOF
./coredns
```
