FROM golang:1.24.3 AS toolchain


FROM debian:bookworm AS builder


ENV DEBIAN_FRONTEND=noninteractive

COPY --from=toolchain /usr/local/go /usr/local/go

ENV PATH="/usr/local/go/bin:${PATH}"

RUN apt-get update -qq \
  && apt-get -y install --no-install-recommends \
  build-essential \
  ca-certificates 

RUN adduser builder --system --disabled-login --home /build \
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

