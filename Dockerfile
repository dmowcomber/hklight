FROM golang:1.13.3-alpine3.10

WORKDIR /go/src/github.com/dmowcomber/hklight
COPY . /go/src/github.com/dmowcomber/hklight

RUN go build -mod=vendor .

CMD "./hklight"
