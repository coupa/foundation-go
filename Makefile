SRCROOT ?= $(realpath .)
BUILD_ROOT ?= $(SRCROOT)

# These are paths used in the docker image
SRCROOT_D = /go/src/github.com/coupa/foundation-go
BUILD_ROOT_D = $(SRCROOT_D)/tmp/dist

GOLANG_VER = 1.11.1

default: test

test: dist
	docker-compose down
	docker-compose up --exit-code-from tests || exit

build: vendor
	CGO_ENABLED=0 GO15VENDOREXPERIMENT=1 go build -x \
	-o $(BUILD_ROOT)/foundation \
	.

vendor: clean $(GOPATH)/bin/dep
	dep ensure

$(GOPATH)/bin/dep:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Used primarily to test, build is not much of use for a library in itself (which doesnt have a main package)
dist:
	docker pull golang:$(GOLANG_VER)

    # Mount using -v source folder
    # Pass in env variables using -e for src root
	docker run --rm \
        -v $(SRCROOT):$(SRCROOT_D) \
        -w $(SRCROOT_D) \
        -e BUILD_ROOT=$(BUILD_ROOT_D) \
        -e SRCROOT=$(SRCROOT_D) \
        -e UID=`id -u` \
        -e GID=`id -g` \
        golang:$(GOLANG_VER) \
        make distbuild

distbuild: clean build
	chown -R $(UID):$(GID) $(SRCROOT)

clean:
	if [ -f Gopkg.lock ]; then rm -f Gopkg.lock; fi
	if [ -d vendor ]; then rm -Rf vendor; fi
