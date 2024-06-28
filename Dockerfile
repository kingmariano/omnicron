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

# Install Go dependencies and required packages
RUN apk add --no-cache \
    go \
    ffmpeg \
    python3 \
    py3-pip \
    python3-dev \
    build-base

RUN go mod download
RUN go get -u ./...
RUN go mod vendor
RUN go mod tidy

# Build Go binary
RUN go build -o ./omnicron

# Set up Python virtual environment
RUN python3 -m venv /build/venv
ENV PATH="/build/venv/bin:$PATH"

# Install Python dependencies in virtual environment
COPY ./python/requirements.txt ./python/requirements.txt
RUN /build/venv/bin/pip install --upgrade pip
RUN /build/venv/bin/pip install --upgrade --no-cache-dir -r ./python/requirements.txt

# Remove the default uvloop
RUN /build/venv/bin/pip uninstall -y uvloop

# Deploy stage
FROM gcr.io/distroless/base-debian11

WORKDIR /app

# Copy the built Go binary
COPY --from=builder /build/omnicron ./omnicron

# Copy ffmpeg binary from the builder stage
COPY --from=builder /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=builder /usr/share/ffmpeg /usr/share/ffmpeg

# Copy Python virtual environment
COPY --from=builder /build/venv /app/venv

# Copy Python scripts
COPY ./python /app/python

# Define the health check command
HEALTHCHECK --interval=1m --timeout=10s --retries=10 \
  CMD curl -f $HEALTHCHECK_ENDPOINT || exit 1

# Expose port 9000
EXPOSE 9000

# Run the application
CMD ["/app/omnicron"]
