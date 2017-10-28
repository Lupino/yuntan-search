FROM docker.io/phusion/baseimage:0.9.21
MAINTAINER Li Meng Jun <lmjubuntu@gmail.com>

RUN apt-get update && apt-get install -y git && \
    curl -o /tmp/go.tar.gz https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz && \
    tar xvf /tmp/go.tar.gz -C /usr/local

ENV GOPATH /root/go
ENV PATH /root/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

RUN go get github.com/Lupino/yuntan-search && \
    go get github.com/Lupino/tokenizer/tokenizer

WORKDIR /root
