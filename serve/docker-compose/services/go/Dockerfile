ARG GO_IMAGE_NAME
ARG GO_IMAGE_VERSION

FROM ${GO_IMAGE_NAME}:${GO_IMAGE_VERSION}

# Set working directory
WORKDIR /home/app/src

RUN go install github.com/air-verse/air@v1.52.2

ENTRYPOINT /home/app/serve/serve.bash