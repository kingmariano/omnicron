"""
This module defines the FastAPI router for the omnicron python backend server.
"""

import os
from pathlib import Path  # Standard library imports
from typing import List

import requests  # Third-party imports
from fastapi import APIRouter, Depends, HTTPException, File, UploadFile
from fastapi.security.api_key import APIKeyHeader
from dotenv import load_dotenv
from pydantic import BaseModel
from shazamio import Shazam

from utils.g4f_utils import get_chat_completion_response  # Local application imports

router = APIRouter()
env_path = Path(__file__).resolve().parents[2] / ".env"
load_dotenv(env_path)

# Get the API key from the environment
api_key = os.getenv('MY_API_KEY')
grok_api_key = os.getenv('GROK_API_KEY')
gemini_api_key = os.getenv('GEMINI_PRO_API_KEY')
print(api_key)
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
