.PHONY: build docker

TARGET_DIR ?= ./

# TEST_ARGS = -v | grep -c RUN
DOCKER_OPTS ?= --rm
DOCKER_REPO ?= github.com/influenzanet/counter-service
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse HEAD)
TAG ?= $(DOCKER_REPO):$(VERSION)

build:
	go build -o $(TARGET_DIR) ./cmd/counter-service

run:
	go run ./cmd/counter-service

docker:
	docker build -t $(TAG) -f build/docker/Dockerfile $(DOCKER_OPTS) --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(REVISION) .
