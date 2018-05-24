FROM golang:1.8
ADD . /go/src/github.com/efy/placeholder

WORKDIR /go/src/github.com/efy/placeholder

RUN go get -u github.com/golang/dep/...
RUN dep ensure

RUN go build -o /go/bin/placeholder cmds/server/*.go

FROM alpine:latest
COPY --from=0 /go/bin/placeholder .
ENV PORT 8080
CMD ["/go/bin/placeholder"]
