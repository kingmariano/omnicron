# import asyncio
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from endpoints import router
import uvicorn

# loop = asyncio.get_event_loop()
# print(f"Using event loop: {type(loop)}\n")
# print(f"Current event loop policy: {asyncio.get_event_loop_policy()}")

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
    return {"Hello": "World"}

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
