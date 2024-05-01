# SilentPress

> SilentPress is yet another Wiki, Blog & CMS framework, based on [silent](https://github.com/fritx/silent).

## v.s. VitePress & WordPress

| | Silent | SilentPress | VitePress | WordPress |
| :-- | :--: | :--: | :--: | :--: |
| Markdown first | √ | √ | √ |  |
| Static first | √ | √ | √ |  |
| Build-stage free | √ | √ |  | √ |
| CMS admin |  | √ | <a target="_blank" href="https://vitepress.dev/guide/cms">🔧</a> | √ |
| Database free | √ | √ | √ |  |
| Soooo simple | √ | √ |  |  |

## Live Demo

- Wiki Home: https://fritx.me/silentpress/
- Wiki Admin: https://fritx.me/silentpress/admin
  - (Username: `admin`, Password: `SilentPress`)

## Deploy via Docker-Commpose

```yml
# docker-compose.yml
services:
  silentpress:
    image: fritx/silentpress  # coming soon..
    restart: unless-stopped
    volumes:
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
sed -i.bak "s/^COOKIE_SECRET=.*/COOKIE_SECRET=\"$(openssl rand -base64 32)\"/" .env
sed -i.bak "s/^ADMIN_PASSWORD=.*/ADMIN_PASSWORD=\"$(openssl rand -base64 32)\"/" .env

# Install dependencies
git submodule update --init --recursive
go mod download

# Develop
go run .

# Develop with live reload
go install github.com/cosmtrek/air@latest
air --build.exclude_dir "p,p_example,silent,silent_ext,static"

# Build & Deploy
go build && ./silentpress

# Deploy via PM2
pm2 start pm2.json && pm2 log

# Develop via Docker-Compose
docker compose up
```
