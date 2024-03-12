#--- Build stage
FROM golang:1.21-bullseye AS go-builder

WORKDIR /src

COPY . /src/

RUN make build CGO_ENABLED=0

#--- Image stage
FROM alpine:3.19.1

COPY --from=go-builder /src/target/dist/minio-auth-plugin /usr/bin/minio-auth-plugin

WORKDIR /opt

ENTRYPOINT ["/usr/bin/minio-auth-plugin"]
