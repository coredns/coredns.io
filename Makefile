all:
	hugo -d /opt/www/coredns.io

.PHONY: clean
clean:
	rm -rf /opt/www/coredns.io/*

.PHONY: test
test:
	hugo server

# Scan all markdown files for Corefile snippets and check validity
# github.com/miekg/corecheck
.PHONY: scan
scan:
	corecheck -dir content/blog

.PHONY: sync-from-coredns
sync-from-coredns:
	( cd bin; ./sync-from-coredns.py $$GOPATH/src/github.com/coredns/coredns/plugin )
