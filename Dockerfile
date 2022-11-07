FROM golang:1.19.3 AS builder

WORKDIR /build
COPY . .
RUN make build

FROM alpine

WORKDIR /etc/git2consul

COPY --from=builder /build/git2consul /usr/local/bin/git2consul
COPY --from=builder /build/config.sample.yaml /etc/git2consul/config.yaml

ARG TINI_ARCH=amd64
ARG TINI_VERSION=v0.19.0
ENV TINI_VERSION=${TINI_VERSION}

RUN apk add --update --no-cache curl ; \
        rm -rf /tmp/*; \
        rm -rf /var/cache/apk/*

RUN curl -LZo/tini \
    https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static-${TINI_ARCH}; \
    chmod +x /tini


ENTRYPOINT ["/tini", "--"]
CMD ["git2consul", "-config", "/etc/git2consul/config.yaml"]
