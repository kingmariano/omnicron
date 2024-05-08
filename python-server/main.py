from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from endpoints import router
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



