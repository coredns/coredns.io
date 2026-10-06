+++
title = "git"
description = "*git* - pull git repositories."
weight = 10
tags = [  "plugin" , "git" ]
categories = [ "plugin", "external" ]
date = "2021-01-07T00:12:00+08:00"
repo = "https://github.com/ganawaj/coredns-git"
home = "https://github.com/ganawaj/coredns-git/blob/master/README.md"
+++

## Description

*git* clones a git repository into the site. This makes it possible to deploy your zones with a
simple git push.

The *git* plugin starts a service routine that runs during the lifetime of the server. When the
service starts, it clones the repository. While the server is still up, it pulls the latest every
so often. You can also set up a webhook to pull immediately after a push. In regular git fashion, a
pull only includes changes, so it is very efficient.

If a pull fails, the service will retry up to three times. If the pull was not successful by then,
it won't try again until the next interval.

This plugin *requires* `git` to be installed on the system.

Webhooks are not supported, this is a pure pull model.

## Syntax

~~~ txt
git REPO [PATH]
~~~

 *  **REPO** is the URL to the repository; SSH and HTTPS URLs are supported

 *  **PATH** is the path, relative to site root, to clone the repository into; default is site root

This simplified syntax pulls from master every 3600 seconds (1 hour) and only works for public
repositories.

For more control or to use a private repository, use the following syntax:

~~~
git [REPO PATH] {
	repo        REPO
	path        PATH
	branch      BRANCH
	interval    INTERVAL
	args        ARGS
	pull_args   PULL_ARGS
}
~~~

 *  **REPO** is the URL to the repository; only HTTPS URLs are supported.

 *  **PATH** is the path to clone the repository into; default is site root (if set). It can be
    absolute or relative (to site root). See the *root* plugin.

 *  **BRANCh** is the branch or tag to pull; default is master branch. **`{latest}`** is a
    placeholder for latest tag which ensures the most recent tag is always pulled.

 *  **INTERVAl** is the number of seconds between pulls; default is 3600 (1 hour), minimum 5. An
    interval of -1 disables periodic pull.

 *  **ARGS** is the additional cli args to pass to `git clone` e.g. `--depth=1`. `git clone` is
    called when the source is being fetched the first time.

 *  **PULL_ARGS** is the additional cli args to pass to `git pull` e.g. `-s recursive -X theirs`.
    `git pull` is used when the source is being updated.

## Examples

Public repository pulled into site root every hour:

~~~
git github.com/user/myproject
~~~

Public repository pulled into the "subdir" directory in the site root:

~~~
git github.com/user/myproject subdir
~~~

## Also See

The *root* plugin for setting the root.
