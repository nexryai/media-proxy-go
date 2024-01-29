FROM alpine:edge as builder
WORKDIR /build

COPY . ./

RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.xtom.com.hk/alpine#g' /etc/apk/repositories \
 && apk add --no-cache ca-certificates go git alpine-sdk g++ build-base cmake clang libressl-dev vips vips-cpp vips-dev vips-heif \
 && go build -ldflags="-s -w" -trimpath -o mediaproxy main.go

FROM alpine:edge

COPY --from=builder /build/mediaproxy /app/mediaproxy

RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.xtom.com.hk/alpine#g' /etc/apk/repositories \
 && apk add --no-cache ca-certificates tini vips vips-heif \
 && addgroup -g 981 app \
 && adduser -u 981 -G app -D -h /app app \
 && chown -R app:app /app \
 && chmod +x /app/mediaproxy \
 && mkdir /cache \
 && chown -R app:app /cache

USER app
ENV CACHE_DIR=/cache

CMD ["tini", "--", "/app/mediaproxy"]
