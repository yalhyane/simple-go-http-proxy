# Variables
APP_NAME := simple-go-http-proxy
VERSION := 0.0.1
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GOOS_LIST := darwin linux windows
GOARCH_LIST := 386 amd64
PLATFORMS_LIST="windows/amd64/.exe" "windows/386/.exe" "darwin/amd64" "linux/amd64" "linux/386"

# Targets
build:
	@echo "Building $(APP_NAME) $(VERSION) for all platforms"
	@for PLATFORM in $(PLATFORMS_LIST); do \
  	  	echo 'Building binary for platform: '$$PLATFORM; \
  	  	PLATFORM_SPLIT=($${PLATFORM//\// }); \
  	  	GOOS=$${PLATFORM_SPLIT[0]}; \
  	  	GOARCH=$${PLATFORM_SPLIT[1]}; \
  	  	EXT=$${PLATFORM_SPLIT[2]-""}; \
  	  	OUTPUT_NAME=$(APP_NAME)'_'$(VERSION)'_'$$GOOS'-'$$GOARCH$$EXT;\
        GOOS=$$GOOS GOARCH=$$GOARCH go build -o dist/$$OUTPUT_NAME -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"; \
    done

clean:
	@echo "Cleaning up dist directory"
	@rm -rf dist/*

.PHONY: build clean
