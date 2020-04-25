+++
title = "alias"
description = "*alias* - replaces zone apex CNAMEs."
weight = 10
tags = [  "plugin" , "alias" ]
categories = [ "plugin", "external" ]
date = "2020-04-25T19:12:00+08:00"
repo = "https://github.com/serverwentdown/alias"
home = "https://github.com/serverwentdown/alias/blob/master/README.md"
+++

## Description

The *alias* plugin eliminates CNAME records from zone apex by making the subsequent resolved records look like they belong to the zone apex. This behaves similarily to [CloudFlare's Zone Flattening](https://support.cloudflare.com/hc/en-us/articles/200169056-CNAME-Flattening-RFC-compliant-support-for-CNAME-at-the-root).

This plugin works only with plugins that produce A or AAAA records alongside the CNAME record. Examples include `auto` and `file`. However, you might need to adjust the order of this plugin to use it with other plugins. 

> Preferrably, this should not be used in favour of the RFC drafts for the new [ANAME](https://tools.ietf.org/html/draft-ietf-dnsop-aname-00) records, but the DNS library used by CoreDNS does not support ANAME records yet. 

## Syntax

```
alias
```

## Examples

```
example.com {
	file db.example.com
	alias
}
# This is used to resolve CNAME records by the `file` plugin. Modify accordingly
. {
	forward . 1.1.1.1 1.0.0.1
}
```

This will transform responses like this:

```
;; ANSWER SECTION:
example.com.		3600	IN	CNAME	two.example.org.
two.example.org.	3600	IN	CNAME	one.example.net.
one.example.net.	3600	IN	A	127.0.0.1
```

into

```
;; ANSWER SECTION:
example.com.		3600	IN	A	127.0.0.1
```

See [`example/`](https://github.com/serverwentdown/alias/tree/master/example/) for a more extensive example. 
