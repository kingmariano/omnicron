# Copyright (c) 2024 Charles Ozochukwu

# Permission is hereby granted, free of charge, to any person obtaining a copy
#  of this software and associated documentation files (the "Software"), to deal
#  in the Software without restriction, including without limitation the rights
#  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
#  copies of the Software, and to permit persons to whom the Software is
#  furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in all
#  copies or substantial portions of the Software.

#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
#  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
#  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
#  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
#  SOFTWARE
"""
This module defines the FastAPI router for the omnicron python backend server.
"""

import os
import subprocess
import aiofiles
import asyncio
from pathlib import Path  # Standard library imports
from tempfile import NamedTemporaryFile
from typing import List
from fastapi.responses import JSONResponse
import requests  # Third-party imports
from fastapi import APIRouter, Depends, HTTPException, File, UploadFile
from fastapi.security.api_key import APIKeyHeader
from dotenv import load_dotenv
from pydantic import BaseModel
from shazamio import Shazam
from PIL import Image
import pytesseract
from utils.g4f_utils import get_chat_completion_response
from utils.doc_utils import process_page
# Local application imports

router = APIRouter()
env_path = Path(__file__).resolve().parents[2] / ".env"
load_dotenv(env_path)

# Get the API key from the environment
api_key = os.getenv('MY_API_KEY')
grok_api_key = os.getenv('GROK_API_KEY')
gemini_api_key = os.getenv('GEMINI_PRO_API_KEY')
tessdata_prefix = os.getenv('TESSDATA_PREFIX')
# Check if essential environment variables are set
if not api_key:
    raise RuntimeError("MY_API_KEY environment variable is not set.")
if not grok_api_key:
    raise RuntimeError("GROK_API_KEY environment variable is not set.")
if not gemini_api_key:
    raise RuntimeError("GEMINI_PRO_API_KEY environment variable is not set.")
if not tessdata_prefix:
    raise RuntimeError("TESSERACT_PREFIX environment variable is not set.")

os.environ['TESSDATA_PREFIX'] = tessdata_prefix
import  fitz

# Define the API key name
API_KEY_NAME = "Api-Key"
# Define a dependency to check for the API key in the headers
api_key_header = APIKeyHeader(name=API_KEY_NAME)

# Custom dependency to check for the API key in the headers


def check_api_key(key: str = Depends(api_key_header)):
    """
    Check if the API key is valid.
    """
    if key is None:
        raise HTTPException(status_code=401, detail="Provide Api Key")
    elif key != api_key:
        raise HTTPException(status_code=401, detail="Invalid Api Key")
    return key

# Define the request body models


class Message(BaseModel):
    """
    Represents a message.
    """
    role: str
    content: str


class ChatCompletionRequest(BaseModel):
    """
    Represents a chat completion request.
    """
    model: str = None
    messages: List[Message]
    stream: bool = False
    proxy: str = None
    timeout: int = None
    shuffle: bool = False
    image_url: str = None


class SearchSongRequest(BaseModel):
    """
    Represents an audio music request.
    """
    song: str
    limit: int = 6
    proxy: str = None


@router.post("/chat/completion")
async def chat_completion(request: ChatCompletionRequest, _: str = Depends(check_api_key)):
    """
    Perform chat completion using the GPT-3.5 and other state of the art model.

    Args:
        request (ChatCompletionRequest): The request object containing chat completion parameters.
        api_key (str, optional): The API key for authentication. Defaults to None.

    Returns:
        dict: A dictionary containing the chat completion response.
    """
    # Set API keys

    if request.image_url:
        try:
            print("image requested")
            image_response = requests.get(
                request.image_url, timeout=request.timeout)
            image_response.raise_for_status()
        except requests.exceptions.HTTPError as err:
            raise HTTPException(
                status_code=400, detail=f"Failed to download image: {err}") from err
        image_bytes = image_response.content
    else:
        image_bytes = None

    try:
        response = get_chat_completion_response(
            grok_api_key,
            gemini_api_key,
            request.messages,
            request.model,
            request.stream,
            request.proxy,
            request.timeout,
            request.shuffle,
            image_bytes
        )
        return {"response": response}
    except Exception as err:
        raise HTTPException(
            status_code=400, detail=f"Failed to complete chat: {err}") from err


@router.post("/shazam")
async def shazam_endpoint(file: UploadFile = File(...),
                          _: str = Depends(check_api_key)):
    """
    Recognize a song from an uploaded audio file.

    Args:
        file (UploadFile): The audio file to be recognized.
        _: The API key, validated by the check_api_key dependency.

    Returns:
        dict: The recognition result from Shazam.

    Raises:
        HTTPException: If an error occurs during processing.
    """
    try:
        mp3_bytes = await file.read()
        shazam = Shazam()
        out = await shazam.recognize(data=mp3_bytes)
        return out
    except Exception as err:
        raise HTTPException(status_code=500,
                            detail=f"Error processing shazam: {err}") from err


@router.post("/search-song")
async def search_song(
        request: SearchSongRequest,
        _: str = Depends(check_api_key)):
    """
    Search for a song using Shazam's search capabilities.

    Args:
        request (SearchSongRequest): The search request containing 
        the song query, limit, and optional proxy.
        _: The API key, validated by the check_api_key dependency.

    Returns:
        dict: The search result from Shazam.

    Raises:
        HTTPException: If an error occurs during the search.
    """
    try:
        shazam = Shazam()
        tracks = await shazam.search_track(query=request.song,
                                           limit=request.limit, proxy=request.proxy)
        return tracks
    except Exception as err:
        raise HTTPException(
            status_code=500,
            detail=f"Error searching song: {err}") from err
@router.post("/doc_analyze")
async def doc_analyze(file: UploadFile = File(...), _: str = Depends(check_api_key)):
    """
    Analyzes contents in a document file using the OCR engine if the document contains images and returns its text.
    Note: Supported file formats are PDF, XPS, EPUB, MOBI, FB2, CBZ, SVG, TXT.
    """
    # Save the uploaded file to a temporary location
    try:
        temp_file_path = Path(f"temp_{file.filename}")
        async with aiofiles.open(temp_file_path, "wb") as temp_file:
            await temp_file.write(await file.read())

        doc = fitz.open(str(temp_file_path))

        text_output = []
        for page in doc:
            text_output.append(await process_page(page))

        doc.close()
        return JSONResponse(content={"text": text_output})
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        if temp_file_path.exists():
            temp_file_path.unlink()


@router.post("/image_to_text")
async def image_to_text(file: UploadFile = File(...), _: str = Depends(check_api_key)):
    """
    Extract text from an uploaded image file using Tesseract OCR.

    Args:
        file (UploadFile): The image file to be processed.
        _: The API key, validated by the check_api_key dependency.

    Returns:
        dict: The extracted text from the image.

    Raises:
        HTTPException: If an error occurs during processing.
    """
    try:
        # Save the uploaded file to a temporary location
        temp_file_path = Path(f"temp_{file.filename}")
        async with aiofiles.open(temp_file_path, "wb") as temp_file:
            await temp_file.write(await file.read())

        # Open the image file using PIL
        image = Image.open(temp_file_path)

        # Use Tesseract to extract text from the image
        text = pytesseract.image_to_string(image)

        # Return the extracted text as JSON
        return JSONResponse(content={"text": text})
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    finally:
        # Delete the temporary file
        if temp_file_path.exists():
            temp_file_path.unlink()     


@router.get("/check_tesseract_installation_path")
async def check_tesseract_installation_path(_: str = Depends(check_api_key)):
    """
    Check the installation path of Tesseract and if the eng.traineddata file exists.
    """
    try:
        # Find the Tesseract installation path
        tesseract_path = subprocess.check_output(['which', 'tesseract']).decode().strip()

        # Check if the eng.traineddata file exists in the specified TESSDATA_PREFIX or common directories
        traineddata_paths = [
            Path(tessdata_prefix) / 'eng.traineddata',
            Path('/usr/share/tesseract-ocr/4.00/tessdata') / 'eng.traineddata',
            Path('/usr/share/tesseract-ocr/tessdata') / 'eng.traineddata',
            Path('/usr/local/share/tessdata') / 'eng.traineddata',
            Path('/usr/share/tessdata') / 'eng.traineddata',
            Path('/usr/share/tesseract/tessdata') / 'eng.traineddata',
            Path('/var/lib/tesseract-ocr/tessdata') / 'eng.traineddata'
        ]

        for path in traineddata_paths:
            if path.exists():
                eng_traineddata_path = path
                break
        else:
            raise FileNotFoundError("eng.traineddata not found in any known directories")

        return JSONResponse(content={"tesseract_path": tesseract_path, "eng_traineddata_path": str(eng_traineddata_path)})
    except subprocess.CalledProcessError:
        raise HTTPException(status_code=500, detail="Tesseract is not installed or not found in PATH")
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))