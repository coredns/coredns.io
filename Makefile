all:
	hugo -d /opt/www/coredns.io

.PHONY: clean
clean:
	rm -rf /opt/www/coredns.io/*

.PHONY: test
test:
	hugo server
