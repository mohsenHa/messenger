ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))/src

lint:
	cd $(ROOT) && which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.0)
	cd $(ROOT) && golangci-lint run --config=$(ROOT)/.golangci.yml $(ROOT)/...

test:
	cd $(ROOT) && go test ./...

format:
	cd $(ROOT) && @which gofumpt || (go install mvdan.cc/gofumpt@latest)
	cd $(ROOT) && @gofumpt -l -w $(ROOT)
	cd $(ROOT) && @which gci || (go install github.com/daixiang0/gci@latest)
	cd $(ROOT) && @gci write $(ROOT)
	cd $(ROOT) && @which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.0)
	cd $(ROOT) && @golangci-lint run --fix

build:
	cd $(ROOT)
	echo "Building stage"
	echo ${GOLANG_VERSION}

deploy:
	cd $(ROOT)
	echo "Deploy stage"
	echo ${GOLANG_VERSION}

