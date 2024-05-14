"""
This module defines the main FastAPI application for the omnicron python backend server.
"""

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from endpoints import router
import uvicorn

app = FastAPI(title="omnicron python backend server")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.include_router(router, prefix="/api/v1")


@app.get("/")
def read_root():
    """
    Root endpoint returning a simple greeting.
    """
    return {"Hi": "World"}


if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8080, reload=True)
