# Build stage
FROM golang:1.22.4-alpine AS builder

WORKDIR /build
# Set build arguments for API keys
ARG MY_API_KEY
ARG GROK_API_KEY
ARG GEMINI_PRO_API_KEY
ARG REPLICATE_API_TOKEN
ARG CLOUDINARY_URL
ARG YOUTUBE_DEVELOPER_KEY
# Set environment variables from build arguments
ENV MY_API_KEY=${MY_API_KEY}
ENV GROK_API_KEY=${GROK_API_KEY}
ENV GEMINI_PRO_API_KEY=${GEMINI_PRO_API_KEY}
ENV REPLICATE_API_TOKEN=${REPLICATE_API_TOKEN}
ENV YOUTUBE_DEVELOPER_KEY=${YOUTUBE_DEVELOPER_KEY}
ENV PORT=9000
ENV HEALTHCHECK_ENDPOINT=http://localhost:${PORT}/api/v1/readiness

COPY . .

# Install Go and Python dependencies
RUN apk add --no-cache \
        ffmpeg \
        python3 \
        py3-pip \
        python3-dev \
        build-base && \
    go mod download && \
    go get -u ./... && \
    go mod vendor && \
    go mod tidy && \
    go build -o ./omnicron

# Install Python dependencies
COPY ./python/requirements.txt ./python/requirements.txt
RUN pip3 install --upgrade --no-cache-dir -r ./python/requirements.txt && \
    pip3 uninstall -y uvloop

# Deploy stage
FROM gcr.io/distroless/base-debian11

WORKDIR /app

# Copy the built Go binary
COPY --from=builder /build/omnicron ./omnicron

# Copy ffmpeg binary from the builder stage
COPY --from=builder /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=builder /usr/share/ffmpeg /usr/share/ffmpeg

# Copy Python and dependencies from the builder stage
COPY --from=builder /usr/local/lib/python3.11 /usr/local/lib/python3.11
COPY --from=builder /usr/local/bin/python3.11 /usr/local/bin/python3.11
COPY --from=builder /usr/local/bin/pip3 /usr/local/bin/pip3
COPY ./python /app/python

# Define the health check command
HEALTHCHECK --interval=1m --timeout=10s --retries=10 \
  CMD curl -f $HEALTHCHECK_ENDPOINT || exit 1

# Expose port 9000
EXPOSE 9000

# Run the application
CMD ["/app/omnicron"]
