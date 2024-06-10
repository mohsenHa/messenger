# Messenger

Simple end-to-end encrypted messenger.

## Getting Started
For start project follow these steps:

- Clone project
```shell
git clone https://github.com/mohsenHa/messenger.git
cd messenger
```
- Create .env file
```shell
cp serve/.env.compose.example serve/.env
```
if you want you can update variables.
In windows please check the `EOL` of the line it must be `Unix(LF)` otherwise your config not work.

- Start docker compose

Windows:
```shell
.\serve\docker-compose.bat up -d
```


Linux:
```shell
./serve/docker-compose.bash up -d
```

You must set below host in host file:

```
127.0.0.1 SERVICE_NAME.local
```

- Done.

Now you can test APIs with postman or use sample golang client. 
For run client run this command.

Windows:
```shell
.\serve\docker-compose.bat exec app go run ./cmd/client/ ARG1 ARG2
```

Linux:
```shell
./serve/docker-compose.bash exec app go run ./cmd/client/ ARG1 ARG2
```
Hint: 
- ARG1: this argument (default:user) is used to store user details in the file user.json.
- ARG2: this argument (default:127.0.0.1:8088) is the base url of the server.

You can run the client command in separate command line with different user file and start chat between together. 

## Useful details

### URLs
- Base API URL: `http://SERVICE_NAME.local`
- Profiler URL: `http://SERVICE_NAME.local/profiler/debug/pprof/`
- PHPMyAdmin URL: `http://SERVICE_NAME.local/db_manager/`
- Rabbitmq manager URL: `http://SERVICE_NAME.local/rabbitmq_manager/`
- Traefik dashboard URL: `http://127.0.0.1:8080`

Note: The `SERVICE_NAME` is configured in `.env` file.

### Credentials

#### RabbitMQ
```
username: guest
password: guest
```

#### PHPMyAdmin
```
username: root
password: password
```

### .ENV 

Example of .env file:

```
SERVICE_NAME=messenger
COMPOSE_PROJECT_NAME=messenger
#SERVE_APPLICATION=not_serve
SERVE_APPLICATION=do_serve
#SERVE_WITH_AIR=yes
SERVE_WITH_AIR=no
GO_IMAGE_NAME=golang
GO_IMAGE_VERSION=1.22
```   
- `SERVICE_NAME`: Used for URLs of the services
- `COMPOSE_PROJECT_NAME`: Used for compose project name
- `SERVE_APPLICATION`: Used for serve or not serve the application
- `SERVE_WITH_AIR`: Used for serve with air for development. It is a live reload for golang.
