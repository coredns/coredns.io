# Writing Plugins

As mentioned before in this manual, plugins are the thing that make CoreDNS tick. We've seen
a bunch of configuration in the [previous section](#setups). But how can you write your own plugin?

See [Writing Plugins for CoreDNS](/2016/12/19/writing-plugins-for-coredns/) for an older post on
this subject. The [plugin.md](https://github.com/coredns/coredns/blob/master/plugin.md) documented
in CoreDNS' source also has some background and talks about styling the README.md.

The canonical example plugin, is the [*example*](/plugins/example) plugin. It's [github
repository](https://github.com/coredns/example) shows the most minimal code (with tests!) that is
needed to create plugin.

It has:

1. `setup.go` and `setup_test.go` that implement that parse the configuration from the Corefile.
   The (usually named) `setup` function is called whenever the Corefile parser see the plugin's
   name. In this case "example".
2. `example.go` (usually named `<plugin_name>.go`) handled the logic of handling the query, and
   `example_test.go` has basic units tests to check if the plugin works.
3. The `README.md` that documents in a Unix manual style how this plugin can be configured.
4. A LICENSE file; for inclusion in CoreDNS this needs to an APL like license.

TODO: extend
