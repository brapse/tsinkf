MAJOR    := 0
MINOR    := 1
PATCH    := 4
VERSION  := $(MAJOR).$(MINOR).$(PATCH)
LDFLAGS := -ldflags "-X main.Version $(VERSION)"
TARGET := tsinkf

default: test

build:
	go build $(LDFLAGS) -o $(TARGET)

test:
	go test
