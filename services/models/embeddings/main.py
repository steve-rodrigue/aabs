from fastapi import FastAPI
from pydantic import BaseModel
from FlagEmbedding import BGEM3FlagModel

app = FastAPI()

print("Loading BGE-M3...")
model = BGEM3FlagModel(
    "BAAI/bge-m3",
    use_fp16=False
)
print("BGE-M3 loaded")

class EmbeddingRequest(BaseModel):
    text: str

class EmbeddingResponse(BaseModel):
    dimensions: int
    embedding: list[float]

@app.post("/embed", response_model=EmbeddingResponse)
async def embed(request: EmbeddingRequest):
    result = model.encode(
        [request.text],
        batch_size=1,
        max_length=8192
    )

    vector = result["dense_vecs"][0].tolist()

    return EmbeddingResponse(
        dimensions=len(vector),
        embedding=vector
    )

@app.get("/health")
async def health():
    return {
        "status": "ok",
        "model": "BAAI/bge-m3"
    }