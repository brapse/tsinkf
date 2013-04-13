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
test:
	go test

example: build
	./sink -from="cat example_input" -to="wc -l" -v
