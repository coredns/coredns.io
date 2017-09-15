+++
title = "startup"
description = "*startup* executes a command when the server begins."
weight = 10
tags = [  "plugin" , "startup" ]
categories = [ "plugin", "external" ]
date = "2017-07-22T12:37:19+01:00"
repo = "https://github.com/mholt/caddy"
home = "https://caddyserver.com/docs/startup"
enabled = "default"
+++

This is useful for preparing to serve a zone by running a script or starting a background process.
(Also see [shutdown](/explugins/shutdown).)

Each command that is executed at startup is blocking, unless you suffix the command with a space and
`&`, which will cause the command to be run in the background. The output and error of the command go
to stdout and stderr, respectively. There is no stdin.

A command will only be executed once for each time it appears in the Corefile.

## Syntax

~~~
startup COMMAND
~~~

* **COMMAND** is the command to execute; it may be followed by arguments.

## Examples

Start command before the server starts listening:

~~~
startup /etc/init.d/command start
~~~

On windows, you might need to use quotes when the command path contains spaces:

~~~
startup "\"C:\Program Files\command.exe\" -b 127.0.0.1:9123" &
~~~

## How to Enable

This external plugin is included in the default CoreDNS configuration.
