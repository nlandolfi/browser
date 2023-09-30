TOOL_NAME := bro
TOOL_VERSION := 0.0.0
BUILT_AT_RFC3339 := $(shell date +"%Y-%m-%dT%H:%M:%S%z" | sed 's/\([0-9]\{2\}\)\([0-9]\{2\}\)$$/\1:\2/')

# Determine if Git repository is clean and get the corresponding SHA
GIT_CLEAN := $(shell if [ -z "$$(git status --porcelain)" ]; then echo "yes"; else echo "no"; fi)
ifeq ($(GIT_CLEAN),yes)
GIT_SHA := $(shell git rev-parse HEAD)
else
GIT_SHA := XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
endif

.PHONY: install build

install:
	@go build -ldflags "-X main.ToolVersion=$(TOOL_VERSION) -X main.GitSHA=$(GIT_SHA) -X main.BuiltAtRFC3339=$(BUILT_AT_RFC3339)" -o $(TOOL_NAME) *.go
	@mv $(TOOL_NAME) $$GOPATH/bin

# distribution directory, build writes binaries and .tar.gz files here
DIST_DIR := dist

build:
	@# macOS Intel
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.ToolVersion=$(TOOL_VERSION) -X main.GitSHA=$(GIT_SHA) -X main.BuiltAtRFC3339=$(BUILT_AT_RFC3339)" -o $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).darwin-amd64
	@tar -czf $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).darwin-amd64.tar.gz $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).darwin-amd64

	@# macOS ARM
	@GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.ToolVersion=$(TOOL_VERSION) -X main.GitSHA=$(GIT_SHA) -X main.BuiltAtRFC3339=$(BUILT_AT_RFC3339)" -o $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).darwin-arm64
	@tar -czf $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).darwin-arm64.tar.gz $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).darwin-arm64

	@# Linux Intel
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X main.ToolVersion=$(TOOL_VERSION) -X main.GitSHA=$(GIT_SHA) -X main.BuiltAtRFC3339=$(BUILT_AT_RFC3339)" -o $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).linux-amd64
	@tar -czf $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).linux-amd64.tar.gz $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).linux-amd64

	@# Linux ARM 32-bit
	@GOOS=linux GOARCH=arm go build -ldflags "-X main.ToolVersion=$(TOOL_VERSION) -X main.GitSHA=$(GIT_SHA) -X main.BuiltAtRFC3339=$(BUILT_AT_RFC3339)" -o $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).linux-arm
	@tar -czf $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).linux-arm.tar.gz $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).linux-arm

	@# Linux ARM 64-bit
	@GOOS=linux GOARCH=arm64 go build -ldflags "-X main.ToolVersion=$(TOOL_VERSION) -X main.GitSHA=$(GIT_SHA) -X main.BuiltAtRFC3339=$(BUILT_AT_RFC3339)" -o $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).linux-arm64
	@tar -czf $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).linux-arm64.tar.gz $(DIST_DIR)/$(TOOL_NAME)$(TOOL_VERSION).linux-arm64

	@shasum -a 256 $(DIST_DIR)/*.tar.gz
