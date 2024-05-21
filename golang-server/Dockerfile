# Build stage
FROM golang:1.22.2-alpine AS builder

WORKDIR /build
# Set build arguments for API keys
ARG MY_API_KEY
ARG GROK_API_KEY

# Set environment variables from build arguments
ENV MY_API_KEY=$MY_API_KEY
ENV GROK_API_KEY=$GROK_API_KEY
ENV PORT=9000
ENV HEALTHCHECK_ENDPOINT=http://localhost:${PORT}/api/v1/readiness

COPY . .

# Install Go dependencies
RUN go mod download

# Run tests
RUN CGO_ENABLED=0 go test ./tests

# Build Go binary
RUN go build -o ./omnicron

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 

WORKDIR /app

COPY --from=builder /build/omnicron ./omnicron

# Define the health check command
HEALTHCHECK --interval=1m --timeout=10m --retries=10 \
  CMD curl -f $HEALTHCHECK_ENDPOINT || exit 1

# Expose port 9000
EXPOSE 9000

# Run the application
CMD ["/app/omnicron"]

