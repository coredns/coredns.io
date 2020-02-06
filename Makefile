all:
	hugo -d public/

.PHONY: deps
deps:
	go get github.com/miekg/corecheck

.PHONY: run
run:
	hugo server

# Sync CoreDNS' plugin README.md's to coredns.io. Also sync the release notes from the notes directory.
.PHONY: sync
sync:
	@GO111MODULE=on go run bin/sync.go
	@cp -vu .coredns/notes/coredns-* content/blog

# Scan all markdown files for Corefile snippets and check validity github.com/miekg/corecheck
.PHONY: test
test:
	corecheck -exe $(PWD)/.coredns/coredns -dir content/blog
	corecheck -exe $(PWD)/.coredns/coredns -dir content/manual
