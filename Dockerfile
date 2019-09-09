FROM golang:1.12.4 as builder
WORKDIR /go/src/github.com/videocoin/cloud-pipelines
COPY . .
RUN make build

FROM bitnami/minideb:jessie
COPY --from=builder /go/src/github.com/videocoin/cloud-pipelines/bin/pipelines /opt/videocoin/bin/pipelines
CMD ["/opt/videocoin/bin/pipelines"]

