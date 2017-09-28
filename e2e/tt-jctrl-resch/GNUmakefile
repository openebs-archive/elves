# Lint code. Reference: https://golang.org/cmd/vet/
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods \
         -nilfunc -printf -rangeloops -shift -structtags -unsafeptr

# list only our namespaced directories
PACKAGES = $(shell go list ./... | grep -v '/vendor/')

.PHONY: all
all: format lint deps test image

.PHONY: init
init: clean setup

.PHONY: setup
setup:
	@echo "------------------"
	@echo "--> Running setup"
	@echo "------------------"
	@go get github.com/golang/lint/golint
	@go get golang.org/x/tools/cmd/goimports
	@go get -u github.com/golang/dep/cmd/dep

.PHONY: clean
clean:
	@echo "------------------"
	@echo "--> Running clean"
	@echo "------------------"
	@rm -rf ./app
	@rm -rf vendor/*/

.PHONY: format
format:
	@echo "------------------"
	@echo "--> Running go fmt"
	@echo "------------------"
	@go fmt $(PACKAGES)

.PHONY: lint
lint:
	@echo "------------------"
	@echo "--> Running golint"
	@echo "------------------"
	@golint $(PACKAGES)
	@echo "------------------"
	@echo "--> Running go vet"
	@echo "------------------"
	@go vet $(PACKAGES)

.PHONY: deps
deps:
	@echo "-------------------------"
	@echo "--> Updating dependencies"
	@echo "-------------------------"
	@dep ensure
	@echo "-------------------------"
	@echo "--> Status of dependencies"
	@echo "-------------------------"
	@dep status

# This is the main command of this project.
.PHONY: test
test:
	@echo "------------------"
	@echo "--> Running tests"
	@echo "------------------"
	@go test

.PHONY: image
image:
	@rm -rf ./app
	@GOOS=linux go build -o ./app .
	@sudo docker build -t openebs/m-e2e-tt-jctrl-resch:ci .
