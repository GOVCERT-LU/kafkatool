FROM golang:1.16 AS toolchain

FROM debian:buster AS builder

# set the GOVCERT proxy
ENV HTTP_PROXY="http://proxy.int.govcert.etat.lu:8080"
ENV HTTPS_PROXY="http://proxy.int.govcert.etat.lu:8080"
ENV NO_PROXY=localhost,127.0.0.1,.govcert.etat.lu
ENV http_proxy="http://proxy.int.govcert.etat.lu:8080"
ENV https_proxy="http://proxy.int.govcert.etat.lu:8080"
ENV no_proxy=localhost,127.0.0.1,.govcert.etat.lu

ENV DEBIAN_FRONTEND noninteractive

COPY --from=toolchain /usr/local/go /usr/local/go

ENV PATH="/usr/local/go/bin:${PATH}"

RUN apt-get update -qq \
  && apt-get -y install --no-install-recommends \
  build-essential \
  ca-certificates 

RUN adduser builder --system --disabled-login \
  && mkdir /build \
  && chown builder: -R /build

COPY --chown=builder:root . /kafkatool

USER builder

WORKDIR /kafkatool

RUN cd src/kafkatool \
  && go build -ldflags="-s -w"

ENTRYPOINT ["/bin/bash"]

RUN mkdir -p /build/ \
  #
  && cp -ax /kafkatool/src/kafkatool/kafkatool /build/

