# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
#GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
MKDIR_P=mkdir -p
GITCMD=git
BINARY_NAME=shurli
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_OSX=$(BINARY_NAME)_osx
BINARY_WIN=$(BINARY_NAME).exe
DIST_DIR=dist
DIST_OSX=$(DIST_DIR)_osx
DIST_OSX_PATH=$(DIST_DIR)/$(DIST_OSX)
DIST_UNIX=$(DIST_DIR)_unix
DIST_UNIX_PATH=$(DIST_DIR)/$(DIST_UNIX)
DIST_WIN=$(DIST_DIR)_win
DIST_WIN_PATH=$(DIST_DIR)/$(DIST_WIN)
DIST_FILES=chains.json config.json.sample favicon.png LICENSE README.md assets public sagoutil templates
CP_AV=cp -av

all: build
build: deps
	$(GITCMD) checkout grewal
	$(GOBUILD) -o $(BINARY_NAME) -v
# test: 
#	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -rf $(DIST_DIR)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_OSX)
run: deps
	$(GITCMD) checkout grewal
	$(GOBUILD) -o $(BINARY_NAME) -v 
	./$(BINARY_NAME) start
deps:
	$(GOGET) -u github.com/satindergrewal/kmdgo
	$(GOGET) -u github.com/Meshbits/shurli
	$(GOGET) -u github.com/gorilla/mux
	$(GOGET) -u github.com/gorilla/websocket

# Cross compilation
build-linux: deps
	$(GITCMD) checkout grewal
	$(MKDIR_P) $(DIST_UNIX_PATH)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DIST_UNIX_PATH)/$(BINARY_NAME) -v
	$(CP_AV) $(DIST_FILES) $(DIST_UNIX_PATH)
build-osx: deps
	$(GITCMD) checkout grewal
	$(MKDIR_P) $(DIST_OSX_PATH)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DIST_OSX_PATH)/$(BINARY_NAME) -v
	$(CP_AV) $(DIST_FILES) $(DIST_OSX_PATH)
build-win: deps
	$(GITCMD) checkout grewal
	$(MKDIR_P) $(DIST_WIN_PATH)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(DIST_WIN_PATH)/$(BINARY_WIN) -v
	$(CP_AV) $(DIST_FILES) $(DIST_WIN_PATH)
# docker-build:
# 	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/Meshbits/shurli golang:latest go build -o "$(BINARY_UNIX)" -v
