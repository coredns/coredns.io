+++
date = "2017-09-16T07:28:00Z"
description = "Semantic versioning for CoreDNS."
tags = ["Versioning", "Semver", "Backwards", "Compatiblity"]
title = "Semantic Versioning"
author = "miek"
+++

CoreDNS' next release is around the corner and it is going to be **1.0.0**. With this release to
move to [semantic versioning](http://semver.org/). This will allow us to make changes, some of
which may be backwards incompatible, in a sane manor:

> Given a version number MAJOR.MINOR.PATCH, increment the:
>
> MAJOR version when you make incompatible API changes,
> MINOR version when you add functionality in a backwards-compatible manner, and
> PATCH version when you make backwards-compatible bug fixes.

With respect to Go, we will support the **last two** released versions. At the time of the writing
this means we develop in Go1.9, and support Go1.8.
