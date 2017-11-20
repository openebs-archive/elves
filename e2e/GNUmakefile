# Lint code. Reference: https://golang.org/cmd/vet/
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods \
         -nilfunc -printf -rangeloops -shift -structtags -unsafeptr

# list only our namespaced directories
PACKAGES = $(shell go list ./... | grep -v '/vendor/')

.PHONY: all
all: format lint test

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
	@rm -rf ./tt-jctrl-resch/app
	@rm -rf ./demo-replica-count/app

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

.PHONY: test
test:
	@echo "------------------"
	@echo "--> Running tests"
	@echo "------------------"
	@go test $(PACKAGES)

.PHONY: tt-jctrl-resch-image
tt-jctrl-resch-image:
	@rm -rf ./tt-jctrl-resch/app
	@GOOS=linux go build -o ./tt-jctrl-resch/app ./tt-jctrl-resch
	@sudo docker build -t openebs/m-e2e-tt-jctrl-resch:ci ./tt-jctrl-resch

.PHONY: demo-replica-count-image
demo-replica-count-image:
	@rm -rf ./demo-replica-count/app
	@GOOS=linux go build -o ./demo-replica-count/app ./demo-replica-count
	@sudo docker build -t openebs/m-e2e-demo-replica-count:ci ./demo-replica-count

.PHONY: storagepool-policy
storagepool-policy-image:
	@rm -rf ./storagepool-policy/app
	@GOOS=linux go build -o ./storagepool-policy/app ./storagepool-policy
	@sudo docker build -t openebs/m-e2e-storagepool-policy:ci ./storagepool-policy
