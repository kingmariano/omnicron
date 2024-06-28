FROM python:3.11-slim

# Set the working directory in the container
WORKDIR /app

# Copy the requirements file into the container
COPY requirements.txt .

# Install any needed packages specified in requirements.txt
RUN pip install --upgrade --no-cache-dir -r requirements.txt

# Install ffmpeg
RUN apt-get update && apt-get install -y ffmpeg

# Remove the default uvloop.
RUN pip uninstall -y uvloop

# Set build arguments for API keys
ARG MY_API_KEY
ARG GROK_API_KEY
ARG GEMINI_PRO_API_KEY
ARG UVLOOP_ENABLED
# Set environment variables from build arguments
ENV MY_API_KEY=$MY_API_KEY
ENV GROK_API_KEY=$GROK_API_KEY
ENV GEMINI_PRO_API_KEY=$GEMINI_PRO_API_KEY
ENV HEALTHCHECK_ENDPOINT=http://localhost:8000


COPY . .

# Define the health check command
HEALTHCHECK --interval=1m --timeout=10m --retries=10 \
  CMD curl -f $HEALTHCHECK_ENDPOINT || exit 1
  

# Expose port 8000
EXPOSE 8000

# Run uvicorn with the specified command
CMD ["python", "main.py"]