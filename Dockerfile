FROM golang:1.20 as builder
WORKDIR /app

COPY . ./

RUN apt update && apt -y install libwebp-dev \
 && go build -o mediaproxy main.go

FROM alpine:latest

COPY --from=builder /app /app

RUN apk add tini libwebp --no-cache \
 && addgroup -g 1000 app \
 && adduser -D -h /app -s /bin/sh -u 1000 -G app app \
 && chown -R app:app /app \
 && chmod +x /app/mediaproxy \
 && chmod -R 777 /app \


USER app
CMD ["tini", "--", "/app/mediaproxy"]