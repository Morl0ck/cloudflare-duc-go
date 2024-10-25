FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ddns .

FROM alpine:latest as runner

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/ddns .

# Set environment variables (can be overridden by docker-compose)
ENV CF_API_TOKEN=""
ENV CF_ZONE_NAME=""
ENV CF_RECORD_NAME=""
ENV UPDATE_INTERVAL="5"

CMD ["./ddns"]