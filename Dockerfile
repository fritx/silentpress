FROM golang:alpine AS build
WORKDIR /app

COPY go.* .
RUN go mod download

COPY *.go .
RUN go build

FROM alpine
WORKDIR /app

COPY silent ./silent
COPY silent_ext ./silent_ext
COPY static ./static

COPY --from=build /app/silentpress .
CMD ["./silentpress"]
