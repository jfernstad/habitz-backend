FROM arm64v8/golang:1.18.2-alpine3.15 as builder

RUN apk --no-cache add gcc g++

WORKDIR /go/src/github.com/jfernstad/habitz/web/
ADD ./internal ./internal
ADD ./cmd/backend ./cmd/backend
ADD ./vendor ./vendor
ADD ./go.mod ./go.mod
ADD ./go.sum ./go.sum

RUN CGO_ENABLED=1 GOOS=linux GOARH=arm64 go build -ldflags -a -installsuffix cgo -o app github.com/jfernstad/habitz/web/cmd/backend

#FROM alpine:latest
FROM arm64v8/alpine:3.15
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/
RUN mkdir -p cmd/backend/templates
RUN mkdir data
COPY --from=builder /go/src/github.com/jfernstad/habitz/web/app .

ENV SQLITE_DB data/habitz.sqlite
CMD ["./app"]

# docker run -p 3000:3000 -v ${PWD}:/root/data -d habitz:latest