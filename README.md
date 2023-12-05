# SilentPress

```sh
# Prepare
cp -r p_example p
cp .env.example .env
vim .env  # set your own config & secrets

# Develop
go run .

# Develop with live reload
go install github.com/cosmtrek/air@latest
air --build.exclude_dir "p,p_example,silent,silent_ext,static"

# Build & Deploy
go build && ./silentpress

# Deploy via Docker
docker run --env-file .env xxx

# Deploy via Docker-Compose
# See below
```

## Deploy via Docker-Commpose

```yml
# docker-compose.yml
services:
  silentpress:
    # todo
```
