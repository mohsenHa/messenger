#!/bin/bash

set -o allexport
source /home/app/serve/docker-compose/.env
set +o allexport

if [ $SERVE_APPLICATION = 'not_serve' ]; then
   echo "NOT SERVE"
  while [ true ]; do date; sleep 5; done
else
  if [ $SERVE_WITH_AIR = 'yes' ]; then
    echo "Serve with air"
    air -c /home/app/serve/docker-compose/services/go/.air.toml
  else
    echo "Build and serve"
    go build -o /home/app/src/cmd/httpserver/temp/app /home/app/src/cmd/httpserver/main.go
    /home/app/src/cmd/httpserver/temp/app
  fi
fi