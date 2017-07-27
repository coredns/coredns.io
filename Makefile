all:
	hugo -d /opt/www/coredns.io

.PHONY: clean
clean:
	rm -rf /opt/www/coredns.io/*

.PHONY: test
test:
	hugo server

.PHONY: sync-from-coredns
sync-from-coredns:
	( cd bin; ./readme-to-hugo.py $$GOPATH/src/github.com/coredns/coredns/middleware )
