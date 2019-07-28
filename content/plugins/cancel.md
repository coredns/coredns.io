+++
title = "cancel"
description = "*cancel* a plugin that cancels a request's context after 5001 milliseconds."
weight = 6
tags = [ "plugin", "cancel" ]
categories = [ "plugin" ]
date = "2019-07-28T20:04:45.448852"
+++

## Description

The *cancel* plugin creates a canceling context for each request. It adds a timeout that gets
triggered after 5001 milliseconds.

The 5001 number is chosen because the default timeout for DNS clients is 5 seconds, after that they
give up.

A plugin interested in the cancellation status should call `plugin.Done()` on the context. If the
context was canceled due to a timeout the plugin should not write anything back to the client and
return a value indicating CoreDNS should not either; a zero return value should suffice for that.

~~~ txt
cancel [TIMEOUT]
~~~

* **TIMEOUT** allows setting a custom timeout. The default timeout is 5001 milliseconds (`5001 ms`)

## Examples

~~~ corefile
. {
    cancel
    whoami
}
~~~

Or with a custom timeout:

~~~ corefile
. {
    cancel 1s
    whoami
}
~~~

## Also See

The Go documentation for the context package.
