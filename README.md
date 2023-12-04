# bec-wiki

```sh
# Prepare
cp -r p_example p
cp .env.example .env
vim .env  # set your own config & secrets

# Develop
go run .

# Build & Deploy
go build && ./bec-wiki

# Deploy via Docker
docker run --env-file .env xxx

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
