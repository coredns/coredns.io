+++
title = "federation"
description = "The *federation* plugin enables [federated](https://kubernetes.io/docs/tasks/federation/federation-service-discovery/) queries to be resolved via the kubernetes plugin."
weight = 12
tags = [ "plugin", "federation" ]
categories = [ "plugin" ]
date = "2017-09-10T18:11:52.763770"
+++

Enabling *federation* without also having *kubernetes* is a noop.

## Syntax

~~~
federation [ZONES...] {
    NAME DOMAIN
~~~

* Each **NAME** and **DOMAIN** defines federation membership. One entry for each. A duplicate
  **NAME** will silently overwrite any previous value.

## Examples

Here we handle all service requests in the `prod` and `stage` federations.

~~~ txt
. {
    kubernetes cluster.local 
    federation cluster.local {
        prod prod.feddomain.com
        staging staging.feddomain.com
    }
}
~~~

Or slightly shorter:

~~~ txt
cluster.local {
    kubernetes
    federation {
        prod prod.feddomain.com
        staging staging.feddomain.com
    }
}
~~~
