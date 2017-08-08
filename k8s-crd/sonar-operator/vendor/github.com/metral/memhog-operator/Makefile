default: build

# Pkgs to build
PKGS = $(shell glide nv)

# Name of app/binary
APP_NAME=memhog-operator

# Go pkg path for app
GO_PKG_PATH=github.com/metral/$(APP_NAME)

# Binary output dir
OUTPUT_DIR = _output
OUTPUT_PATH = $(OUTPUT_DIR)/$(APP_NAME)

build:
	# Build binary
	go install $(PKGS)

clean:
	rm -rf $(APP_NAME)
	rm -rf $(OUTPUT_DIR)

static-build:
	# Static build for Linux x86_64
	$(shell mkdir -p $(OUTPUT_PATH)/linux_amd64)
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $(OUTPUT_PATH)/linux_amd64/$(APP_NAME) -a -installsuffix no_cgo -ldflags '-w -extld ld -extldflags -static' -x $(GO_PKG_PATH)
