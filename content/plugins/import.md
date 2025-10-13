+++
title = "import"
description = "*import* includes files or references snippets from a Corefile."
weight = 26
tags = ["plugin", "import"]
categories = ["plugin"]
date = "2025-10-13T05:58:44.87744810"
+++

## Description

The *import* plugin can be used to include files into the main configuration. Another use is to
reference predefined snippets. Both can help to avoid some duplication.

This is a unique plugin in that *import* can appear outside of a server block. In other words, it
can appear at the top of a Corefile where an address would normally be.

## Syntax

~~~
import PATTERN
~~~

*   **PATTERN** is the file, glob pattern (`*`) or snippet to include. Its contents will replace
    this line, as if that file's contents appeared here to begin with.

Corefile may contain at most 10000 import statements. A glob pattern counts as a single import. The limit protects the configuration from recursive imports.

## Files

You can use *import* to include a file or files. This file's location is relative to the
Corefile's location. It is an error if a specific file cannot be found, but an empty glob pattern is
not an error.

## Snippets

You can define snippets to be reused later in your Corefile by defining a block with a single-token
label surrounded by parentheses:

~~~ corefile
(mysnippet) {
	...
}
~~~

Then you can invoke the snippet with *import*:

~~~
import mysnippet
~~~

## Examples

Import a shared configuration:

~~~
. {
   import config/common.conf
}
~~~

Where `config/common.conf` contains:

~~~
prometheus
errors
log
~~~

This imports files found in the zones directory:

~~~
import ../zones/*
~~~

## See Also

See corefile(5).
