+++
title = "CoreDNS-1.1.1 Release"
description = "CoreDNS-1.1.1 Release Notes."
tags = ["Release", "1.1.1", "Notes"]
release = "1.1.1"
date = "2018-03-24T09:33:29+00:00"
author = "coredns"
+++

We are pleased to announce the [release](https://github.com/coredns/coredns/releases/tag/v1.1.1) of
CoreDNS-1.1.1!

This release fixes a **critical bug** in the *cache* plugin found by [Cure53](/2018/03/15/cure53-security-assessment/).

All users are encouraged to upgrade.

## Core

Fix a corner case when scrubbing the reply to fit the request's buffer.

## Plugins

* [*cache*](/plugins/cache) fixes the critical spoof vulnerability.
* [*route53*](/plugins/route53) adds support for PTR records.

## Contributors

The following people helped with getting this release done:

Chris O'Haver,
Mario Kleinsasser,
Miek Gieben,
Yong Tang.

And of course the people in [Cure53](https://cure53.de). Also special shout out to Mario Kleinsasser
for helping to debug.

For documentation see our (in progress!) [manual](/manual). For help and other resources, see our
[community page](https://coredns.io/community/).
