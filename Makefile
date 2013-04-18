MAJOR    := 0
MINOR    := 0
PATCH    := 1
VERSION  := $(MAJOR).$(MINOR).$(PATCH)

LDFLAGS := -ldflags "-X main.Version $(VERSION)"

SRC=$(wildcard *.go)
TGT=tsinkf

default: example

build:
	go build

test: build
	go test

example: build
	./tsinkf run -from="find /bin -type f|head" -to="wc -l"
	./tsinkf show -v

