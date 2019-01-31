FROM golang:1.11-alpine3.8

RUN apk update && apk add git

ENV GOPATH /go

RUN go get -v github.com/Lupino/yuntan-search

FROM alpine:3.8

COPY --from=0 /go/bin/yuntan-search /usr/bin/yuntan-search

ENTRYPOINT ["yuntan-search"]
