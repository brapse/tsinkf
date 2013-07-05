MAJOR    := 0
MINOR    := 1
PATCH    := 4
VERSION  := $(MAJOR).$(MINOR).$(PATCH)

default: test

build:
	go build

test:
	go test
