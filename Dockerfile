FROM golang:1.27.3-alpine@sha256:b8c9b82affbbb82a020bbcf998dad2ba9516944ac7e883524642a481a9505e9c AS builder

COPY . /go/src/birthday-notifier
WORKDIR /go/src/birthday-notifier

RUN set -ex \
 && mkdir -p /rootfs/usr/bin \
 && apk add --update \
      git \
 && go build \
      -ldflags "-s -w -X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath \
      -o /rootfs/usr/bin/birthday-notifier


FROM alpine:3.24.3@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a

LABEL org.opencontainers.image.authors="Knut Ahlers <knut@ahlers.me>" \
      org.opencontainers.image.url="https://git.luzifer.io/registry/-/packages/container/birthday-notifier" \
      org.opencontainers.image.source="https://git.luzifer.io/luzifer/birthday-notifier" \
      org.opencontainers.image.title="birthday-notifier"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates \
      tzdata

COPY --from=builder /rootfs/ /

USER 1000:1000

ENTRYPOINT ["/usr/bin/birthday-notifier"]
