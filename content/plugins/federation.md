+++
title = "federation"
description = "*federation* enables federated queries to be resolved via the kubernetes plugin."
weight = 12
tags = [ "plugin", "federation" ]
categories = [ "plugin" ]
date = "2018-10-17T18:39:57.646722"
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
    upstream [ADDRESS...]
}
~~~

* Each **NAME** and **DOMAIN** defines federation membership. One entry for each. A duplicate
  **NAME** will silently overwrite any previous value.
* `upstream` [**ADDRESS**...] defines the upstream resolvers used for resolving the `CNAME` target
  produced by this plugin.  If no **ADDRESS** is given, CoreDNS
  will resolve External Services against itself. **ADDRESS** can be an IP, an IP:port, or a path
  to a file structured like resolv.conf.

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
