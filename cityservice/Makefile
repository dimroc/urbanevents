DOCKER = docker
REPO = tutum.co/dimroc/cityservice

TAG = $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
ifeq ($(TAG), master)
	TAG = latest
else ifeq ($(TAG), HEAD)
	TAG = latest
endif

all: build

run: build
	$(DOCKER) run -i -t \
		-e "ELASTICSEARCH_URL=$(ELASTICSEARCH_URL)" \
		-e "AWS_REGION=$(AWS_REGION)" -e "AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY)" -e "AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID)" \
		$(REPO)

build:
	$(DOCKER) build -t $(REPO):$(TAG) .

push:
	$(DOCKER) push $(REPO):$(TAG)
