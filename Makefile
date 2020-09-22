# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
#GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
CGO_CFLAGS=$(shell env CGO_CFLAGS="-I$HOME/go/src/github.com/satindergrewal/saplinglib/src/")
CGO_LDFLAGS_DARWIN=$(shell env CGO_LDFLAGS="-L$HOME/go/src/github.com/satindergrewal/saplinglib/dist/darwin -lsaplinglib -framework Security")
CGO_LDFLAGS_WIN="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/win64 -lsaplinglib -lws2_32 -luserenv"
CGO_LDFLAGS_LINUX="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/linux -lsaplinglib -lpthread -ldl -lm"
CGO_CC_WIN="x86_64-w64-mingw32-gcc"
MKDIR_P=mkdir -p
GITCMD=git
ROOT_DIR=$(shell pwd)
BINARY_NAME=shurli
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_OSX=$(BINARY_NAME)_osx
BINARY_WIN=$(BINARY_NAME).exe
DIST_DIR=dist
DIST_OSX=shurli_osx
DIST_OSX_PATH=$(DIST_DIR)/$(DIST_OSX)
DIST_UNIX=shurli_unix
DIST_UNIX_PATH=$(DIST_DIR)/$(DIST_UNIX)
DIST_WIN=shurli_win
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

# OS condition reference link: https://gist.github.com/sighingnow/deee806603ec9274fd47
UNAME_S=$(shell uname -s)

all:
	@echo $(OSFLAG)

all: build
build:
	$(GITCMD) checkout dev
	$(GOBUILD) -o $(BINARY_NAME) -v
# test: 
#	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -rf $(DIST_DIR)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_OSX)
run:
	$(GITCMD) checkout dev
	$(GOBUILD) -o $(BINARY_NAME) -v 
	./$(BINARY_NAME) start
deps-linux:
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/linux -lsaplinglib -lpthread -ldl -lm" $(GOGET) -u github.com/satindergrewal/kmdgo
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/linux -lsaplinglib -lpthread -ldl -lm" $(GOGET) -u github.com/Meshbits/shurli
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/linux -lsaplinglib -lpthread -ldl -lm" $(GOGET) -u github.com/gorilla/mux
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/linux -lsaplinglib -lpthread -ldl -lm" $(GOGET) -u github.com/gorilla/websocket

deps-osx:
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/darwin -lsaplinglib -framework Security" $(GOGET) -u github.com/satindergrewal/kmdgo
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/darwin -lsaplinglib -framework Security" $(GOGET) -u github.com/Meshbits/shurli
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/darwin -lsaplinglib -framework Security" $(GOGET) -u github.com/gorilla/mux
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/darwin -lsaplinglib -framework Security" $(GOGET) -u github.com/gorilla/websocket


deps-win:
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/win64 -lsaplinglib -lws2_32 -luserenv" CC="x86_64-w64-mingw32-gcc" $(GOGET) -u github.com/satindergrewal/kmdgo
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/win64 -lsaplinglib -lws2_32 -luserenv" CC="x86_64-w64-mingw32-gcc" $(GOGET) -u github.com/Meshbits/shurli
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/win64 -lsaplinglib -lws2_32 -luserenv" CC="x86_64-w64-mingw32-gcc" $(GOGET) -u github.com/gorilla/mux
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/win64 -lsaplinglib -lws2_32 -luserenv" CC="x86_64-w64-mingw32-gcc" $(GOGET) -u github.com/gorilla/websocket

# Cross compilation
build-linux: deps-linux
	rm -rf $(DIST_UNIX_PATH)
	$(GITCMD) checkout dev
	$(MKDIR_P) $(DIST_UNIX_PATH)
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/linux -lsaplinglib -lpthread -ldl -lm" CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(DIST_UNIX_PATH)/$(BINARY_NAME) -v
	$(CP_AV) $(DIST_FILES) $(DIST_UNIX_PATH)
	$(CURL_DL) $(KMD_UNIX_URL)
	$(UNZIP) -o komodo_linux.zip -d $(DIST_UNIX_PATH)/assets
	$(CURL_DL) $(SUBATOMIC_UNIX_URL)
	$(UNZIP) -o subatomic_linux.zip -d $(DIST_UNIX_PATH)/assets
	$(RM_RFV) komodo_linux.zip subatomic_linux.zip
	cd $(DIST_UNIX_PATH); zip -r ../shurli_linux.zip *; ls -lha ../; pwd
	$(RM_RFV) $(DIST_UNIX_PATH)
	cd $(ROOT_DIR)
build-osx: deps-osx
	$(GITCMD) checkout dev
	$(MKDIR_P) $(DIST_OSX_PATH)
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/darwin -lsaplinglib -framework Security" CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(DIST_OSX_PATH)/$(BINARY_NAME) -v
	$(CP_AV) $(DIST_FILES) $(DIST_OSX_PATH)
	$(CURL_DL) $(KMD_OSX_URL)
	$(UNZIP) -o komodo_macos.zip -d $(DIST_OSX_PATH)/assets
	$(CURL_DL) $(SUBATOMIC_OSX_URL)
	$(UNZIP) -o subatomic_macos.zip -d $(DIST_OSX_PATH)/assets
	$(RM_RFV) komodo_macos.zip subatomic_macos.zip
	cd $(DIST_OSX_PATH); zip -r ../shurli_macos.zip *
	$(RM_RFV) $(DIST_OSX_PATH)
	cd $(ROOT_DIR)
build-win:
	$(GITCMD) checkout dev
	$(MKDIR_P) $(DIST_WIN_PATH)
	CGO_CFLAGS="-I$(HOME)/go/src/github.com/satindergrewal/saplinglib/src/" CGO_LDFLAGS="-L$(HOME)/go/src/github.com/satindergrewal/saplinglib/dist/win64 -lsaplinglib -lws2_32 -luserenv" CC="x86_64-w64-mingw32-gcc" CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(DIST_WIN_PATH)/$(BINARY_WIN) -v
	@echo ".\shurli.exe start" > $(DIST_WIN_PATH)/start_shurli.cmd
	@echo ".\shurli.exe stop" > $(DIST_WIN_PATH)/stop_shurli.cmd
	$(CP_AV) $(DIST_FILES) $(DIST_WIN_PATH)
	$(CURL_DL) $(KMD_SUBATOMIC_WIN_URL)
	$(UNZIP) -o subatomic_kmd_win.zip -d $(DIST_WIN_PATH)/assets
	$(RM_RFV) subatomic_kmd_win.zip
	cd $(DIST_WIN_PATH); zip -r ../shurli_win.zip *
	$(RM_RFV) $(DIST_WIN_PATH)
	cd $(ROOT_DIR)