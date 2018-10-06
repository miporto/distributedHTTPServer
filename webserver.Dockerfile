FROM golang:alpine
RUN go version

COPY . /go/src/github.com/manuporto/distributedHTTPServer/
WORKDIR /go/src/github.com/manuporto/distributedHTTPServer/

RUN go install -v ./...

EXPOSE 8080

CMD ["webserver", ":8080"]