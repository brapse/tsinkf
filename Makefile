MAJOR    := 0
MINOR    := 0
PATCH    := 1
VERSION  := $(MAJOR).$(MINOR).$(PATCH)

default:
	go build

test:
	go test
