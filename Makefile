PLUGINDIR:=.coredns/plugin

all:
	hugo -d public/
	cp _redirects public/

# this is the old website
old:
	hugo -d /var/www/coredns.io


.PHONY: deps
deps:
	go get github.com/miekg/corecheck

.PHONY: run
run:
	hugo server

# Sync CoreDNS' plugin README.md's to coredns.io. Also sync the release notes from the notes directory.
.PHONY: sync
sync:
	@GO111MODULE=on go run bin/sync.go -plugindir $(PLUGINDIR)
	@cp -vu $(PLUGINDIR)/../notes/coredns-* content/blog

# Scan all markdown files for Corefile snippets and check validity github.com/miekg/corecheck
.PHONY: test
test:
	hugo server
