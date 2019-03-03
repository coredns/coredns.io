+++
title = "federation"
description = "*federation* enables federated queries to be resolved via the kubernetes plugin."
weight = 12
tags = [ "plugin", "federation" ]
categories = [ "plugin" ]
date = "2019-03-03T09:28:16.706382"
+++

## Description

Enabling this plugin allows
[Federated](https://kubernetes.io/docs/tasks/federation/federation-service-discovery/) queries to be
resolved via the kubernetes plugin.

Enabling *federation* without also having *kubernetes* is a noop.

## Syntax

~~~
federation [ZONES...] {
    NAME DOMAIN
    upstream
}
~~~

* Each **NAME** and **DOMAIN** defines federation membership. One entry for each. A duplicate
  **NAME** will silently overwrite any previous value.
* `upstream` [**ADDRESS**...] resolve the `CNAME` target produced by this plugin.  CoreDNS
  will resolve External Services against itself.

## Examples

Here we handle all service requests in the `prod` and `stage` federations.

~~~
. {
    kubernetes cluster.local
    federation cluster.local {
        prod prod.feddomain.com
        staging staging.feddomain.com
        upstream
    }
}
~~~

Or slightly shorter:

~~~
cluster.local {
    kubernetes
    federation {
        prod prod.feddomain.com
        staging staging.feddomain.com
        upstream
    }
}
~~~
