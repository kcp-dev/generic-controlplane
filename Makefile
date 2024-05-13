# Copyright 2024 The KCP Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# We need bash for some conditional logic below.
SHELL := /usr/bin/env bash -e

#-----------------------------------------------------------------------------
# Workaround git issues on OpenShift Prow CI, where the user running in the
# job is not guaranteed to own the repo checkout.
#-----------------------------------------------------------------------------
ifeq ($(CI),true)
   $(shell git config --global --add safe.directory '*')
endif

GO_INSTALL = ./hack/go-install.sh

TOOLS_DIR=hack/tools
TOOLS_GOBIN_DIR := $(abspath $(TOOLS_DIR))
GOBIN_DIR=$(abspath ./bin)
PATH := $(GOBIN_DIR):$(TOOLS_GOBIN_DIR):$(PATH)
TMPDIR := $(shell mktemp -d)
KIND_CLUSTER_NAME ?= kind

# Detect the path used for the install target
ifeq (,$(shell go env GOBIN))
INSTALL_GOBIN=$(shell go env GOPATH)/bin
else
INSTALL_GOBIN=$(shell go env GOBIN)
endif

OPENSHIFT_GOIMPORTS_VER := c72f1dc2e3aacfa00aece3391d938c9bc734e791
OPENSHIFT_GOIMPORTS_BIN := openshift-goimports
OPENSHIFT_GOIMPORTS := $(TOOLS_DIR)/$(OPENSHIFT_GOIMPORTS_BIN)-$(OPENSHIFT_GOIMPORTS_VER)
export OPENSHIFT_GOIMPORTS # so hack scripts can use it

GOLANGCI_LINT_VER := v1.54.2
GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT := $(TOOLS_GOBIN_DIR)/$(GOLANGCI_LINT_BIN)-$(GOLANGCI_LINT_VER)

STATICCHECK_VER := 2023.1
STATICCHECK_BIN := staticcheck
STATICCHECK := $(TOOLS_GOBIN_DIR)/$(STATICCHECK_BIN)-$(STATICCHECK_VER)

GOTESTSUM_VER := v1.8.1
GOTESTSUM_BIN := gotestsum
GOTESTSUM := $(abspath $(TOOLS_DIR))/$(GOTESTSUM_BIN)-$(GOTESTSUM_VER)


ARCH := $(shell go env GOARCH)
OS := $(shell go env GOOS)

KUBE_MAJOR_VERSION := $(shell go mod edit -json | jq '.Require[] | select(.Path == "k8s.io/kubernetes") | .Version' --raw-output | sed 's/v\([0-9]*\).*/\1/')
KUBE_MINOR_VERSION := $(shell go mod edit -json | jq '.Require[] | select(.Path == "k8s.io/kubernetes") | .Version' --raw-output | sed "s/v[0-9]*\.\([0-9]*\).*/\1/")
GIT_COMMIT := $(shell git rev-parse --short HEAD || echo 'local')
GIT_DIRTY := $(shell git diff --quiet && echo 'clean' || echo 'dirty')
GIT_VERSION := $(shell go mod edit -json | jq '.Require[] | select(.Path == "k8s.io/kubernetes") | .Version' --raw-output)+kcp-$(shell git describe --tags --match='v*' --abbrev=14 "$(GIT_COMMIT)^{commit}" 2>/dev/null || echo v0.0.0-$(GIT_COMMIT))
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := \
	-X k8s.io/client-go/pkg/version.gitCommit=${GIT_COMMIT} \
	-X k8s.io/client-go/pkg/version.gitTreeState=${GIT_DIRTY} \
	-X k8s.io/client-go/pkg/version.gitVersion=${GIT_VERSION} \
	-X k8s.io/client-go/pkg/version.gitMajor=${KUBE_MAJOR_VERSION} \
	-X k8s.io/client-go/pkg/version.gitMinor=${KUBE_MINOR_VERSION} \
	-X k8s.io/client-go/pkg/version.buildDate=${BUILD_DATE} \
	\
	-X k8s.io/component-base/version.gitCommit=${GIT_COMMIT} \
	-X k8s.io/component-base/version.gitTreeState=${GIT_DIRTY} \
	-X k8s.io/component-base/version.gitVersion=${GIT_VERSION} \
	-X k8s.io/component-base/version.gitMajor=${KUBE_MAJOR_VERSION} \
	-X k8s.io/component-base/version.gitMinor=${KUBE_MINOR_VERSION} \
	-X k8s.io/component-base/version.buildDate=${BUILD_DATE} \
	-extldflags '-static'
all: build
.PHONY: all

ldflags:
	@echo $(LDFLAGS)

.PHONY: require-%
require-%:
	@if ! command -v $* 1> /dev/null 2>&1; then echo "$* not found in ${PATH}"; exit 1; fi

build: WHAT ?= ./cmd/...
build: require-jq require-go require-git verify-go-versions ## Build the project
	set -x; for W in $(WHAT); do \
  		W=$$(echo "$${W}" | sed 's,^\./,github.com/kcp-dev/generic-controlplane/,') && \
    	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build $(BUILDFLAGS) -ldflags="$(LDFLAGS)" -o bin $${W}; \
    done
.PHONY: build

.PHONY: build-all
build-all:
	GOOS=$(OS) GOARCH=$(ARCH) $(MAKE) build WHAT='./cmd/...'

install: WHAT ?= ./cmd/...
install: require-jq require-go require-git verify-go-versions ## Install the project
	set -x; for W in $(WHAT); do \
  		W=$$(echo "$${W}" | sed 's,^\./,github.com/kcp-dev/generic-controlplane/,') && \
  		GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go install -ldflags="$(LDFLAGS)" $${W}; \
  	done
.PHONY: install

$(GOLANGCI_LINT):
	GOBIN=$(TOOLS_GOBIN_DIR) $(GO_INSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint $(GOLANGCI_LINT_BIN) $(GOLANGCI_LINT_VER)

$(STATICCHECK):
	GOBIN=$(TOOLS_GOBIN_DIR) $(GO_INSTALL) honnef.co/go/tools/cmd/staticcheck $(STATICCHECK_BIN) $(STATICCHECK_VER)

$(LOGCHECK):
	GOBIN=$(TOOLS_GOBIN_DIR) $(GO_INSTALL) sigs.k8s.io/logtools/logcheck $(LOGCHECK_BIN) $(LOGCHECK_VER)

lint: $(GOLANGCI_LINT) $(STATICCHECK) $(LOGCHECK) ## Verify lint
	$(GOLANGCI_LINT) run --timeout 20m ./...
	$(STATICCHECK) -checks ST1019,ST1005 ./...
.PHONY: lint

$(GOTESTSUM):
	GOBIN=$(TOOLS_GOBIN_DIR) $(GO_INSTALL) gotest.tools/gotestsum $(GOTESTSUM_BIN) $(GOTESTSUM_VER)

tools: $(GOLANGCI_LINT) $(CONTROLLER_GEN) $(KCP_APIGEN_GEN) $(YAML_PATCH) $(GOTESTSUM) $(OPENSHIFT_GOIMPORTS) $(CODE_GENERATOR) ## Install tools
.PHONY: tools

$(OPENSHIFT_GOIMPORTS):
	GOBIN=$(TOOLS_GOBIN_DIR) $(GO_INSTALL) github.com/openshift-eng/openshift-goimports $(OPENSHIFT_GOIMPORTS_BIN) $(OPENSHIFT_GOIMPORTS_VER)

.PHONY: imports
imports: $(OPENSHIFT_GOIMPORTS) verify-go-versions
	$(OPENSHIFT_GOIMPORTS) -m github.com/kcp-dev/generic-controlplane

$(TOOLS_DIR)/verify_boilerplate.py:
	mkdir -p $(TOOLS_DIR)
	curl --fail --retry 3 -L -o $(TOOLS_DIR)/verify_boilerplate.py https://raw.githubusercontent.com/kubernetes/repo-infra/master/hack/verify_boilerplate.py
	chmod +x $(TOOLS_DIR)/verify_boilerplate.py

.PHONY: verify-boilerplate
verify-boilerplate: $(TOOLS_DIR)/verify_boilerplate.py ## Verify boilerplate
	$(TOOLS_DIR)/verify_boilerplate.py --boilerplate-dir=hack/boilerplate --skip docs/venv

GO_TEST = go test
ifdef USE_GOTESTSUM
GO_TEST = $(GOTESTSUM) $(GOTESTSUM_ARGS) --
endif

test: WHAT ?= ./...
# We will need to move into the sub package, of sdk to run those tests.
test: ## Run tests
	$(GO_TEST) -race $(COUNT_ARG) -coverprofile=coverage.txt -covermode=atomic $(TEST_ARGS) $$(go list "$(WHAT)")

.PHONY: verify-imports
verify-imports: ## Verify imports
	hack/verify-imports.sh

.PHONY: verify-go-versions
verify-go-versions: ## Verify go versions
	hack/verify-go-versions.sh

.PHONY: clean
clean: clean-workdir ## Clean all
	rm -fr $(TOOLS_DIR)
	rm -f $(GOBIN_DIR)/*

.PHONY: clean-workdir
clean-workdir: WORK_DIR ?= .
clean-workdir: ## Clean workdir
	rm -fr $(WORK_DIR)/.gcp*


.PHONY: help
help: ## Show this help
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

include Makefile.venv
