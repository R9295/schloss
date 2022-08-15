FROM golang:1.18-alpine3.14 as builder

WORKDIR /app 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .

FROM alpine:latest

WORKDIR /app
RUN apk add git --no-cache

COPY --from=builder /app/schloss /usr/bin/

ENTRYPOINT ["schloss"]
