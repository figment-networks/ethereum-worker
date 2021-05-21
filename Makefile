LDFLAGS      := -w -s
MODULE       := github.com/figment-networks/ethereum-worker
VERSION_FILE ?= ./VERSION


# Git Status
GIT_SHA ?= $(shell git rev-parse --short HEAD)

ifneq (,$(wildcard $(VERSION_FILE)))
VERSION ?= $(shell head -n 1 $(VERSION_FILE))
else
VERSION ?= n/a
endif



.PHONY: build-live
build-live: LDFLAGS += -X $(MODULE)/cmd/ethereum-worker/config.Timestamp=$(shell date +%s)
build-live: LDFLAGS += -X $(MODULE)/cmd/ethereum-worker/config.Version=$(VERSION)
build-live: LDFLAGS += -X $(MODULE)/cmd/ethereum-worker/config.GitSHA=$(GIT_SHA)
build-live:
	go build -o ethereum-worker-live -ldflags '$(LDFLAGS)'  ./cmd/ethereum-worker-live

.PHONY: pack-release
pack-release:
	@mkdir -p ./release
	@make build-live
	@mv ./ethereum-worker-live ./release/worker-live
	@zip -r ethereum-worker ./release
	@rm -rf ./release
