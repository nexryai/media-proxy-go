FROM golang:1.20 as builder
WORKDIR /build

COPY . ./

RUN apt update && apt -y install libwebp-dev \
 && go build -o mediaproxy main.go

FROM debian:bullseye-slim

COPY --from=builder /build/mediaproxy /app/mediaproxy

RUN apt install -y tini libwebp6 --no-cache \
 && groupadd -g "991" misskey \
 && useradd -l -u "991" -g "991" -m -d /app app \
 && chown -R app:app /app \
 && chmod +x /app/mediaproxy \
 && chmod -R 777 /app


USER app
CMD ["tini", "--", "/app/mediaproxy"]