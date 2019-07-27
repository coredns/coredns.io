REPO:=github.com/coredns/coredns.io.git

all:
	hugo -d public/

.PHONY: clean
clean:
	rm -rf coredns
	rm -rf public/*

.PHONY: deps
deps:
	go install github.com/miekg/corecheck
	git clone https://github.com/coredns/coredns.git
	( cd coredns; git checkout -b $$(git describe --abbrev=0 --tags); make; git stash )

.PHONY: test
test:
	hugo server

# Sync CoreDNS' plugin README.md's to coredns.io. Also sync the release notes from the notes directory.
.PHONY: sync
sync:
	( rm -f content/plugins/* )
	( cd bin; python2.7 sync-from-coredns.py ../coredns/plugin )
	cp coredns/notes/coredns-* content/blog

# Scan all markdown files for Corefile snippets and check validity github.com/miekg/corecheck
.PHONY: scan
scan:
	@echo "Corefile checking within blog posts might fail due to old Corefiles"
	corecheck -exe coredns/coredns -dir content/blog
	@echo "Corefile checking within manuals might fail due to not running within Kubernetes context and external plugins (unbound)"
	corecheck -exe coredns/coredns -dir content/manual

.PHONY: netlify-push
netlify-push:
ifeq ($(CONTEXT), production)
	@echo "Automatic content sync from coredns/coredns"
	git config user.name "auto_sync"
	git config user.email auto_sync@coredns.io
	git add content
	git commit -m "Automatic content sync from coredns/coredns at tag: $(shell cd coredns; git describe --abbrev=0 --tags)"
	git remote add origin https://$(REPO)
	git config credential.helper store
	@echo "https://$(GIT_USER):$(GIT_PW)@$(REPO)" >> ~/.git-credentials
	git push origin master:master
endif

.PHONY: netlify
netlify: clean deps sync netlify-push
	hugo -d public/