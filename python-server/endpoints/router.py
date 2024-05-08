import os
import requests
from fastapi import APIRouter, Depends, HTTPException, File, UploadFile
from pathlib import Path
from dotenv import load_dotenv
from fastapi.security.api_key import APIKeyHeader
from typing import List
from pydantic import BaseModel
from shazamio import Shazam

from utils.g4f_utils import get_chat_completion_response, handle_image_processing

router = APIRouter()
env_path = Path(__file__).resolve().parents[2] / ".env"
load_dotenv(env_path)

# Get the API key from the environment
api_key = os.getenv('MY_API_KEY') 
grok_api_key = os.getenv('GROK_API_KEY')
gemini_api_key = os.getenv('GEMINI_PRO_API_KEY')
# Define the API key name
API_KEY_NAME = "Api-Key"
# Define a dependency to check for the API key in the headers
api_key_header = APIKeyHeader(name=API_KEY_NAME, auto_error=False)

# Custom dependency to check for the API key in the headers
def check_api_key(api_key: str = Depends(api_key_header)):
    if api_key is None or api_key != api_key:
        raise HTTPException(status_code=401, detail="Invalid API key")
    return api_key
# Define the request body model
class Message(BaseModel):
    role: str
    content: str

class ChatCompletionRequest(BaseModel):
    model: str = None
    messages: List[Message]
    stream: bool = False
    proxy: str = None
    api_key: str = grok_api_key
    timeout: int = None
    shuffle: bool = False
    
class VisionRequest(BaseModel):
    messages: List[Message]
    api_key: str = gemini_api_key
    proxy: str = None
    image_url: str
    timeout: int = None
    model: str  = "gemini-pro-vision"
    stream: bool = False
class SearchSongRequest(BaseModel):
    song: str
    limit: int = 6
    proxy: str = None    
@router.get("/endpoint")
def get_endpoint(api_key: str = Depends(check_api_key)):
    return {"message": "Hello, World!"}

@router.post("/chat/completion")
async def chat_completion(request: ChatCompletionRequest, api_key: str = Depends(check_api_key)):
    try:
        response = get_chat_completion_response(request.messages, request.api_key, request.proxy, request.stream, request.timeout, request.model, request.shuffle)
        return {"response": response}
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e)) 

@router.post("/vision")
async def vision_endpoint(request: VisionRequest, api_key: str = Depends(check_api_key)):
     # Download the image from the URL
    try:
        image_response = requests.get(request.image_url)
        image_response.raise_for_status()
    except requests.exceptions.HTTPError as err:
        raise HTTPException(status_code=400, detail=f"Failed to download image: {err}")

    # Use the downloaded image as a file-like object
    image_bytes = image_response.content

    # Create the response using g4f.ChatCompletion.create
    try:
        response = handle_image_processing(image_bytes, request.messages, request.api_key, request.proxy, request.stream, request.timeout, request.model)
        return response
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error processing image: {e}")


@router.post("/shazam")
async def shazam(file: UploadFile = File(...), api_key: str = Depends(check_api_key)):
    try:
          
        mp3_bytes = await file.read()
        shazam = Shazam()
        out = await shazam.recognize(data=mp3_bytes)
        return out  
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error processing shazam: {e}")
        
@router.post("/search-song")
async def search_song(request: SearchSongRequest, api_key: str = Depends(check_api_key)):
    try:
        shazam = Shazam()
        tracks = await shazam.search_track(query=request.song, limit=request.limit, proxy=request.proxy)
        return tracks
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error searching song: {e}")