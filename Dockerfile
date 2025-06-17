FROM golang:1.22.0-alpine as builder
RUN apk add --no-cache ca-certificates git openssh-client

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o app ./cmd

FROM scratch
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/app /usr/bin/
ENTRYPOINT ["app"]