# Stage 1: Build the Go application
FROM golang:latest AS builder
ARG GOARCH=amd64
WORKDIR /go/src/github.com/AntonioDaria/ltp_service
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Use the build argument GOARCH to build for the specified architecture
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -o app main.go

# Stage 2: Create the final minimal Docker image
FROM alpine:3.19

# Install necessary packages for healthcheck and certificates
RUN apk add --no-cache curl ca-certificates

EXPOSE 8000

# Copy the built application from the builder stage
COPY --from=builder /go/src/github.com/AntonioDaria/ltp_service/app /ltp_service.service

# Ensure the binary has execution permissions
RUN chmod +x /ltp_service.service

# Create a user and group 'app'
RUN addgroup -S app && adduser -S app -G app

# Temporarily comment out the USER directive for debugging
USER app
ENTRYPOINT ["/ltp_service.service"]
