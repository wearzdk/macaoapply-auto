# Makefile
SHELL := /bin/bash
.PHONY: all build-web move-web package clean build-go-linux build-go-darwin build-go-windows

all: clean build-web move-web package build-go-linux build-go-darwin build-go-windows

build-web:
	cd web && VITE_APP_BASE_URL=/webui/ && npm run build

move-web:
	mkdir -p out/webui
	cp -r web/dist/* out/webui/

build-go-linux:
	GOOS=linux GOARCH=amd64 go build -o out/apply-linux main.go

build-go-darwin:
	GOOS=darwin GOARCH=amd64 go build -o out/apply-darwin main.go

build-go-windows:
	GOOS=windows GOARCH=amd64 go build -o out/apply-win.exe main.go

package:
	cp msyh.ttf out/

clean:
	rm -rf out
