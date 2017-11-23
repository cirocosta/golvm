FROM golang:alpine as builder

ADD ./VERSION /go/src/github.com/cirocosta/golvm/VERSION
ADD ./main.go /go/src/github.com/cirocosta/golvm/main.go
ADD ./lvmctl /go/src/github.com/cirocosta/golvm/lvmctl
ADD ./vendor /go/src/github.com/cirocosta/golvm/vendor
ADD ./driver /go/src/github.com/cirocosta/golvm/driver
ADD ./lib /go/src/github.com/cirocosta/golvm/lib

WORKDIR /go/src/github.com/cirocosta/golvm

RUN set -ex && \
  CGO_ENABLED=0 go build \
        -tags netgo -v -a \
        -ldflags "-X main.version=$(cat ./VERSION) -extldflags \"-static\"" && \
  mv ./golvm /usr/bin/golvm

FROM alpine
COPY --from=builder /usr/bin/golvm /golvm

RUN set -x && \
  apk add --update lvm2 cryptsetup util-linux e2fsprogs xfsprogs && \
  mkdir -p /var/log/golvm /mnt

CMD [ "golvm" ]
