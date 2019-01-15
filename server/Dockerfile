FROM golang:1.11-alpine3.8 as builder

RUN apk add --update git curl libarchive-dev build-base linux-headers

RUN mkdir -p $$GOPATH/bin && \
    curl https://glide.sh/get | sh

ADD . /go/src/github.com/gustavosbarreto/uhmicro/server
WORKDIR /go/src/github.com/gustavosbarreto/uhmicro/server

RUN glide i && go build

FROM alpine:3.8

RUN apk add --update libarchive

COPY --from=builder \
     /go/src/github.com/gustavosbarreto/uhmicro/server/server \
     /usr/bin/updatehub-server

ENTRYPOINT ["/usr/bin/updatehub-server"]
