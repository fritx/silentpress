FROM golang:alpine
WORKDIR /app

COPY go.* .
RUN go mod download

COPY *.go .
RUN go build

COPY silent ./silent
COPY silent_ext ./silent_ext
COPY static ./static

CMD ["./silentpress"]
