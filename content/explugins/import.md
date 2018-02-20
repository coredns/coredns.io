+++
title = "import"
description = "*import* - allows including configuration from another file."
weight = 10
tags = [  "plugin" , "import" ]
categories = [ "plugin", "external" ]
date = "2017-07-22T12:37:19+01:00"
repo = "https://github.com/mholt/caddy"
home = "https://caddyserver.com/docs/import"
default = "yes"
+++

## Description

This is a unique directive in that *import* can appear outside of a server block. In other words, it
can appear at the top of a Corefile where an address would normally be. Like other directives,
however, it cannot be used inside of other directives.

Note that the the import path is relative to the Corefile, not the current working directory.

## Syntax

~~~
import PATTERN
~~~

* **PATTERN** is the file or glob pattern (`*`) to include. Its contents will replace this line, as
  if that file's contents appeared here to begin with. This value is relative to the file's
  location. It is an error if a specific file cannot be found, but an empty glob pattern is not an
  error.

## Examples

Import a shared configuration:

~~~
import config/common.conf
~~~

Imports any files found in the zones directory:

~~~
import ../zones/*
~~~
