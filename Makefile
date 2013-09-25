MAJOR    := 0
MINOR    := 1
PATCH    := 9
VERSION  := $(MAJOR).$(MINOR).$(PATCH)
LDFLAGS := -ldflags "-X main.Version $(VERSION)"
TARGET := tsinkf

default: test

build:
	go build $(LDFLAGS) -o $(TARGET)

test:
	go test

dist/$(TARGET)-darwin-v$(VERSION):
	mkdir -p dist
	go build $(LDFLAGS) -o dist/$(TARGET)-darwin-v$(VERSION)

dist/$(TARGET)-linux-v$(VERSION):
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(TARGET)-linux-v$(VERSION)

dist: dist/$(TARGET)-linux-v$(VERSION) dist/$(TARGET)-darwin-v$(VERSION)
	@echo DONE
