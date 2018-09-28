GLIDE = $(GOPATH)/bin/glide

SRCROOT ?= $(realpath .)
BUILD_ROOT ?= $(SRCROOT)

# These are paths used in the docker image
SRCROOT_D = /go/src/common_go
BUILD_ROOT_D = $(SRCROOT_D)/tmp/dist

VERSION = $$(git name-rev --tags --name-only $$(git rev-parse HEAD))
COMMIT = $$(git rev-parse --short HEAD)

default: test

# test runs the unit tests and vets the code
test: dist
	docker-compose down
	docker-compose up --exit-code-from tests || exit

cover: dep
	@go tool cover 2>/dev/null; if [ $$? -eq 3 ]; then \
		go get -u golang.org/x/tools/cmd/cover; \
	fi
	@go test $(TEST) -cover -test.v=true -test.coverprofile=c.out
	sed -i -e "s#.*/\(.*\.go\)#\./\\1#" c.out
	@go tool cover -html=c.out -o c.html

trace: dep
	@go test $(TEST) -trace trace.out common_go
	@go tool trace common_go.test trace.out

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@echo "go tool vet *.go"
	@go tool vet ./*.go ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
	fi

build: dep vet

	@echo $(GOPATH)

	CGO_ENABLED=0 GO15VENDOREXPERIMENT=1 go build -x \
	-o $(BUILD_ROOT)/common_go \
	-ldflags "-s -X main.version=$(VERSION) -X main.commit=$(COMMIT)" \
	.

# Establish dependencies
dep:
	go get -u github.com/Masterminds/glide

	if [ -f glide.lock ]; then rm -f glide.lock; fi
	if [ -d vendor ]; then rm -Rf vendor; fi
	GO15VENDOREXPERIMENT=1 $(GLIDE) install --force

dist:
	docker pull golang:1.9.1

    # Mount using -v source folder
    # Pass in env variables using -e for src root
	docker run --rm \
	           -v $(SRCROOT):$(SRCROOT_D) \
	           -w $(SRCROOT_D) \
	           -e BUILD_ROOT=$(BUILD_ROOT_D) \
               -e SRCROOT=$(SRCROOT_D) \
               -e UID=`id -u` \
               -e GID=`id -g` \
	           golang:1.9.1 \
	           make distbuild

distbuild: clean build
	chown -R $(UID):$(GID) $(SRCROOT)

clean:

.PHONY: bin default dep test vet dist clean
