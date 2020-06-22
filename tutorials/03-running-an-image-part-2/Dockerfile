FROM golang:latest
EXPOSE 8080
COPY build /go/src/build
WORKDIR /go/src/build/
RUN go build .
ENTRYPOINT ["/go/src/build/build"]