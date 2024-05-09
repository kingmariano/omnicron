import asyncio
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from endpoints import router

loop = asyncio.get_event_loop()
print(f"Using event loop: {type(loop)}\n")
print(f"Current event loop policy: {asyncio.get_event_loop_policy()}")

app = FastAPI(title="omnicron backend server")

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
    return {"Hello": "World"}



