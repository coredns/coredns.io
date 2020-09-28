+++
title = "alternate"
description = "*alternate* - allow redirecting queries to an alternate set of upstreams based on RCODE"
weight = 10
tags = [  "plugin" , "alternate" ]
categories = [ "plugin", "external" ]
date = "2020-09-28T00:00:00-05:00"
repo = "https://github.com/coredns/alternate"
home = "https://github.com/coredns/alternate/blob/master/README.md"
+++

## Description

The *alternate* plugin is able to selectively forward queries to another upstream server, depending the error result provided by the initial resolver.
It allows an alternate set of upstreams be specified which will be used
if the plugin chain returns specific error messages. The *alternate* plugin utilizes the *forward* plugin (<https://coredns.io/plugins/forward>) to query the specified upstreams.

> The *alternate* plugin supports only DNS protocol and random policy w/o additional *forward* parameters, so following directives will fail:

```
. {
    forward . 8.8.8.8
    alternate NXDOMAIN . tls://192.168.1.1:853 {
        policy sequential
    }
}
```

As the name suggests, the purpose of the *alternate* is to allow a alternate when, for example,
the desired upstreams became unavailable.

## Syntax

```
{
    alternate [original] RCODE_1[,RCODE_2,RCODE_3...] . DNS_RESOLVERS
}
```

* **original** is optional flag. If it is set then alternate uses original request instead of potentially changed by other plugins
* **RCODE** is the string representation of the error response code. The complete list of valid rcode strings are defined as `RcodeToString` in <https://github.com/miekg/dns/blob/master/msg.go>, examples of which are `SERVFAIL`, `NXDOMAIN` and `REFUSED`. At least one rcode is required, but multiple rcodes may be specified, delimited by commas.
* **DNS_RESOLVERS** accepts dns resolvers list.

## Examples

### Alternate to local DNS server

The following specifies that all requests are forwarded to 8.8.8.8. If the response is `NXDOMAIN`, *alternate* will forward the request to 192.168.1.1:53, and reply to client accordingly.

```
. {
	forward . 8.8.8.8
	alternate NXDOMAIN . 192.168.1.1:53
	log
}

```
### Alternate with original request used

The following specify that `original` query will be forwarded to 192.168.1.1:53 if 8.8.8.8 response is `NXDOMAIN`. `original` means no changes from next plugins on request. With no `original` flag alternate will forward request with EDNS0 option (set by rewrite).

```
. {
	forward . 8.8.8.8
	rewrite edns0 local set 0xffee 0x61626364
	alternate original NXDOMAIN . 192.168.1.1:53
	log
}

```

### Multiple alternates

Multiple alternates can be specified, as long as they serve unique error responses.

```
. {
    forward . 8.8.8.8
    alternate NXDOMAIN . 192.168.1.1:53
    alternate original SERVFAIL,REFUSED . 192.168.100.1:53
    log
}

```
