FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /stratego ./main.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates wget \
    && addgroup -S app && adduser -S app -G app
WORKDIR /
COPY --from=builder /stratego /stratego
USER app
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/stratego", "--server"]
