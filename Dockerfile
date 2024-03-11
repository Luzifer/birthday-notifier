FROM golang:alpine as builder

COPY . /go/src/birthday-notifier
WORKDIR /go/src/birthday-notifier

RUN set -ex \
 && apk add --update git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath


FROM alpine:latest

LABEL maintainer "Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates \
      tzdata

COPY --from=builder /go/bin/birthday-notifier /usr/local/bin/birthday-notifier

ENTRYPOINT ["/usr/local/bin/birthday-notifier"]
CMD ["--"]

# vim: set ft=Dockerfile:
