# Lint code. Reference: https://golang.org/cmd/vet/
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods \
         -nilfunc -printf -rangeloops -shift -structtags -unsafeptr

# list only the project's .go files i.e. exlcudes any .go files from the vendor 
# directory
GOFILES_NOVENDOR = $(shell find ./e2e -type f -name '*.go' -not -path "./e2e/vendor/*")

# Specify the name for the e2e binary
CTLNAME=m-e2e

.PHONY: all
all: format lint deps test

.PHONY: init
init: clean setup

.PHONY: setup
setup:
	@echo "--> Running setup"
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/Masterminds/glide

.PHONY: clean
clean:
	@echo "--> Running clean"
	rm -rf ./e2e/vendor/

.PHONY: format
format:
	@echo "--> Running go fmt"
	@go fmt $(GOFILES_NOVENDOR)

.PHONY: lint
lint:
	@echo "--> Running golint"
	@golint $(GOFILES_NOVENDOR)
	@echo "--> Running go vet"
	@go vet $(GOFILES_NOVENDOR)

.PHONY: deps
deps:
	@echo "--> Updating dependencies"
	@cd ./e2e && @glide install -v

# This is the main command of this project.
.PHONY: test
test:
	@echo "--> Running tests"
	@go test -v --kubeconfig=$(HOME)/.kube/config $(GOFILES_NOVENDOR)

.PHONY: image
image:
	@cd ./e2e/deploy/docker && sudo docker build -t openebs/m-e2e:ci .
	@sh ./e2e/buildscripts/push.sh
