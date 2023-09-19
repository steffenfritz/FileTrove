ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=ftrove
VERSION=1.0.0-DEV-1
BUILD=`git rev-parse --short HEAD`
PLATFORMS=darwin linux windows
ARCHITECTURES=amd64

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -w -s" 
BLDFLAGS=-buildmode=pie

all: clean build_all

build:
	go build ${BLDFLAGS} ${LDFLAGS} -o ${BINARY}

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build $(BLDFLAGS) $(LDFLAGS) -v -o $(BINARY)-$(GOOS))))
	tar cvfz ${BINARY}_macos_$(VERSION).tar.gz ${BINARY}-darwin
	tar cvfz ${BINARY}_linux_$(VERSION).tar.gz ${BINARY}-linux
	mv ${BINARY}-windows ${BINARY}.exe && tar cvfz ${BINARY}_win_$(VERSION).tar.gz ${BINARY}.exe


clean:
	rm -f '${BINARY}-linux'
	rm -f '${BINARY}-darwin'
	rm -f '${BINARY}.exe'

.PHONY: clean build build_all all
