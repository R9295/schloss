FROM golang:1.18-alpine3.14 as builder

WORKDIR /app 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .

FROM scratch

WORKDIR /app

COPY --from=builder /app/schloss /usr/bin/

ENTRYPOINT ["schloss"]
