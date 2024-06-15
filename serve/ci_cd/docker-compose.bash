#! /bin/bash
docker compose --env-file ./serve/docker-compose/.env -f ./serve/ci_cd/services/rabbitmq.yml -f ./serve/ci_cd/services/mysql.yml "$@"