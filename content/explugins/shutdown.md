+++
title = "shutdown"
description = "*shutdown* - executes a command when the server exits."
weight = 10
tags = [  "plugin" , "shutdown" ]
categories = [ "plugin", "external" ]
date = "2017-07-22T12:37:19+01:00"
repo = "https://github.com/mholt/caddy"
home = "https://caddyserver.com/docs/shutdown"
default = "yes"
+++

## Description

This is useful for performing cleanup or stopping a background process. (Also see
[startup](/explugins/startup).)

Each command that is executed at shutdown is blocking. The output and error of the command go to
stdout and stderr, respectively. There is no stdin.

Note that shutdown commands will not execute if CoreDNS is force-terminated, for example, by using
a "Force Quit" feature provided by your operating system. However, a typical SIGINT (Ctrl+C) will
allow the shutdown commands to execute.

Even if this directive is shared by more than one host, the command will only execute once per
appearance in the Corefile.

## Syntax

~~~
shutdown COMMAND
~~~

* **COMMAND** is the command to execute; it may be followed by arguments.

## Examples

Stop command when the server quits:

~~~
shutdown /etc/init.d/command stop
~~~
