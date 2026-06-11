FROM golang:1.26.4-alpine@sha256:7a3e50096189ad57c9f9f865e7e4aa8585ed1585248513dc5cda498e2f41812c AS builder

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


FROM alpine:3.23@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

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
