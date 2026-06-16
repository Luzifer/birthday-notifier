FROM golang:1.26.4-alpine@sha256:f1ddd9fe14fffc091dd98cb4bfa999f32c5fc77d2f2305ea9f0e2595c5437c14 AS builder

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


FROM alpine:3.24@sha256:a2d49ea686c2adfe3c992e47dc3b5e7fa6e6b5055609400dc2acaeb241c829f4

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
