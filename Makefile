GO ?= go
GOCACHE ?= $(PWD)/.gocache
GOMODCACHE ?= $(PWD)/.gocache/mod
CACHE_ENV = GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE)

.PHONY: all fmt test build plugin integ integ-incident clean

all: test

fmt:
	$(GO)fmt ./...

test:
	$(CACHE_ENV) $(GO) test ./...

build:
	$(CACHE_ENV) $(GO) build ./...

plugin:
	$(CACHE_ENV) $(GO) build -o bin/incidentplugin ./cmd/incidentplugin

integ-incident:
	$(CACHE_ENV) $(GO) run ./integ/incident.go

integ: integ-incident

clean:
	rm -rf $(GOCACHE) bin
