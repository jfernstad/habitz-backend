FROM golang:1.18 as builder-arm64
WORKDIR /go/src/github.com/jfernstad/habitz/web/
ADD ./internal ./internal
ADD ./cmd/backend ./cmd/backend
ADD ./vendor ./vendor

RUN CGO_ENABLED=1 GOOS=linux GOARH=arm64 go build  -ldflags="-extldflags=-static" -tags sqlite_omit_load_extension -a -installsuffix cgo -o app github.com/jfernstad/habitz/web/cmd/backend

#FROM alpine:latest
# FROM frolvlad/alpine-glibc
FROM arm64v8/alpine:3.15
RUN apk --no-cache add ca-certificates sqlite

# RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# removing apk cache
# RUN rm -rf /var/cache/apk/*

WORKDIR /root/
RUN mkdir -p cmd/backend/templates
RUN mkdir data
COPY --from=builder-arm64 /go/src/github.com/jfernstad/habitz/web/app .

ENV SQLITE_DB data/habitz.sqlite
CMD ["./app"]

# docker run -v ${PWD}:/root/data habitz:latest