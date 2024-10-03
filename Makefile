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
	rm -fr pkg

clean_debian:	
	rm -fr debian/kafkatool
	rm -f debian/files
	rm -f debian/kafkatool.*
	rm -fr debian/.debhelper
	rm -fr debian/DEBIAN

tar: clean
	cd .. \
	&& tar \
	--exclude='.git' \
	--exclude='bin' \
	--exclude='build' \
	--exclude='pkg' \
	--exclude='src/kafkatool/kafkatool' \
	--exclude='.gitignore' \
	-cJvf $(BINARY_NAME)_$(VERSION).orig.tar.xz kafkatool