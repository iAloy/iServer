ifndef CONFIG_FILE
CONFIG_FILE = config.json
endif

MAJOR       = $(shell jq -r .version.major ${CONFIG_FILE})
MINOR       = $(shell jq -r .version.minor ${CONFIG_FILE})
PATCH       = $(shell jq -r .version.patch ${CONFIG_FILE})
BUILD       = $(shell jq -r .version.build ${CONFIG_FILE})

VERSION = ${MAJOR}.${MINOR}.${PATCH}.${BUILD}

ARCH = amd64 arm64

all: clean
	@for i in $(ARCH);do echo "Building for $$i ..."; env GOOS=linux GOARCH=$$i go build  --ldflags="-s -w -X 'main.versionStr=${VERSION}' -buildid=" -tags "netgo" -o build/iServer_$$i main.go; done
	@echo "Done.";

clean:
	@rm -rf ./build