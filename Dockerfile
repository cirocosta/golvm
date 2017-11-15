FROM golang:alpine as builder

ADD ./VERSION /go/src/github.com/cirocosta/golvm/VERSION
ADD ./main.go /go/src/github.com/cirocosta/golvm/main.go
ADD ./vendor /go/src/github.com/cirocosta/golvm/vendor
ADD ./driver /go/src/github.com/cirocosta/golvm/driver
ADD ./lib /go/src/github.com/cirocosta/golvm/lib

WORKDIR /go/src/github.com/cirocosta/golvm

RUN set -ex && \
  CGO_ENABLED=0 go build \
        -tags netgo -v -a \
        -ldflags "-X main.version=$(cat ./VERSION) -extldflags \"-static\"" && \
  mv ./golvm /usr/bin/golvm

FROM busybox
COPY --from=builder /usr/bin/golvm /golvm

RUN set -x && \
  mkdir -p /var/log/golvm /mnt

CMD [ "golvm" ]
