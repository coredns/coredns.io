all:
	hugo -d /opt/www/new.coredns.io

.PHONY: clean
clean:
	rm -rf /opt/www/new.coredns.io/*

.PHONY: test
test:
	hugo server
