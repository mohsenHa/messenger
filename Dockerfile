FROM golang:1.22 as build
#FROM busybox as build
COPY ./src/ /home/app/src/
# Set working directory
WORKDIR /home/app/src
RUN go build  -o /home/app/build/app ./cmd/httpserver/

FROM debian:stable-slim as run
COPY --from=build /home/app/build/app /app
COPY --from=build /home/app/src/config.yml /config.yml
ENTRYPOINT ["/app"]