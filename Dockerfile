# Use the Ubuntu Slim base image
FROM ubuntu:22.04

WORKDIR /app

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
ENV PATH="/app/venv/bin:$PATH"

# Install necessary packages
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    python3-venv \
    gcc \
    g++ \
    musl-dev \
    build-essential \
    curl \
    cargo \
    ffmpeg \
    tesseract-ocr \
    && apt-get clean

# Extract Go version from go.mod and install Go
COPY go.mod go.sum ./

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

# Set up Python virtual environment and install dependencies
COPY ./python/requirements.txt ./python/requirements.txt
RUN python3 -m venv /app/venv && \
    /app/venv/bin/pip install --upgrade pip && \
    /app/venv/bin/pip install --upgrade --no-cache-dir -r ./python/requirements.txt && \
    /app/venv/bin/pip uninstall -y uvloop

# Verify installations
RUN python3 --version
RUN /app/venv/bin/pip --version
RUN go version
RUN ffmpeg -version
RUN tesseract --version

# Copy Python scripts
COPY ./python /app/python

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./omnicron

# Ensure the Go binary is executable
RUN chmod +x /app/omnicron

# Define the health check command
HEALTHCHECK --interval=1m --timeout=10s --retries=10 \
  CMD curl -f $HEALTHCHECK_ENDPOINT || exit 1

# Expose port 8000 for the FastAPI server
EXPOSE 8000

# Expose port 9000
EXPOSE 9000

# Run the application
ENTRYPOINT ["/app/omnicron"]
