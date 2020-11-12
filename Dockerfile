FROM golang:alpine

ENV GO1111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /build
ADD . /build
RUN go build ./cmd/orderPushScheduler
EXPOSE 50051
CMD ["./orderPushScheduler", "server"]

