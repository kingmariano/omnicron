## Omnicron

![Omnicron](omnicron_logo.jpg "omnicron")

[![Build Status](https://img.shields.io/github/actions/workflow/status/kingmariano/omnicron/ci.yml?branch=main)](https://github.com/kingmariano/omnicron/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/kingmariano/omnicron)](https://goreportcard.com/report/github.com/kingmariano/omnicron)
[![Release](https://img.shields.io/github/release/kingmariano/omnicron.svg?label=Release)](https://github.com/kingmariano/omnicron/releases)
[![License](https://img.shields.io/github/license/kingmariano/omnicron)](https://github.com/kingmariano/omnicron/blob/main/LICENSE)

## 📜Description

Omnicron is a powerful multimodel API designed to cater to a wide range of advanced capabilities. Think of it as an all in one package that suites all your needs. Here is a list of its features:

1. **Chat Completion**: Includes models like the open ai gpt-4 and gpt-40 to handle chat completion.
2. **Image Generation**: Image Generation is handled by [replicate](https://replicate.com/) state of the art models for image generations like the [sdxl-lightning-4step](https://replicate.com/bytedance/sdxl-lightning-4step) and more ...
3. **Video Generation**
4. **Speech To Text**
5. **Text To Speech (voice cloning)**
6. **Image Upscale**
7. **Music Generation**
8. **Song Recognition (shazam)**
9. **Youtube Summarization**
10. **Video downloader tool**: Download videos from various website.
11. **Music downloader tool**: Search and download song/music to your device. A drop in replacement for spotify.
12. **Doc GPT**: Chat with your pdf files.
13. **MP3 Converter**: Convert your video and audio files to mp3.
14. **Text To Image**: Scan text written in your image using the tesseract ocr engine.

### Efficiency and Model Strength

Omnicron leverages robust and efficient models from [Replicate](https://replicate.com/), ensuring high performance, accuracy, and scalability across all its functionalities. Rigorous testing and meticulous design using go makes Omnicron suitable for production environments.

### Projects Utilizing Omnicron

One notable project that successfully integrates Omnicron is [Omnicron-Telebot](https://github.com/kingmariano/omnicron-telebot), demonstrating its versatility and reliability in real-world applications.

> note: If you want your project to be listed here. Send an email to me with a full description of your project. charlesozochukwu2004@gmail.com

## Requirements

- Go 1.20+
- Python 3.8 or higher
- Git
- Docker and Docker-compose.
- Tesseract-ocr
- ffmpeg

## Installation

To install Omnicron, follow these steps:

1. **Clone the Repository**: Clone the Git repository to your local machine.

   ```sh
   git clone https://github.com/kingmariano/omnicron.git
   ```

2. **Set Up Environment Variables**: Change the environment variables as required.

   - **MY API KEY**: Generate a random API key. It can be done using various methods. Either by using openssl or any other method suitable for you.
     using openssl

   ```bash
   openssl rand -base64 32
   ```

   - **Port**: Set the port to `9000` or any other port number.

   - **GROK API KEY**: Signup/Login to the [Grok Console Cloud](https://console.groq.com/login) and retrieve your API token.

   - **GEMINI PRO API KEY**: Signup/Login to the [Google Gemini Studio](https://ai.google.dev/aistudio/) and retrieve your API token.

   - **REPLICATE API TOKEN**: Signup/Login to [replicate](https://replicate.com/) and retrieve your API token.

   - **CLOUDINARY URL**: Signup/Login to Cloudinary. Navigate to the [cloudinary console](https://console.cloudinary.com/) and get your cloudinary URL secret. It should look like this `cloudinary://<your_api_key>:<your_api_secret>@djagytapi`

   - **YOUTUBE DEVELOPER KEY**: Retrieve your youtube developer key. You can check this [post](https://blog.hubspot.com/website/how-to-get-youtube-api-key) on how to get it.

   - **TESSDATA PREFIX**: This is the location of where tesseract is installed on your machine. Install [tessract](https://tesseract-ocr.github.io/tessdoc/Installation.html) for your os and set the location to `/usr/local/share/tessdata` for linux or `C:\Program Files\Tesseract-OCR\tessdata` for windows. If yours is configured to a different file location. Set it to where the **tessdata** location is on your machine.

   - **FAST API BASE URL**: Set this to `http://0.0.0.0:8000` or `http://localhost:8000` to connect the fast api server to the go code.

   - **Environment Variables**: Create a `.env` file in the project root directory and add the following:

   ```env
   PORT=9000 // or some other port
   MY_API_KEY=YOUR_API_KEY_HERE
   GROK_API_KEY=YOUR_GROK_API_KEY_HERE
   GEMINI_PRO_API_KEY=YOUR_GEMINI_PRO_API_KEY_HERE
   REPLICATE_API_TOKEN=YOUR_REPLICATE_API_TOKEN_HERE
   CLOUDINARY_URL=YOUR_CLOUDINARY_URL_HERE
   YOUTUBE_DEVELOPER_KEY=YOUR_YOUTUBE_DEVELOPER_KEY_HERE
   TESSDATA_PREFIX=/usr/local/share/tessdata //or C:\Program  Files\Tesseract-OCR\tessdata for windows
   FAST_API_BASE_URL=http://0.0.0.0:8000
   ```

3. **Build and run the Application**:
   First install all the python dependencies by running

   ```sh
   pip install --upgrade -r ./python/requirements.txt
   ```

   then run

   ```sh
   go build -o omnicron &&  ./omnicron
   ```

4. **Using docker compose**: You can use docker compose to build the image.

   ```sh
   docker-compose up --build
   ```

## 💡Usage

After setting up and running the application, you can navigate to `http://localhost:9000/readiness` to check the health of the application.

**Check out the full api documentation [here](https://omnicron-docs.com)**

## Client Libraries📚

- Golang: A robust wrapper for the Omnicron API has also already been written check it out [here](https://github.com/kingmariano/omnicron-go)

# Contributing

Please feel free to submit issues, fork the repository and send pull requests! If you notice any bug in the code. You can always submit an issue and it will be reviewed.

---

## Citation

If you use Omnicron in your research or project, please cite it as follows:

```bibtex
@software{omnicron,
  author = {Charles Ozochukwu},
  title = {Omnicron: A Multimodel API for Advanced AI Capabilities},
  year = {2024},
  url = {https://github.com/kingmariano/omnicron}
}
```

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/kingmariano/omnicron/blob/main/LICENSE) file for details.
