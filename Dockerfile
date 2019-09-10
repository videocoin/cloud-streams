FROM golang:1.12.4 as builder
WORKDIR /go/src/github.com/videocoin/cloud-streams
COPY . .
RUN make build

FROM bitnami/minideb:jessie
COPY --from=builder /go/src/github.com/videocoin/cloud-streams/bin/streams /opt/videocoin/bin/streams
CMD ["/opt/videocoin/bin/streams"]

