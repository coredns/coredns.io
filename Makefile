all:
	hugo -d /opt/www/coredns.io

.PHONY:
clean:
	rm -rf /opt/www/coredns.io/*

.PHONY:
test:
	hugo server
