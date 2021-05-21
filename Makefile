BINARY_NAME=kafkatool

prefix = /usr/local

VERSION = $(shell dpkg-parsechangelog -S Version | cut -d'-' -f 1)

all: compile bash_completion

install: all
	mkdir -p $(DESTDIR)$(prefix)/bin
	cp build/$(BINARY_NAME) $(DESTDIR)$(prefix)/bin

	mkdir -p $(DESTDIR)/etc/bash_completion.d
	cp kafkatool_completion.sh $(DESTDIR)/etc/bash_completion.d

bash_completion: compile
	build/$(BINARY_NAME) completion

compile:
	bash build.sh

clean: 
	rm -f build/$(BINARY_NAME)
	rm -f kafkatool_completion.sh

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