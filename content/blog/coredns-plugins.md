+++
date = "2016-12-19T08:00:24Z"
description = "A introduction into writing plugin for CoreDNS."
tags = ["plugin", "coredns", "plugins", "documentation"]
title = "Writing Plugins for CoreDNS"
author = "miek"
+++

A plugin adds functionality to CoreDNS, i.e. *caching*, *metrics* and basic *zone* file serving are
all plugins.

If you want to write a new plugin and want it to be included by default, i.e. merged in the code
base please open an [issue](https://github.com/coredns/coredns/issues) first to discuss initial
design and other things that may come up. Starting with a README file to explain how things work
from a user perspective is usually a good idea.

See the [example plugin](https://github.com/coredns/example) for, uh, an example for how to
structure, write and test a plugin. There are plenty of comments in the code to help you along.

## How to Register a CoreDNS Plugin?

When writing your plugin code you will need to register it with CoreDNS. This can be done by calling
the following function:

```go
func init() { plugin.Register("foo", setup) }
```

Every plugin must have a name, `foo`, in this case. When `foo` is encountered in the configuration
the `setup` function will be called in this package.

### The Setup Function

The `setup` function (it may be called different, but pretty much every plugin just calls it
`setup`) parses the configuration and populates internal data structures.

The setup function a `caddy.Controller` and returns an error: (We use
[plugin.Error](https://godoc.org/github.com/coredns/coredns/plugin#Error) to prefix returned error
with `plugin/foo:` to improve error reporting).

``` go
func setup(c *caddy.Controller) error {
  if err != nil {
    return plugin.Error("foo", err)
  }

  // various other code

  return nil
}
```

If we see a line in the Corefile such as:

```
foo gizmo
```

We can get the value of the first argument ("gizmo") like so:

```go
for c.Next() {              // Skip the plugin name, "foo" in this case.
    if !c.NextArg() {       // Expect at least one value.
        return c.ArgErr()   // Otherwise it's an error.
    }
    value := c.Val()        // Use the value.
}
```
You parse the tokens present for your plugin by iterating over `c.Next()` which is true as long
as there are more tokens to parse. Since a plugin may appear multiple times, you must iterate over
`c.Next()` to get all the appearances of your plugin and consume the tokens.

### Adding to CoreDNS

To plug your plugin into CoreDNS, put in
[`plugin.cfg`](https://github.com/coredns/coredns/blob/master/plugin.cfg) and run `go generate`.

## How A Plugin Works in CoreDNS

Check out the [godoc for the plugin
package](https://godoc.org/github.com/coredns/coredns/plugin). The most important type is
[plugin.Handler](https://godoc.org/github.com/coredns/coredns/plugin#Handler).

A `Handler` is a function that handles a DNS request. CoreDNS will do all the bookkeeping of setting
up an DNS server for you, but you need to implement these two types.

### Writing a Handler

`plugin.Handler` is an interface similar to `http.Handler` except that it deals with DNS and the
`ServeDNS` method returns `(int, error)`. The `int` is status code, and the `error` is logged (if
not nil) See [plugin.md](https://github.com/coredns/coredns/blob/master/plugin.md) for more details
about these return values.

Handlers are usually a struct with at least one field, the next Handler in the chain:

```go
type MyHandler struct {
  Next plugin.Handler
}
```

To implement the `plugin.Handler` interface, we write a method called `ServeDNS`. This method is the
actual handler function, and, unless it fully handles the request by itself, it should call the next
handler in the chain:

```go
func (h MyHandler) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
  return h.Next.ServeDNS(ctx, w, r)
}
```

The interface also needs a method `func Name() string`.,

```go
func (h MyHandler) Name() string { return "foo" }
```

That's all there is to it (apart from writing all code that actually does something with the DNS
request of course).

## Further Reading

Simple examples of plugin that can be found in CoreDNS are:

* [root](https://godoc.org/github.com/coredns/coredns/plugin/root); does not register itself as
  a plugin. It simply performs some setup.
* [chaos](https://godoc.org/github.com/coredns/coredns/plugin/chaos); a DNS plugin that
  responds to `CH txt version.bind` requests.
* [example](https://github.com/coredns/example); an example plugin that prints "example" when
  responding to a query.

**Don't forget: the best documentation is the [godoc](https://godoc.org/github.com/coredns/coredns)
and the [code](https://github.com/coredns/coredns) itself!**
