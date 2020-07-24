# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
#GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
MKDIR_P=mkdir -p
GITCMD=git
ROOT_DIR=$(shell pwd)
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
CURL_DL=curl -LJ
KMD_UNIX_URL=https://github.com/Meshbits/komodo/releases/download/release_buildv0.56_0.6.0/komodo_linux_v0.56.zip -o komodo_linux.zip
KMD_OSX_URL=https://github.com/Meshbits/komodo/releases/download/release_buildv0.56_0.6.0/komodo_macOS_v0.56.zip -o komodo_macos.zip
KMD_SUBATOMIC_WIN_URL=https://github.com/Meshbits/komodo/releases/download/debug_release2/subatomic_komodo_win_debug_bin_12Jul2020.zip -o subatomic_kmd_win.zip
SUBATOMIC_UNIX_URL=https://github.com/Meshbits/komodo_dapps/releases/download/release_buildv0.6/subatomic_linux_v0.6.zip -o subatomic_linux.zip
SUBATOMIC_OSX_URL=https://github.com/Meshbits/komodo_dapps/releases/download/release_buildv0.6/subatomic_macOS_v0.6.zip -o subatomic_macos.zip
RM_RFV=rm -rfv
UNZIP=unzip
TAR_GZ=tar -cvzf

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
	rm -rf $(DIST_UNIX_PATH)
	$(GITCMD) checkout grewal
	$(MKDIR_P) $(DIST_UNIX_PATH)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DIST_UNIX_PATH)/$(BINARY_NAME) -v
	$(CP_AV) $(DIST_FILES) $(DIST_UNIX_PATH)
	$(CURL_DL) $(KMD_UNIX_URL)
	$(UNZIP) -o komodo_linux.zip -d $(DIST_UNIX_PATH)/assets
	$(CURL_DL) $(SUBATOMIC_UNIX_URL)
	$(UNZIP) -o subatomic_linux.zip -d $(DIST_UNIX_PATH)/assets
	$(RM_RFV) komodo_linux.zip subatomic_linux.zip
	cd $(DIST_UNIX_PATH); tar -czvf ../shurli_linux.tar.gz *
	$(RM_RFV) $(DIST_UNIX_PATH)
	cd $(ROOT_DIR)
build-osx: deps
	$(GITCMD) checkout grewal
	$(MKDIR_P) $(DIST_OSX_PATH)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DIST_OSX_PATH)/$(BINARY_NAME) -v
	$(CP_AV) $(DIST_FILES) $(DIST_OSX_PATH)
	$(CURL_DL) $(KMD_OSX_URL)
	$(UNZIP) -o komodo_macos.zip -d $(DIST_OSX_PATH)/assets
	$(CURL_DL) $(SUBATOMIC_OSX_URL)
	$(UNZIP) -o subatomic_macos.zip -d $(DIST_OSX_PATH)/assets
	$(RM_RFV) komodo_macos.zip subatomic_macos.zip
	cd $(DIST_OSX_PATH); tar -czvf ../shurli_macos.tar.gz *
	$(RM_RFV) $(DIST_OSX_PATH)
	cd $(ROOT_DIR)
build-win: deps
	$(GITCMD) checkout grewal
	$(MKDIR_P) $(DIST_WIN_PATH)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(DIST_WIN_PATH)/$(BINARY_WIN) -v
	$(CP_AV) $(DIST_FILES) $(DIST_WIN_PATH)
	$(CURL_DL) $(KMD_SUBATOMIC_WIN_URL)
	$(UNZIP) -o subatomic_kmd_win.zip -d $(DIST_WIN_PATH)/assets
	$(RM_RFV) subatomic_kmd_win.zip
	cd $(DIST_WIN_PATH); tar -czvf ../shurli_win.tar.gz *
	$(RM_RFV) $(DIST_WIN_PATH)
	cd $(ROOT_DIR)
# docker-build:
# 	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/Meshbits/shurli golang:latest go build -o "$(BINARY_UNIX)" -v
