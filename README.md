# bec-wiki

```sh
# Develop
vim .env  # reset your own secret & password
go run .

# Build & Deploy
vim .env  # reset your own secret & password
go build && ./bec-wiki

# Deploy via Docker
docker run xxx -e xxx -e xxx

# Deploy via Docker-Compose
# See below
```

## Deploy via Docker-Commpose

```yml
# docker-compose.yml
services:
  bec-wiki:
    # todo
```
