FROM python:3.11

# Set the working directory in the container
WORKDIR /app
 
# Copy the requirements file into the container
COPY requirements.txt .

# Install any needed packages specified in requirements.txt
RUN pip install --upgrade --no-cache-dir -r requirements.txt

# Install ffmpeg
RUN apt-get update && apt-get install -y ffmpeg

ENV MY_API_KEY = "HwikyHTrh4DMP9ZV1PzOn3K+C4Il7N/RKurD5AjyoIE="

ENV GROK_API_KEY = "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"

ENV GEMINI_PRO_API_KEY = "AIzaSyCWolXXFAdE1Or1O9UVprfclMVqJl74xSc"

ENV HEALTHCHECK_ENDPOINT=http://localhost:8080

ENV UVICORN_LOOP=asyncio

COPY . .

# Define the health check command
HEALTHCHECK --interval=30s --timeout=2m --retries=10 \
  CMD curl -f $HEALTHCHECK_ENDPOINT || exit 1

# Expose port 8080
EXPOSE 8080

# Run uvicorn with the specified command
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
