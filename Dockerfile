FROM golang:1.20-bookworm as builder
WORKDIR /build

COPY . ./

RUN apt update && apt -y install libwebp-dev \
 && go build -o mediaproxy main.go

FROM debian:bookworm-slim

COPY --from=builder /build/mediaproxy /app/mediaproxy

RUN apt update \
 && apt install -y tini libwebp7 ca-certificates openssl \
 && groupadd -g "991" app \
 && useradd -l -u "991" -g "991" -m -d /app app \
 && chown -R app:app /app \
 && chmod +x /app/mediaproxy \
 && chmod -R 777 /app


USER app
CMD ["tini", "--", "/app/mediaproxy"]