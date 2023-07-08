FROM golang:1.20-alpine3.18 as builder
WORKDIR /build

COPY . ./

RUN apk add build-base imagemagick imagemagick-libs imagemagick-dev \
 && go build -ldflags="-s -w" -trimpath -o mediaproxy main.go

FROM alpine:3.18

COPY --from=builder /build/mediaproxy /app/mediaproxy

RUN apk add ca-certificates tini imagemagick-libs --no-cache \
 && addgroup -g 909 app \
 && adduser -D -h /app -s /bin/sh -u 909 -G app app \
 && chown -R app:app /app \
 && chmod +x /app/mediaproxy \
 && chmod -R 777 /app


USER app
CMD ["tini", "--", "/app/mediaproxy"]
