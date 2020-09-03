SHELL := /bin/bash

all: run

.PHONY: run
run: 
	make -C tests/always-use-https-on run
	make -C tests/always-use-https-off run
	make -C tests/browser-cache-ttl run
	make -C tests/cache-level run
	make -C tests/minify run
	make -C tests/mobile-redirect run