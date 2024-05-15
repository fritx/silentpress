# SilentPress

> [SilentPress](https://github.com/fritx/silentpress) is yet another Wiki, Blog & CMS framework, based on [silent](https://github.com/fritx/silent).

## v.s. VitePress & WordPress

| | Silent | SilentPress | VitePress | WordPress |
| :-- | :--: | :--: | :--: | :--: |
| Markdown first | √ | √ | √ |  |
| Static first | √ | √ | √ |  |
| Build-stage free | √ | √ |  | √ |
| CMS admin |  | √ | [🔧](https://vitepress.dev/guide/cms) | √ |
| Access control |  | √ |  | √ |
| Database free | √ | √ | √ |  |
| Soooo simple | √ | √ |  |  |

## Live Demo

- Wiki Home: https://fritx.me/silentpress/
- Wiki Admin: https://fritx.me/silentpress/admin
  - (Username: `admin`, Password: `SilentPress`)

## Deploy via Docker-Compose

```yml
# docker-compose.yml
services:
  silentpress:
    image: fritx/silentpress
    volumes:
      # - ~/silentpress/.env.example:/app/.env
      # - ~/silentpress/p_example:/app/p
      # set your own config & secrets, see below..
      - ./path/to/.env:/app/.env
      - ./path/to/p:/app/p
    environment:
      HOST: '0.0.0.0'  # required here
    ports:
      - '127.0.0.1:8080:8080'
```

## Build from Source

```sh
# Prepare
cp -r p_example p
cp .env.example .env
# set your own config & secrets
# vim .env
# for example
sed -i.bak "s|^COOKIE_SECRET=.*|COOKIE_SECRET=\"$(openssl rand -base64 32)\"|" .env
sed -i.bak "s|^ADMIN_PASSWORD=.*|ADMIN_PASSWORD=\"$(openssl rand -base64 32)\"|" .env

# Install dependencies
(cd silent && git stash -u)
git submodule update --init --recursive
(cd silent && git apply ../silent.patch)
go mod download

# Develop
go run .

# Develop with live reload
go install github.com/cosmtrek/air@latest
pkill -f silentpress/tmp/main; \
  air --build.exclude_dir "p,p_example,silent,silent_ext,static"

# Build & Deploy
go build && ./silentpress

# Deploy via PM2
pm2 start pm2.json && pm2 log

# Develop via Docker-Compose
docker compose up

# Push to Docker-Hub
docker login
docker push fritx/silentpress

# Save silent patch if changed
(cd silent && git add -A && git diff --cached > ../silent.patch && git reset .)
```
