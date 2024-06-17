ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))/src

lint:
	which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.0)
	golangci-lint run --config=$(ROOT)/.golangci.yml $(ROOT)/...

test:
	cd $(ROOT) && go test ./...

format:
	cd $(ROOT)
	@which gofumpt || (go install mvdan.cc/gofumpt@latest)
	@gofumpt -l -w $(ROOT)
	@which gci || (go install github.com/daixiang0/gci@latest)
	@gci write $(ROOT)
	echo "which golangci-lint"
	@which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.0)
	echo "golangci-lint run --fix"
	@golangci-lint run --fix

build:
	cd $(ROOT)
	echo "Building stage"
	echo ${GOLANG_VERSION}

deploy:
	cd $(ROOT)
	echo "Deploy stage"
	echo ${GOLANG_VERSION}

