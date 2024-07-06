# Stage 1: Build stage with Ubuntu
FROM ubuntu:22.04 AS builder

WORKDIR /build

# Set build arguments for API keys
ARG MY_API_KEY
ARG GROK_API_KEY
ARG GEMINI_PRO_API_KEY
ARG REPLICATE_API_TOKEN
ARG CLOUDINARY_URL
ARG YOUTUBE_DEVELOPER_KEY
ARG TESSERACT_PREFIX
# Set environment variables from build arguments
ENV MY_API_KEY=${MY_API_KEY}
ENV GROK_API_KEY=${GROK_API_KEY}
ENV GEMINI_PRO_API_KEY=${GEMINI_PRO_API_KEY}
ENV REPLICATE_API_TOKEN=${REPLICATE_API_TOKEN}
ENV CLOUDINARY_URL = ${CLOUDINARY_URL}
ENV YOUTUBE_DEVELOPER_KEY=${YOUTUBE_DEVELOPER_KEY}
ENV TESSERACT_PREFIX=${TESSERACT_PREFIX}
ENV PORT=9000
ENV HEALTHCHECK_ENDPOINT=http://localhost:${PORT}/api/v1/readiness

# Install necessary packages
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    python3-venv \
    gcc \
    g++ \
    musl-dev \
    ffmpeg \
    build-essential \
    curl \
    cargo \
    tesseract-ocr \
    && apt-get clean

# Copy Go module files
COPY go.mod go.sum ./

# Extract Go version from go.mod
RUN grep '^go ' go.mod | awk '{print $2}' > goversion.txt
RUN curl -OL https://golang.org/dl/go$(cat goversion.txt).linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$(cat goversion.txt).linux-amd64.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go

# Verify Go installation
RUN go version

# Download Go dependencies
RUN go mod download

# Copy the entire project directory
COPY . .

# Build the Go application
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

# Stage 2: Final stage with Alpine for a smaller image
FROM alpine:latest

WORKDIR /app

# Install runtime packages
RUN apk add --no-cache ffmpeg curl

# Copy the built Go binary
COPY --from=builder /build/omnicron ./omnicron

# Copy ffmpeg binary from the builder stage
COPY --from=builder /usr/bin/ffmpeg /usr/bin/ffmpeg
COPY --from=builder /usr/share/ffmpeg /usr/share/ffmpeg

# Copy Tesseract binaries from the builder stage
COPY --from=builder /usr/bin/tesseract /usr/bin/tesseract
COPY --from=builder /usr/share/tesseract-ocr /usr/share/tesseract-ocr

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
