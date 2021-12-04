

FROM golang:1.17-alpine AS builder

WORKDIR go/src/github.com/payments-api
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
COPY . .
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o /app services/payment-processor/cmd/main.go
#
# Create final image
FROM scratch
WORKDIR /
COPY --from=builder app app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY services/payment-processor/config.yaml config.yaml
EXPOSE 50051
CMD ["./app"]