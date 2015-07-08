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
		-e "CITYSERVICE_SETTINGS=config/nyc.json" \
		-e "GO_ENV=qa" -e "PORT=8080" \
		-e "TWITTER_CONSUMER_SECRET=$(TWITTER_CONSUMER_SECRET)" -e "TWITTER_CONSUMER_KEY=$(TWITTER_CONSUMER_KEY)" \
		-e "TWITTER_TOKEN=$(TWITTER_TOKEN)" -e "TWITTER_TOKEN_SECRET=$(TWITTER_TOKEN_SECRET)" \
		-e "INSTAGRAM_CLIENT_ID=$(INSTAGRAM_CLIENT_ID)" -e "INSTAGRAM_CLIENT_SECRET=$(INSTAGRAM_CLIENT_SECRET)" \
		-e "TUTUM_SERVICE_FQDN=http://4975e96c.ngrok.io" \
		$(REPO)

build:
	$(DOCKER) build -t $(REPO):$(TAG) .

push:
	$(DOCKER) push $(REPO):$(TAG)
