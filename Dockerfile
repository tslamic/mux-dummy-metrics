FROM golang:1.12.5-alpine AS builder

WORKDIR /metrics
COPY ./ .
RUN apk update && apk add --no-cache git=~2.20 ca-certificates=~20190108
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o app .

FROM alpine:3.9.4
RUN addgroup -S -g 1710 metrics && adduser -S metrics -G metrics
USER metrics
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /metrics/app .

EXPOSE 8088

CMD ["./app"]
