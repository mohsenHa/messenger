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
	echo $(DOCKER_HUB_TOKEN) | docker login --username "$(DOCKER_HUB_USERNAME)" --password-stdin
	mkdir ./src/key/
	touch ./src/key/key
	touch ./src/key/key.pub
	echo $(PRIVATE_KEY) > ./src/key/key
	echo $(PUBLIC_KEY) > ./src/key/key.pub
	docker build . -t $(IMAGE_NAME):$(IMAGE_VERSION)-$(GITHUB_RUN_ID) --build-arg GO_VERSION=$(GO_VERSION) --build-arg MESSENGER_APPLICATION_ENABLE__PROFILING=$(MESSENGER_APPLICATION_ENABLE__PROFILING) --build-arg MESSENGER_LOGGER_STORE__TO__FILE=$(MESSENGER_LOGGER_STORE__TO__FILE) --build-arg MESSENGER_RABBITMQ_USER=$(MESSENGER_RABBITMQ_USER) --build-arg MESSENGER_RABBITMQ_PASSWORD=$(MESSENGER_RABBITMQ_PASSWORD) --build-arg MESSENGER_RABBITMQ_HOST=$(MESSENGER_RABBITMQ_HOST) --build-arg MESSENGER_RABBITMQ_PORT=$(MESSENGER_RABBITMQ_PORT) --build-arg MESSENGER_RABBITMQ_VHOST=$(MESSENGER_RABBITMQ_VHOST) --no-cache
	docker push $(IMAGE_NAME):$(IMAGE_VERSION)-$(GITHUB_RUN_ID)
	docker logout

deploy:
	echo "Deploy stage"
	npm install -g @liara/cli
	liara deploy --image="$(IMAGE_NAME):$(IMAGE_VERSION)-$(GITHUB_RUN_ID)" --api-token="$(LIARA_API_TOKEN)" --app="$(APP_NAME)" --port="$(APP_PORT)" --no-app-logs


