.PHONY: build docker

TARGET_DIR ?= ./

# TEST_ARGS = -v | grep -c RUN
DOCKER_OPTS ?= --rm
DOCKER_REPO ?=github.com/influenzanet/counter-service
VERSION := $(shell git describe --tags --abbrev=0)

TAG ?= $(VERSION)

build:
	go build -o $(TARGET_DIR) ./cmd/counter-service

docker:
	docker build -t $(DOCKER_REPO):$(VERSION)  -f build/docker/Dockerfile $(DOCKER_OPTS) .