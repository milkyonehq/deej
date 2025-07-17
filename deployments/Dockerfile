FROM golang:1.24-alpine AS build

ENV CGO_ENABLED=1

WORKDIR /app

RUN apk update && apk add --no-cache \
    build-base \
    opus-dev \
    pkgconf

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd

COPY pkg ./pkg

RUN go build -o deej ./cmd/deej

FROM alpine:3.22

ARG UID="1001"
ARG GID="1001"

WORKDIR /app

RUN apk update && apk add --no-cache \
    ffmpeg \
    python3 \
    py3-pip

RUN pip install --upgrade --break-system-packages \
    yt-dlp

RUN addgroup -g $GID deej && adduser -u $UID -G deej -D deej

RUN chown deej:deej /app

COPY --from=build /app/deej .

USER deej

ENTRYPOINT ["./deej"]
