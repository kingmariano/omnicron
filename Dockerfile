# Build stage
FROM golang:1.22 AS builder

WORKDIR /build

# Set build arguments for API keys
ARG MY_API_KEY
ARG GROK_API_KEY
ARG GEMINI_PRO_API_KEY
ARG REPLICATE_API_TOKEN
ARG CLOUDINARY_URL
ARG YOUTUBE_DEVELOPER_KEY
ARG TESSDATA_PREFIX
# Set environment variables from build arguments
ENV MY_API_KEY=${MY_API_KEY}
ENV GROK_API_KEY=${GROK_API_KEY}
ENV GEMINI_PRO_API_KEY=${GEMINI_PRO_API_KEY}
ENV REPLICATE_API_TOKEN=${REPLICATE_API_TOKEN}
ENV CLOUDINARY_URL=${CLOUDINARY_URL}
ENV YOUTUBE_DEVELOPER_KEY=${YOUTUBE_DEVELOPER_KEY}
ENV PORT=9000
# Set the TESSDATA_PREFIX environment variable
ENV TESSDATA_PREFIX=${TESSDATA_PREFIX}

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o omnicron

# Final stage
FROM python:3.12-slim

WORKDIR /app

COPY --from=builder /build/omnicron /app/

# copy pthon requiremts
COPY ./python/requirements.txt ./python/requirements.txt

RUN python3 -m venv /app/venv && \
    /app/venv/bin/pip install --upgrade pip && \
    /app/venv/bin/pip install --upgrade --no-cache-dir -r ./python/requirements.txt && \
    /app/venv/bin/pip uninstall -y uvloop

# copy python scripts
COPY ./python /app/python

# Install FFmpeg and other dependencies
RUN apt-get update && \
    apt-get install -y \
    ffmpeg \
    wget \
    unzip \
    bc \
    vim \
    python3-pip \
    libleptonica-dev \
    git \
    make \
    g++ \
    autoconf \
    automake \
    libtool \
    pkg-config \
    libpng-dev \
    libjpeg62-turbo-dev \
    libtiff5-dev \
    libicu-dev \
    libpango1.0-dev \
    autoconf-archive \
    poppler-utils && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install Tesseract OCR
RUN mkdir -p /app/src && cd /app/src && \
    wget https://github.com/tesseract-ocr/tesseract/archive/5.0.0-alpha-20201224.zip && \
    unzip 5.0.0-alpha-20201224.zip && \
    cd tesseract-5.0.0-alpha-20201224 && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install && \
    ldconfig && \
    make training && \
    make training-install && \
    cd /usr/local/share/tessdata && \
    wget https://github.com/tesseract-ocr/tessdata_best/raw/main/eng.traineddata && \
    cd /app && \
    rm -rf /app/src

# Verify installations    
RUN ffmpeg -version
RUN tesseract --version   

# Ensure the Go binary is executable
RUN chmod +x /app/omnicron

ENV PATH="/app/venv/bin:$PATH"

ENV HEALTHCHECK_ENDPOINT=http://localhost:${PORT}/api/v1/readiness

# Define  health check command
HEALTHCHECK --interval=1m --timeout=10s --retries=10 \
  CMD curl -f $HEALTHCHECK_ENDPOINT || exit 1

EXPOSE 8000 9000

ENTRYPOINT ["/app/omnicron"]