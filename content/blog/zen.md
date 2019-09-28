+++
title = "The Zen of CoreDNS"
description = "Philosophy behind CoreDNS"
tags = ["Zen", "Meta"]
draft = false
date = "2019-09-28T13:55:25Z"
author = "CoreDNS Authors"
+++

CoreDNS is a DNS server, written in Go. It links *plugins* to provide a customized DNS service.
This document details the philosophy behind it and provides guidance to future maintainers and
users alike.

# The Zen of CoreDNS

Our main principle is that **we want to maintain an excellent implementation of a DNS server**.

The separation of CoreDNS' functionality into plugins helps enormously with this goal. The quality
of these plugins has a direct impact on CoreDNS. The **better the plugins** are, the **better
CoreDNS** is.

The *non-plugin*, *core* part of CoreDNS is there to support plugins and to abstract as much of the
DNS protocol as possible. This allows a plugin author to focus on the code quality and let CoreDNS
take care of any DNS peculiarities.

With its plugins, CoreDNS brings the [Unix
philosophy](https://en.wikipedia.org/wiki/Unix_philosophy) to a DNS server.

To make good plugins:

 *  Write plugins that do *one* thing and do it *well*.

 *  Write plugins to *work together*.

This means plugins should be focussed and small. This makes them easy to maintain and debug. Good
plugins lead to a DNS server that is highly capable, and straightforward in its architecture and
implementation.

## What Does This Mean For Users?

 *  CoreDNS is easy to use.

 *  Each plugin is easy to use. The default configuration works for the majority of use-cases.

 *  No unnecessary configuration knobs are exposed.

 *  A plugin is well documented.

## What Does This Mean For Plugin Authors?

Ease of use puts a lot of emphasis on the plugin's author to make this happen.

 *  A plugin needs only a minimum number of settings to make it work. If values for settings can be
    determined automatically, they should be *determined* **automatically**.

 *  Adhere to the DNS specification.

 *  Have decent test coverage.

When contemplating a new feature, ask yourself "can this be done in a new, separate plugin?"
