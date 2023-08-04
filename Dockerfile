FROM fedora:38 as builder
WORKDIR /build

COPY . ./

RUN dnf update -y && dnf install -y golang ImageMagick-devel \
 && go build -ldflags="-s -w" -trimpath -o mediaproxy main.go

FROM fedora:38

COPY --from=builder /build/mediaproxy /app/mediaproxy

RUN dnf update -y \
 && dnf -y install https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm \
 && dnf -y install https://download1.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm \
 && dnf update -y \
 && dnf -y install ffmpeg --allowerasing \
 && dnf update -y && dnf install -y tini ImageMagick-libs \
 && groupadd -g 981 app \
 && useradd -d /app -s /bin/sh -u 981 -g app app \
 && chown -R app:app /app \
 && chmod +x /app/mediaproxy \
 && chmod -R 777 /app


USER app
CMD ["tini", "--", "/app/mediaproxy"]
