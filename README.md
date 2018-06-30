# coredns.io

Purpose:

1. What is CoreDNS?
1. How to use it?
1. How to get involved and get help.

The website is created with HUGO, you'll need to download and install that if you want to locally
work on the website. For creating new content it is not needed.

We have three types of pages.

* Home - `themes/coredns/layouts/index.html`
* Blog/News/Documentation - `content/blog/*` (depends on params, see `archetypes`)
* Plugins - `content/plugins/*` and `content/explugins/*`

Any blog post you'll tag with 'stale' will pop up a warning after 9 months that the content of
the post may not reflect the current workings of CoreDNS.

Releases should be tagged with 'release = "number" and `data/coredns.toml` should be updated
with the number.

`data/subtext.toml` controls the buttons for "Docs", "Plugins" and "External Plugins".

Create:

* new release: `hugo new -k release blog/coredns-<number>.md`
* new blog: `hugo new -k blog blog/coredns-<number>.md`
* new blog that will also be documentation: `hugo new -k doc blog/coredns-<number>.md`

## Style

For the colors we have:

* Dark: #280071   (40, 0, 113)
* Middle: #5F259F (95, 37, 159)
* Light: #8246AF  (130, 70, 175)

The open source font Lato is very close to Brandon Grotesque (which is used in the logo). Try to use
this for any imagery you add.

See the `style` directory for Libre Office Draw templates you can use.

## Popup

See the top-level (i.e. not in the themes directory) `layout/partials/popup.html` for adjusting the
text. This can used for important notifications.

## Corefile Snippets

Any Corefile snippets should be use the (fake) language `corefile`, we have a small utility that
checks all these snippets to see if they are still valid.
