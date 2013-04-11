MAJOR    := 0
MINOR    := 0
PATCH    := 1
VERSION  := $(MAJOR).$(MINOR).$(PATCH)

default: example

build:
	go build
test:
	go test

example:
	go run sink.go -from="find . -type f" -to="wc -l"
