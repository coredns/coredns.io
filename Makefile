C:=~/src/github.com/coredns/coredns

all:
	hugo -d /var/www/coredns.io

.PHONY: clean
clean:
	rm -rf /var/www/coredns.io/*

.PHONY: test
test:
	hugo server

# Sync CoreDNS' plugin README.md's to coredns.io. Also sync the release notes from the notes directory.
.PHONY: sync
sync:
	( rm -f content/plugins/* )
	( cd bin; ./sync-from-coredns.py $(C)/plugin )
	cp $(C)/notes/coredns-* content/blog

# Scan all markdown files for Corefile snippets and check validity github.com/miekg/corecheck
.PHONY: scan
scan:
	corecheck -exe $(C)/coredns -dir content/blog
	corecheck -exe $(C)/coredns -dir content/manual
