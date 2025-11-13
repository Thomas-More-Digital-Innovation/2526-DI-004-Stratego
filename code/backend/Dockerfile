# Multi-stage Dockerfile for the Go backend
FROM golang:1.25-alpine AS builder
WORKDIR /src

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /stratego ./main.go

# Final minimal runtime image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY --from=builder /stratego /stratego

EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/stratego", "--server"]
