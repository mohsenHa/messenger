#! /bin/bash
docker compose --env-file ./serve/docker-compose/.env -f ./serve/docker-compose/services/rabbitmq.yml -f ./serve/docker-compose/services/mysql.yml "$@"