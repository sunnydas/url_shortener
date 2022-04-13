FROM golang:1.13.4-alpine3.10 as builder
WORKDIR /
COPY ./ .
RUN GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags="-w -s" -v

FROM alpine:3.10
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /url-shortener /
CMD ["/url-shortener"]