FROM golang:1.20.2 AS build
WORKDIR /go/src/app
COPY . .
RUN GOOS=linux GOARCH=386 go build -v -o app .

FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /root/
COPY --from=build /go/src/app/app .
ENTRYPOINT ["./app"]
