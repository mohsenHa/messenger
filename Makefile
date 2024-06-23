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
	echo "Building stage"
	cd $(ROOT)
	echo $(DOCKER_HUB_USERNAME)
	echo $(IMAGE_NAME)
	echo $(IMAGE_VERSION)
	echo $(GITHUB_RUN_ID)
	echo $(DOCKER_HUB_TOKEN) | docker login --username "$(DOCKER_HUB_USERNAME)" --password-stdin
	docker build . -t $(IMAGE_NAME):$(IMAGE_VERSION)-$(GITHUB_RUN_ID)
	docker push $(IMAGE_NAME):$(IMAGE_VERSION)-$(GITHUB_RUN_ID)
	docker logout

deploy:
	echo "Deploy stage"
	npm install -g @liara/cli
	liara deploy --image="$(IMAGE_NAME):$(IMAGE_VERSION)-$(GITHUB_RUN_ID)" --api-token="$(LIARA_API_TOKEN)" --app="$(APP_NAME)" --port="$(APP_PORT)" --no-app-logs


