# Makefile
SHELL := /bin/bash
.PHONY: all build-web move-web build-go package clean

all: clean build-web move-web build-go package

build-web:
	cd web && npm run build

move-web:
	mkdir -p out/webui
	cp -r web/dist/* out/webui/

build-go:
	go build -o out/main main.go

package:
	cp msyh.ttf out/

clean:
	rm -rf out
