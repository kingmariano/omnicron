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
ENV CLOUDINARY_URL=${CLOUDINARY_URL}
ENV YOUTUBE_DEVELOPER_KEY=${YOUTUBE_DEVELOPER_KEY}
ENV TESSERACT_PREFIX=${TESSERACT_PREFIX}
ENV PORT=9000
ENV HEALTHCHECK_ENDPOINT=http://localhost:${PORT}/api/v1/readiness

# Install necessary packages
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    gcc \
    g++ \
    musl-dev \
    build-essential \
    curl \
    cargo \
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

# Install Python dependencies
COPY ./python/requirements.txt ./python/requirements.txt
RUN pip3 install --upgrade pip
RUN pip3 install --upgrade --no-cache-dir -r ./python/requirements.txt --target /build/python-packages

# Remove the default uvloop
RUN pip3 uninstall -y uvloop

# Determine the Python version
RUN python3 --version 2>&1 | awk '{print $2}' | cut -d. -f1,2 > python-version.txt

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./omnicron

# Stage 2: Final stage with Alpine for a smaller image
FROM alpine:latest

WORKDIR /app

# Install runtime packages including tesseract-ocr
RUN apk add --no-cache ffmpeg curl python3 py3-pip tesseract-ocr

# Copy .env file
COPY .env /app/.env

# Copy the built Go binary
COPY --from=builder /build/omnicron ./omnicron

# Ensure the binary is executable
RUN chmod +x /app/omnicron

# Determine the Python version from the build stage
COPY --from=builder /build/python-version.txt /app/python-version.txt
RUN PYTHON_VERSION=$(cat /app/python-version.txt)

# Install Python dependencies
COPY --from=builder /build/python-packages /usr/local/lib/python$PYTHON_VERSION/site-packages

# Copy Python scripts
COPY ./python /app/python

# Ensure the correct interpreter is used
RUN ln -sf /usr/bin/python3 /usr/local/bin/python3
RUN ln -sf /usr/bin/pip3 /usr/local/bin/pip3

# Verify ffmpeg installation
RUN ffmpeg -version

# Verify Python installation
RUN python3 --version

# Verify Tesseract installation and version
RUN tesseract --version

# Define the health check command
HEALTHCHECK --interval=1m --timeout=10s --retries=10 \
  CMD curl -f $HEALTHCHECK_ENDPOINT || exit 1

# Expose port 8000 for the FastAPI server
EXPOSE 8000

# Expose port 9000
EXPOSE 9000

# Run the application
ENTRYPOINT ["/app/omnicron"]
