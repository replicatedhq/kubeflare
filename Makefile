
SHELL := /bin/bash
VERSION ?=`git describe --tags`
DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"`
VERSION_PACKAGE = github.com/replicatedhq/kubeflare/pkg/version
GIT_TREE = $(shell git rev-parse --is-inside-work-tree 2>/dev/null)
IMAGE_TAG ?= dev
ifneq "$(GIT_TREE)" ""
define GIT_UPDATE_INDEX_CMD
git update-index --assume-unchanged
endef
define GIT_SHA
`git rev-parse HEAD`
endef
else
define GIT_UPDATE_INDEX_CMD
echo "Not a git repo, skipping git update-index"
endef
define GIT_SHA
""
endef
endif

define LDFLAGS
-ldflags "\
	-X ${VERSION_PACKAGE}.version=${VERSION} \
	-X ${VERSION_PACKAGE}.gitSHA=${GIT_SHA} \
	-X ${VERSION_PACKAGE}.buildTime=${DATE} \
"
endef

export GO111MODULE=on
# export GOPROXY=https://proxy.golang.org

all: generate fmt vet manifests kubeflare

.PHONY: clean-and-tidy
clean-and-tidy:
	@go clean -modcache ||:
	@go mod tidy ||:

.PHONY: integration
integration: integration-bin
	make -C integration run

.PHONY: integration-bin
integration-bin: generate fmt vet manifests
	go build \
		${LDFLAGS} \
		-i \
		-o bin/integration \
		./cmd/integration

.PHONY: test
test: generate fmt vet manifests
	go test ./pkg/... ./cmd/... -coverprofile cover.out

.PHONY: kubeflare
kubeflare: clean-and-tidy generate fmt vet bin/kubeflare 

.PHONY: bin/kubeflare
bin/kubeflare:
	go build \
		${LDFLAGS} \
		-i \
		-o bin/kubeflare \
		./cmd/kubeflare

.PHONY: run
run: generate fmt vet bin/kubeflare
	./bin/kubeflare run \
	--log-level debug 

.PHONY: install
install: manifests generate dev
	kubectl apply -f config/crds/v1

.PHONY: deploy
deploy: manifests
	kubectl apply -f config/crds/v1
	kustomize build config/default | kubectl apply -f -

.PHONY: manifests
manifests: controller-gen
	$(CONTROLLER_GEN) \
		rbac:roleName=manager-role webhook \
		crd:crdVersions=v1beta1 \
		output:crd:artifacts:config=config/crds/v1beta1 \
		paths="./..."
	$(CONTROLLER_GEN) \
		rbac:roleName=manager-role webhook \
		crd:crdVersions=v1 \
		output:crd:artifacts:config=config/crds/v1 \
		paths="./..."

.PHONY: fmt
fmt:
	go fmt ./pkg/... ./cmd/...

.PHONY: vet
vet:
	go vet ./pkg/... ./cmd/...

.PHONY: generate
generate: controller-gen client-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./pkg/apis/...
	$(CLIENT_GEN) \
		--output-package=github.com/replicatedhq/kubeflare/pkg/client \
		--clientset-name kubeflareclientset \
		--input-base github.com/replicatedhq/kubeflare/pkg/apis \
		--input crds/v1alpha1 \
		-h ./hack/boilerplate.go.txt

.PHONY: dev
dev: kubeflare
	docker build -t kubeflare/kubeflare-manager -f ./Dockerfile.manager .
	docker tag kubeflare/kubeflare-manager localhost:32000/kubeflare/kubeflare-manager:latest
	docker push localhost:32000/kubeflare/kubeflare-manager:latest

.PHONY: image
image: kubeflare
	docker build -t kubeflare/kubeflare-manager:$(IMAGE_TAG) -f ./Dockerfile.manager .

.PHONY: contoller-gen
controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.1
CONTROLLER_GEN=$(shell go env GOPATH)/bin/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

.PHONY: client-gen
client-gen:
ifeq (, $(shell which client-gen))
	go get k8s.io/code-generator/cmd/client-gen@kubernetes-1.18.0
CLIENT_GEN=$(shell go env GOPATH)/bin/client-gen
else
CLIENT_GEN=$(shell which client-gen)
endif
