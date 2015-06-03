DOCKER = docker
REPO = dimroc/urbanevents-tools

TAG = $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
ifeq ($(TAG), master)
	TAG = latest
else ifeq ($(TAG), HEAD)
	TAG = latest
endif

all: build

run: build
	$(DOCKER) run -i -t $(REPO)

build:
	$(DOCKER) build -t $(REPO):$(TAG) .

push:
	$(DOCKER) push $(REPO):$(TAG)
