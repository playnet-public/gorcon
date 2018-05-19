###################### PlayNet Makefile ######################
#
# This Makefile is used to manage the command-line template
# All possible tools have to reside under their respective folders in cmd/
# and are being autodetected.
# 'make full' would then process them all while 'make toolname' would only
# handle the specified one(s).
# Edit this file with care, as it is also being used by our CI/CD Pipeline
# For usage information check README.md
#
# Parts of this makefile are based upon github.com/kolide/kit
#

NAME		:= gorcon
REPO		:= playnet-public
GIT_HOST	:= github.com
REGISTRY	:= quay.io
IMAGE		:= playnet/gorcon

PATH 		:= $(GOPATH)/bin:$(PATH)

VERSION		:= $(shell git describe --tags --always --dirty)
BRANCH 		:= $(shell git rev-parse --abbrev-ref HEAD)
REVISION 	:= $(shell git rev-parse HEAD)
REVSHORT 	:= $(shell git rev-parse --short HEAD)
USER 		:= $(shell whoami)

STAGING 	?= true
V			?= 0
NAMESPACES	?= debug

DOCKER_TAGS := -t $(REGISTRY)/$(IMAGE):$(VERSION) -t $(REGISTRY)/$(IMAGE):latest

-include .env

include helpers/make_version

.PHONY: build

### MAIN STEPS ###

all: test install run

# install required tools and dependencies
deps:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint
	go get -u github.com/haya14busa/goverage
	go get -u github.com/kisielk/errcheck
	go get -u github.com/maxbrunsfeld/counterfeiter
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega
	go get -u github.com/schrej/godacov
	go get -u golang.org/x/tools/cmd/goimports

# test entire repo
gotest:
	@go test -cover -race $(shell go list ./... | grep -v /vendor/)

test:
	@go get github.com/onsi/ginkgo/ginkgo
	@ginkgo -r -race -p

# install passed in tool project
install:
	GOBIN=$(GOPATH)/bin go install cmd/$(NAME)/*.go

# run tool
run:
	$(NAME)

# format entire repo (excluding vendor)
format:
	@go get golang.org/x/tools/cmd/goimports
	@find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	@find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

# build binary for docker image
buildgo: .pre-build
	CGO_ENABLED=0 GOOS=linux go build -i -o build/app -ldflags ${KIT_VERSION} -a -installsuffix cgo ./cmd/$(NAME)

# build docker image
build:
	@docker build --build-arg GIT_HOST=$(GIT_HOST) --build-arg REPO=$(REPO) --build-arg NAME=$(NAME) --build-arg COMMAND='buildgo' --no-cache --rm=true -t $(REGISTRY)/$(IMAGE)-build:$(VERSION) -f ./Dockerfile.build .
	@docker run -t $(REGISTRY)/$(IMAGE)-build:$(VERSION) /bin/true
	@docker cp `docker ps -q -n=1 -f ancestor=$(REGISTRY)/$(IMAGE)-build:$(VERSION) -f status=exited`:/go/src/$(GIT_HOST)/$(REPO)/$(NAME)/build .
	@docker rm `docker ps -q -n=1 -f ancestor=$(REGISTRY)/$(IMAGE)-build:$(VERSION) -f status=exited` || true
	docker build --no-cache --rm=true $(DOCKER_TAGS) --build-arg TOOL=$(NAME) -f Dockerfile.static .

# run specified tool from code
dev: generate
	@go run -ldflags $(KIT_VERSION) cmd/$(NAME)/*.go \
	-debug

# build the docker image
docker: build

# upload the docker image
upload:
	docker push $(REGISTRY)/$(IMAGE)

# clean build results and delete all images
clean:
	rm -rf build
	docker rmi -f $(shell docker images -q --filter=reference=$(REGISTRY)/$(IMAGE)*)

version:
	@echo $(VERSION)

# create build dir
.pre-build:
	@mkdir -p build

# helper to build new image and kick existing pod
update-deployment: docker upload clean restart-deployment

# delete existing pod to force imagePull (if latest)
restart-deployment:
	kubectl delete po -n $(NAMESPACE) -lapp=$(NAME)

# go quality checks
check: format lint vet

# vet entire repo (excluding vendor)
vet:
	@go vet $(shell go list ./... | grep -v /vendor/)

# lint entire repo (excluding vendor)
lint:
	@go get github.com/golang/lint/golint
	@golint -min_confidence 1 $(shell go list ./... | grep -v /vendor/)

# errcheck entire repo (excluding vendor)
errcheck:
	@go get github.com/kisielk/errcheck
	@errcheck -ignore '(Close|Write)' $(shell go list ./... | grep -v /vendor/)

cover:
	@go get github.com/haya14busa/goverage
	@go get github.com/schrej/godacov
	goverage -v -coverprofile=coverage.out $(shell go list ./... | grep -v /vendor/)

generate:
	@go get github.com/maxbrunsfeld/counterfeiter
	@go generate ./...
