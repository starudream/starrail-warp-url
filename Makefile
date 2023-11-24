PROJECT ?= $(shell basename $(CURDIR))
MODULE  ?= $(shell go list -m)
VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)

BITTAGS :=
LDFLAGS := -s -w
LDFLAGS += -X "github.com/starudream/go-lib/core/v2/config/version.gitVersion=$(VERSION)"

.PHONY: init
init:
	git status -b -s
	go mod tidy

.PHONY: bin
bin: init
	CGO_ENABLED=0 go build -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' -o bin/$(PROJECT) $(MODULE)/cmd

.PHONY: bin-windows
bin-windows: init
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags '$(BITTAGS)' -ldflags '$(LDFLAGS)' -o bin/$(PROJECT).exe $(MODULE)/cmd

.PHONY: run
run: bin
	DEBUG=true bin/$(PROJECT) $(ARGS)

.PHONY: watch-web
watch-web:
	tailwindcss -i ./web/_index.css -o ./web/index.css -w

.PHONY: build-web
build-web:
	tailwindcss -i ./web/_index.css -o ./web/index.css -m
