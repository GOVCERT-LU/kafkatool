GOPATH=$(shell pwd)
GOINSTALL := GOPATH=$(GOPATH) go install
GOGET := GOPATH=$(GOPATH) go get
GODEP := GOPATH=$(GOPATH) $(GOPATH)/bin/dep
GOCLEAN := GOPATH=$(GOPATH) go clean

BINARY_NAME=kafkatool

prefix = /usr/local

VERSION = $(shell dpkg-parsechangelog -S Version | cut -d'-' -f 1)

all: build bash_completion

install: all
	mkdir -p $(DESTDIR)$(prefix)/bin
	cp bin/$(BINARY_NAME) $(DESTDIR)$(prefix)/bin

	mkdir -p $(DESTDIR)/etc/bash_completion.d
	cp kafkatool_completion.sh $(DESTDIR)/etc/bash_completion.d

bash_completion: build
	bin/kafkatool completion

dependencies:
	# make sure golang dep is installed
	$(GOGET) -u github.com/golang/dep/cmd/dep

	# make sure golang depencies as defined in src/kafkatool/Gopkg.toml are installed and up-to-date
	cd src/$(BINARY_NAME) && $(GODEP) ensure

build: dependencies
	$(GOINSTALL) $(BINARY_NAME)

clean: 
	$(GOCLEAN)
	rm -fr bin/
	rm -fr pkg/
	rm -f kafkatool_completion.sh
	rm -fr src/github.com
	rm -fr src/kafkatool/vendor
	rm -f src/kafkatool/Gopkg.lock

clean_debian:	
	rm -fr debian/kafkatool
	rm -f debian/files
	rm -f debian/kafkatool.*
	rm -fr debian/.debhelper

tar: clean
	cd .. \
	&& tar \
	--exclude='.git' \
	-cjvf $(BINARY_NAME)_$(VERSION).orig.tar.bz2 kafkatool	