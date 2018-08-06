+++
title = "on"
description = "*on* - executes a command when a specified event is triggered."
weight = 10
tags = [  "plugin" , "on" ]
categories = [ "plugin", "external" ]
date = "2018-01-22T07:53:19+01:00"
repo = "https://github.com/mholt/caddy"
home = "https://caddyserver.com/docs/on"
default = "yes"
+++

## Description

*on* executes a command when a specified event is triggered. This can be useful for preparing to
 serve a site by running a script or starting a background process when the server starts, or for
 stopping it when the server exits.

Each command that is executed is blocking, unless you suffix the command with a space and &, which
will cause the command to be run in the background. (Do not do this when the server is exiting, or
the command may not finish before its parent process exits.) The output and error of the command go
to standard output and standard error, respectively. There is no standard input.

A command will only be executed once for each time it appears in the Corefile. In other words, even
if this directive is shared by more than one zone, a command will only execute once per appearance
in the Corefile.

Note that commands scheduled for the shutdown event will not execute if CoreDNS is force-terminated,
for example, by using a "Force Quit" feature provided by your operating system. However, a typical
SIGINT (Ctrl+C) will allow the shutdown commands to execute.

## Syntax

~~~
on EVENT COMMAND
~~~

**EVENT** is the name of the *event* on which to execute the **COMMAND** (see list below).
**COMMAND** is the command to execute; it may be followed by arguments.

### Events

Commands can execute on the following events:

* `startup` - The server instance is starting or starting up
* `shutdown` - The server instance is shutting down (not restarting)

## Examples

Start php-fpm before the server starts listening:

~~~
on startup /etc/init.d/php-fpm start
~~~

Stop php-fpm when the server quits:

~~~
on shutdown /etc/init.d/php-fpm stop
~~~

On Windows, you might need to use quotes when the command path contains spaces:

~~~
on startup "\"C:\Program Files\PHP\v7.0\php-cgi.exe\" -b 127.0.0.1:9123" &
~~~
