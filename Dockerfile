FROM golang:1.23

WORKDIR ${GOPATH}/refstudy/
COPY . ${GOPATH}/refstudy/

RUN go mod download

RUN go build -o /build ./cmd \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]