FROM golang:1.27.2-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS builder

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


FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6

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
