FROM golang:1.17.6 AS builder

WORKDIR /build
COPY . .
RUN make build

FROM alpine

COPY --from=builder /build/git2consul /git2consul

ENTRYPOINT ["/git2consul"]
