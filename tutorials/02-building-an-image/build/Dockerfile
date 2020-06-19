FROM golang:latest
COPY . /go/src/build/
WORKDIR /go/src/build/
RUN go build .
ENTRYPOINT ["/go/src/build/build"]