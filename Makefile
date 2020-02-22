NAME		:= takanawa
VERSION		:= 0.0.0
OUTPUT		:= bin

GOOS		:= $(shell go env GOOS)
GOARCH		:= $(shell go env GOARCH)

LDFLAGS		:= -w -X main.version=$(VERSION)
EXTLDFLAGS	:=
TAGS		:=

ifeq ($(GOOS),windows)
LDFLAGS		:= $(LDFLAGS) -H=windowsgui
EXTLDFLAGS	:= $(EXTLDFLAGS) -static
endif

ifneq (,$(filter $(GOOS),linux freebsd netbsd openbsd dragonfly))
TAGS		:= $(TAGS) netgo
EXTLDFLAGS	:= $(EXTLDFLAGS) -static
endif

ifeq ($(GOOS),darwin)
TAGS		:= $(TAGS) netgo
LDFLAGS		:= $(LDFLAGS) -s
EXTLDFLAGS	:= $(EXTLDFLAGS) -sectcreate __TEXT __info_plist Info.plist
endif

ifeq ($(GOOS),android)
LDFLAGS		:= $(LDFLAGS) -s
endif

GO.build	:= GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -tags '$(TAGS)' -ldflags '$(LDFLAGS) -extldflags "$(EXTLDFLAGS)"'
GO.test		:= go test -v -cover -coverprofile=coverage.out -covermode=atomic

.PHONY: default
default: build

.PHONY: build
build: $(OUTPUT)
	$(GO.build) -o $(OUTPUT) ./...

.PHONY: test
test:
	$(GO.test) ./...

.PHONY: clean
clean:
	$(RM) $(wildcard $(OUTPUT)/*)

$(OUTPUT):
	mkdir -p $@
