FROM golang:1.27.3-alpine AS builder
WORKDIR /build
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./ 
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/aibuddy ./cmd/app/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /build/aibuddy /app/aibuddy
CMD ["./aibuddy"]