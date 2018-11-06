+++
title = "loop"
description = "*loop* detect simple forwarding loops and halt the server."
weight = 20
tags = [ "plugin", "loop" ]
categories = [ "plugin" ]
date = "2018-11-06T07:19:41.754524"
+++

## Description

The *loop* plugin will send a random probe query to ourselves and will then keep track of how many times
we see it. If we see it more than twice, we assume CoreDNS is looping and we halt the process.

The plugin will try to send the query for up to 30 seconds. This is done to give CoreDNS enough time
to start up. Once a query has been successfully sent *loop* disables itself to prevent a query of
death.

The query sent is `<random number>.<random number>.zone` with type set to HINFO.

## Syntax

~~~ txt
loop
~~~

## Examples

Start a server on the default port and load the *loop* and *forward* plugins. The *forward* plugin
forwards to it self.

~~~ txt
. {
    loop
    forward . 127.0.0.1
}
~~~

After CoreDNS has started it stops the process while logging:

~~~ txt
plugin/loop: Forwarding loop detected in "." zone. Exiting. See https://coredns.io/plugins/loop#troubleshooting. Probe query: "HINFO 5577006791947779410.8674665223082153551.".
~~~

## Limitations

This plugin only attempts to find simple static forwarding loops at start up time.  To detect a loop, all of the following must be true

* the loop must be present at start up time.
* the loop must occur for at least the `HINFO` query type.

## Troubleshooting

When CoreDNS logs contain the message `Forwarding loop detected ...`, this means that
the `loop` detection plugin has detected an infinite forwarding loop in one of the upstream
DNS servers.  This is a fatal error because operating with an infinite loop will consume
memory and CPU until eventual out of memory death by the host.

A forwarding loop is usually caused by:

* Most commonly, CoreDNS forwarding requests directly to itself. e.g. via a loopback address such as `127.0.0.1`, `::1` or `127.0.0.53`
* Less commonly, CoreDNS forwarding to an upstream server that in turn, forwards requests back to CoreDNS.

To troubleshoot this problem, look in your Corefile for any `proxy` or `forward` to the zone
in which the loop was detected.  Make sure that they are not forwarding to a local address or
to another DNS server that is forwarding requests back to CoreDNS. If `proxy` or `forward` are
 using a file (e.g. `/etc/resolv.conf`), make sure that file does not contain local addresses.

### Troubleshooting Loops In Kubernetes Clusters
When a CoreDNS Pod deployed in Kubernetes detects a loop, the CoreDNS Pod will start to "CrashLoopBackOff".
This is because Kubernetes will try to restart the Pod every time CoreDNS detects the loop and exits.

A common cause of forwarding loops in Kubernetes clusters is an interaction with
`systemd-resolved` on the host node.  `systemd-resolved` will, in certain configurations,
put `127.0.0.53` as an upstream into `/etc/resolv.conf`. Kubernetes (`kubelet`) by default
will pass this `/etc/resolv/conf` file to all Pods using the `default` dnsPolicy (this
includes CoreDNS Pods). CoreDNS then uses this `/etc/resolv.conf` as a list of upstreams
to proxy/forward requests to.  Since it contains a local address, CoreDNS ends up forwarding
requests to itself.

There are many ways to work around this issue, some are listed here:

* Add the following to `kubelet`: `--resolv-conf /run/systemd/resolve/resolv.conf`.  This flag
tells `kubelet` to pass an alternate `resolv.conf` to Pods. For `systemd-resolved`,
`/run/systemd/resolve/resolv.conf` is typically the location of the "original" `/etc/resolv.conf`.
* Disable `systemd-resolved` on host nodes, and restore `/etc/resolv.conf` to the original.
* A quick and dirty fix is to edit your Corefile, replacing `proxy . /etc/resolv.conf` with
the ip address of your upstream DNS, for example `proxy . 8.8.8.8`.
