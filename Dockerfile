FROM golang:1.23.5

WORKDIR ${GOPATH}/merch-shop/
COPY . ${GOPATH}/merch-shop/

RUN go build -o /build ./internal/cmd \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]