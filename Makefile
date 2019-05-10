.PHONY: all dev build-dev run-dev deps build install

DOCKER_IMAGE := freighter-dev
DOCKER_FILE := Dockerfile

GOPATH ?= /go
PROJECTSUBDIR ?= /src/github.com/allankerr/freighter

PROJECTDIR := $(GOPATH)$(PROJECTSUBDIR)

all: deps build install

dev: build-dev run-dev

build-dev:
	docker build -f "$(DOCKER_FILE)" -t "$(DOCKER_IMAGE)" --build-arg PROJECT_DIR="$(PROJECTDIR)" .

run-dev:
	docker run --privileged -v "$(CURDIR):$(PROJECTDIR)" -it "$(DOCKER_IMAGE)"

deps:
	go get

build:
	go build

install:
	go install