BINARY = hex_demo
VET_REPORT = vet.report.txt
TEST_REPORT = test.report.txt
GOARCH = amd64

VERSION=$(shell date '+%Y%m%d.%H%M%S')
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
DOCKER_REPO=docker_repo_here/${BINARY}
DOCKER_LATEST_TAG=${DOCKER_REPO}:latest
DOCKER_VERSION_TAG=${DOCKER_REPO}:${VERSION}-${BRANCH}
HOST=$(shell hostname)

CURRENT_DIR=$(shell pwd)
BUILD_DIR=$(CURRENT_DIR)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH} -X main.BUILDIMAGE=${DOCKER_VERSION_TAG} -X main.BUILDHOST=${HOST}"

# Build the project
all: clean test fmt vet linux docker

runlocal:
	GO111MODULE=on GOARCH=${GOARCH} go run ${LDFLAGS} . -l

runserver:
	CONFIG_FILE="`pwd`/config/localhost.yaml" GO111MODULE=on GOARCH=${GOARCH} go run ${LDFLAGS} .

linux: clean
	GO111MODULE=on GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o tmp/${BINARY}-linux ./cmd/proxy/

test:
	GO111MODULE=on go test -v ../../pkg/... 2>&1

fmt:
	GO111MODULE=on go fmt ../../pkg/...

watch:
	watch -n 1 -d 'mysql hex_demo -uroot -e "select * from cars"'

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -rf tmp

.PHONY: linux test vet fmt clean docker docker-push
