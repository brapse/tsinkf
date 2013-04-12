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
	go run tsinkf.go -from="cat example_input" -to="wc -l"
