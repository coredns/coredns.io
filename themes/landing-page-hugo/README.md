# Landing Page Theme for Hugo

[Hugo](http://gohugo.io) theme based on the
[landing-page theme for Jekyll](https://github.com/swcool/landing-page-theme)

# How to use

## Background Images

Override the default intro and contact backgrounds by putting images in these files:

* `img/intro-bg.jpg`
* `img/contact-bg.jpg`


## About Section

Create a markdown file named `content/services/about/about.md`

```txt
---
Title: About Us
Draft: false
---

WidgetCo is the world leader in widget production.
```

## Services Section

Create a markdown file describing a service you offer in `content/services/` - e.g. `content/services/widgets.md`.

```txt
---
Title: Customized Widgets
Img: widgets.png
Category: Services
Draft: false
---

We specialize in bespoke widgets, built to your specification.
```

Then place a matching image in `static/img/services/` - e.g. `static/img/services/widgets.png`


## Social Contact Buttons

Contact buttons will be automatically created if one or more
`[[params.social]]` is configured in `config.toml`:

```toml
baseurl = "http://yourdomain.com"
languageCode = "en-us"
title = "WidgetCo Inc"
[params]
	description = "The best widgets in the world!"

[[params.social]]
	title = "email"
	icon = "envelope-o"
	url = "mailto:bushbama@whitehouse.gov"
[[params.social]]
	title = "twitter"
	icon = "twitter"
	url = "https://twitter.com/SBootstrap"
[[params.social]]
	title = "github"
	icon = "github"
	url = "https://github.com/IronSummitMedia/startbootstrap"
[[params.social]]
	title = "linkedin"
	icon = "linkedin"
	url = "http://linkedin.com/yourusername"
```

* `title` parameter sets the text to be displayed on the contact button
* `icon` parameter sets which [Font Awesome icon](http://fortawesome.github.io/Font-Awesome/icons/)
will be displayed.


## Google Analytics

Google Analytics support is automatically enabled if you set the
`googleAnalytics` param in `config.toml` to your Google Analytics tracking ID.

```toml
[params]
googleAnalytics = "UA-12345678-1"

```


# Demo
View this equivalent jekyll theme in action [here](https://swcool.github.io/landing-page-theme)

#Have a look, Have a try
```
git clone https://github.com/crakjie/landing-page-hugo.git
cd landing-page-hugo/exampleSite
hugo server -t landing-page-hugo
```
# Screenshot
![screenshot](https://raw.githubusercontent.com/swcool/landing-page-theme/master/img/screenshot.png)
