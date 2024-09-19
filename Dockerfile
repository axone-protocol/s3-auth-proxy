#--- Build stage
FROM golang:1.23-bullseye AS go-builder

WORKDIR /src

COPY . /src/

RUN make build CGO_ENABLED=0

#--- Image stage
FROM alpine:3.20.3

COPY --from=go-builder /src/target/dist/s3-auth-proxy /usr/bin/s3-auth-proxy

WORKDIR /opt

ENTRYPOINT ["/usr/bin/s3-auth-proxy"]
