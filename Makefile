ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))/src

lint:
	cd $(ROOT) && which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.0)
	cd $(ROOT) && golangci-lint run --config=$(ROOT)/.golangci.yml $(ROOT)/...

test:
	cd $(ROOT) && go test ./...

format:
	cd $(ROOT) && which gofumpt || (go install mvdan.cc/gofumpt@latest)
	cd $(ROOT) && gofumpt -l -w $(ROOT)
	cd $(ROOT) && which gci || (go install github.com/daixiang0/gci@latest)
	cd $(ROOT) && gci write $(ROOT)
	cd $(ROOT) && which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.0)
	cd $(ROOT) && golangci-lint run --fix

build:
	cd $(ROOT)
	echo "Building stage"
	echo $(DOCKER_HUB_TOKEN) | docker login --username "$(DOCKER_HUB_USERNAME)" --password-stdin
	docker build . -t $(IMAGE_NAME):$(IMAGE_VERSION)-$(GITHUB_RUN_ID)
	docker push $(IMAGE_NAME):$(IMAGE_VERSION)-$(GITHUB_RUN_ID)
	docker logout

deploy:
	cd $(ROOT)
	echo "Deploy stage"
	echo ${GOLANG_VERSION}

