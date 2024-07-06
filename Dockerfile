# Use the Python 3.11 slim base image
FROM python:3.11-slim

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

# Install necessary packages
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    ffmpeg \
    tesseract-ocr \
    tesseract-ocr-eng \
    libtesseract-dev \
    libleptonica-dev \
    gcc \
    g++ \
    cargo \
    git \
    && rm -rf /var/lib/apt/lists/*

# Create and activate a virtual environment
# RUN python3 -m venv /app/venv
# ENV PATH="/app/venv/bin:$PATH"

# Upgrade pip and install Python dependencies
COPY ./python/requirements.txt ./python/requirements.txt
RUN pip install --upgrade pip
RUN pip install --upgrade --no-cache-dir -r ./python/requirements.txt

# Remove the default uvloop
RUN pip uninstall -y uvloop

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project directory
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./omnicron

# Verify installations
RUN python3 --version
RUN pip --version
RUN go version
RUN ffmpeg -version
RUN tesseract --version

# Copy .env file
COPY .env /app/.env

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
