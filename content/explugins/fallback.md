+++
title = "fallback"
description = "*fallback* - allow redirecting queries to an alternate set of upstreams based on RCODE"
weight = 10
tags = [  "plugin" , "fallback" ]
categories = [ "plugin", "external" ]
date = "2018-02-07T12:36:00-05:00"
repo = "https://github.com/coredns/fallback"
home = "https://github.com/coredns/fallback/blob/master/README.md"
+++

## Description

The *fallback* plugin allows an alternate set of upstreams be specified which will be used
if the plugin chain returns specific error messages. The *fallback* plugin utilizes the *proxy* plugin (<https://coredns.io/plugins/proxy>) to query the specified upstreams.

As the name suggests, the purpose of the *fallback* is to allow a fallback when, for example,
the desired upstreams became unavailable.

It is recommended that this comes before *proxy* and *forward* in the `plugins.cfg` file. It must
come before any plugins for which it wants to react to their response code.

## Syntax

```
fallback RCODE PROXY_PARAMS
```

* **RCODE** is the string representation of the error response code. The complete list of valid rcode strings are defined as `RcodeToString` in <https://github.com/miekg/dns/blob/master/msg.go>, examples of which are `SERVFAIL`, `NXDOMAIN` and `REFUSED`.
* **PROXY_PARAMS** accepts the same parameters as the *proxy* plugin
<https://coredns.io/plugins/proxy>.

## Examples

### Fallback to local DNS server

The following specifies that all requests are proxied to 8.8.8.8. If the response is `NXDOMAIN`, *fallback* will proxy the request to 192.168.1.1:53, and reply to client accordingly.

```
. {
	proxy . 8.8.8.8
	fallback NXDOMAIN . 192.168.1.1:53
	log
}

```

### Multiple fallbacks

Multiple fallbacks can be specified, as long as they serve unique error responses.

```
. {
    proxy . 8.8.8.8
    fallback NXDOMAIN . 192.168.1.1:53
    fallback REFUSED . 192.168.100.1:53
    log
}

```

### Additional proxy parameters

You can specify additional proxy parameters for each of the fallback upstreams.

```
. {
    proxy . 8.8.8.8
    fallback NXDOMAIN . 192.168.1.1:53 192.168.1.2:53 {
        protocol dns force_tcp
    }
    log
}
```
