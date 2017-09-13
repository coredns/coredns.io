+++
date = "2016-12-19T08:00:24Z"
description = "A introduction into writing middleware for CoreDNS"
tags = ["middleware", "coredns", "plugins", "documentation"]
title = "Writing Middleware for CoreDNS"
author = "miek"
+++

As CoreDNS uses Caddy for setting up and using middleware, the process of writing middleware is
remarkably [similar](https://github.com/mholt/caddy/wiki/Writing-a-Plugin:-Directives). This post
slightly reworks (and simplifies in some cases) those pages.

A new middleware adds new functionality to CoreDNS, i.e. *caching*, *metrics* and basic *zone* file
serving are all middlewares.

If you want to write a new middleware and want it to be included by default, i.e. merged in the code
base please open an [issue](https://github.com/coredns/coredns/issues) first to discuss initial design
and other things that may come up.

## How to Register a CoreDNS Middleware

[Caddy](https://caddyserver.com) supports different type of plugins, CoreDNS, for instance,
registers itself as a Caddy server type plugin. CoreDNS currently supports one type of plugins;
a *middleware*.

Start a new Go package with an init function. Then register your middleware:

```go
import "github.com/mholt/caddy"

func init() {
  caddy.RegisterPlugin("foo", caddy.Plugin{
    ServerType: "dns",
    Action:     setup,
  })
}
```

Every middleware must have a name, `foo`, in this case. The *ServerType* must be `dns`. The *Action*
speficied here is to say it will call a function called `setup` whenever the *directive* `foo` is
encountered in the Corefile.

### The Setup Function

The *Action* field of the `caddy.Plugin` struct is what makes a directive plugin unique. This is the
function to run when CoreDNS is parsing and executing the Corefile.

The action is simply a function that takes a caddy.Controller and returns an error:
(We use [middleware.Error](https://godoc.org/github.com/coredns/coredns/middleware#Error) to prefix
returned error with `middleware/foo: ` to improve error reporting).

``` go
func setup(c *caddy.Controller) error {
  if err != nil {
    return middleware.Error("foo", err)
  }

  return nil
}
```

It is the responsibility of the setup function to parse the directive's tokens and configure itself.
The [Controller struct](https://godoc.org/github.com/mholt/caddy#Controller)
makes this quite easy. If we expect a line in the Corefile such as:

```
foo gizmo
```

We can get the value of the first argument ("foobar") like so:

```go
for c.Next() {              // Skip the directive name.
    if !c.NextArg() {       // Expect at least one value.
        return c.ArgErr()   // Otherwise it's an error.
    }
    value := c.Val()        // Use the value.
}
```
You parse the tokens present for your directive by iterating over `c.Next()` which is true as long
as there are more tokens to parse. Since a directive may appear multiple times, you must iterate
over `c.Next()` to get all the appearances of your directive and consume the first token (which is the
directive name).

### Adding to CoreDNS

To plug your middleware into CoreDNS, import it. This is done in
[core/coredns.go](https://github.com/coredns/coredns/blob/master/core/coredns.go):


```go
import _ "your/plugin/package/path/here"
```

This makes CoreDNS compile your plugin, but it is still not available, so the second step is
to add it to [directives.go](https://github.com/coredns/coredns/blob/master/core/dnsserver/directives.go):

Add the name (`foo`) of your middleware at the end of the file in the `directives` string slice.
Note the ordering is important, because this is determines how the plugins are chained together.

## How DNS Middleware Works in CoreDNS

Check out the [godoc for the middleware
package](http://godoc.org/github.com/coredns/coredns/middleware). The most important type is
[middleware.Handler](https://godoc.org/github.com/coredns/coredns/middleware#Handler).

A `Handler` is a function that handles a DNS request. CoreDNS will do all the bookkeeping of setting
up an DNS server for you, but you need to implement these two types.

### Writing a Handler

`middleware.Handler` is an interface similar to `http.Handler` except that it deals with DNS and the
`ServeDNS` method returns `(int, error)`. The `int` is the DNS rcode, and the `error` is one that
should be handled and/or logged. Read the
[middleware.md](https://github.com/coredns/coredns/blob/master/middleware.md) doc for more details
about these return values.

Handlers are usually a struct with at least one field, the next Handler in the chain:

```go
type MyHandler struct {
  Next middleware.Handler
}
```

To implement the `middleware.Handler` interface, we write a method called `ServeDNS`.
This method is the actual handler function, and, unless it fully handles the request by itself, it
should call the next handler in the chain:

```go
func (h MyHandler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
  return h.Next.ServeDNS(ctx, w, r)
}
```

The interface also needs a method `func Name() string`, this is mainly used to see if an middleware
is active, i.e. the *auto* middleware used this to detect if *metrics* are active and if so, adds
updates the zone info.

```go
func (h MyHandler) Name() string { return "foo" }
```

That's all there is to it.

## How to Add the Handler to CoreDNS

So, back in your setup function. You've just parsed the tokens and set up your middleware handler
with all the proper configuration:

```go
func setup(c *caddy.Controller) error {
  for c.Next() {
    // Get configuration.
  }

  // Now what?
}
```

To chain in your new handler, get the config for the current site from the dnsserver package.
Then wrap your handler in a middleware function:

```go
cfg := dnsserver.GetConfig(c)
mid := func(next middleware.Handler) middleware.Handler {
  return MyHandler{Next: next}
}
cfg.AddMiddleware(mid)
```

And you're done! Of course, in this example, we simply allocated a `MyHandler` with no special
configuration. It doesn't really matter as long as you chain in the `next` handler properly!

## Examples

Simple examples of middleware that can be found in CoreDNS are:

* [root](https://godoc.org/github.com/coredns/coredns/middleware/root); does not register itself as
  a middleware. It simply performs some setup.
* [chaos](https://godoc.org/github.com/coredns/coredns/middleware/chaos); a DNS middleware that
  responds to `CH txt version.bind` requests.

**Don't forget: the best documentation is the [godoc](https://godoc.org/github.com/coredns/coredns)
and the [code](https://github.com/coredns/coredns) itself!**
