+++
date = "2017-07-23T22:00:00Z"
description = "When should middleware be external?"
tags = ["Middleware", "External", "Out-of-Tree"]
title = "When Should Middleware be External?"
author = "miek"
+++

The `middleware.md` in the CoreDNS source tree has some pointers on what a middleware for CoreDNS
should have as minimum requirements. It basically boils down to: documentation, tests and working.

It is easier to list when a middleware can be *included* in CoreDNS than to say it should stay
external, so we will do the first:

* First, the middleware should be useful for *other* people. "Useful" is a subjective term, but the
  middleware needs to fill a niche that appeals to more than one person.
* It should be sufficiently different from other middleware to warrant inclusion.
* Current internet standards need be supported: IPv4 and IPv6, so A and AAAA records should be
  handled (if your middleware is in the business of dealing with address records that is).
* It must have tests.
* It must have a README.md for documentation.

Middleware for CoreDNS can live out-of-tree, `middleware.cfg` defaults to CoreDNS' repo but other
repos work just as well, so it is pretty easy to use your middleware even though it is out-of-tree.
