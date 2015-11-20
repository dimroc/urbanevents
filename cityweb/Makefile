BIN = $(GOPATH)/bin
NODE_BIN = $(shell npm bin)
PID = .pid
GO_FILES = $(filter-out app/server/bindata.go, $(shell find app -type f -name "*.go"))
BINDATA = app/server/bindata.go
BINDATA_FLAGS = -pkg=server -prefix=app/server/data
BUNDLE = app/server/data/static/build/bundle.js
APP = $(shell find app/client -type f)
GO_APP = github.com/dimroc/urbanevents/cityweb/app

build: clean $(BIN)/app

clean:
	@rm -rf app/server/data/static/build/*
	@rm -rf app/server/data/bundle.server.js
	@rm -rf $(BINDATA)
	@echo cleaned

$(BUNDLE): $(APP)
	@$(NODE_BIN)/webpack --progress --colors

$(BIN)/app: $(BUNDLE) $(BINDATA)
	#echo $(GO_APP)
	#@go install -ldflags "-w -X main.buildstamp `date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.gittag `git describe --tags || true` -X main.githash `git rev-parse HEAD || true`" $(GO_APP)
	@go install $(GO_APP)
	#echo done

kill:
	@kill `cat $(PID)` || true

serve: clean $(BUNDLE)
	@make restart
	@$(NODE_BIN)/webpack-dev-server --config webpack.hot.config.js $$! > $(PID)_wds &
	@ANYBAR_WEBPACK=yep $(NODE_BIN)/webpack --progress --colors --watch $$! > $(PID)_wp &
	@fswatch $(GO_FILES) | xargs -n1 -I{} make restart || make kill
	@kill `cat $(PID)_wp` || true
	@kill `cat $(PID)_wds` || true

restart: BINDATA_FLAGS += -debug
restart: $(BINDATA)
	@make kill
	@go install $(GO_APP)
	@$(BIN)/app run & echo $$! > $(PID)

$(BINDATA):
	$(BIN)/go-bindata $(BINDATA_FLAGS) -o=$@ app/server/data/...

lint:
	@eslint app/client || true
	@golint $(filter-out app/main.go, $(GO_FILES)) || true
	@golint -min_confidence=1 app
