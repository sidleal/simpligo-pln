FROM golang:1.9.2-alpine

RUN mkdir -p /go/src/github.com/sidleal/simpligo-pln
ADD . /go/src/github.com/sidleal/simpligo-pln
RUN go install github.com/sidleal/simpligo-pln
WORKDIR /go/src/github.com/sidleal/simpligo-pln

ENTRYPOINT /go/bin/simpligo-pln

EXPOSE 8080