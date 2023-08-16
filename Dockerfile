FROM debian:bookworm-slim as builder
WORKDIR /build

COPY . ./

RUN apt update -y && apt install -y golang libvips libvips-dev libde265-0 libde265-dev \
 && go build -ldflags="-s -w" -trimpath -o mediaproxy main.go

FROM debian:bookworm-slim

COPY --from=builder /build/mediaproxy /app/mediaproxy

RUN apt update -y \
 && apt install -y tini libvips libde265-0 libjemalloc2 \
 && groupadd -g 981 app \
 && useradd -d /app -s /bin/sh -u 981 -g app app \
 && chown -R app:app /app \
 && chmod +x /app/mediaproxy \
 && chmod -R 777 /app


USER app

ENV LD_PRELOAD=/usr/lib/x86_64-linux-gnu/libjemalloc.so.2
CMD ["tini", "--", "/app/mediaproxy"]
