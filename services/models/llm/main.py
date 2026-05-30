import json
import os
import re
from typing import List, Optional

from fastapi import FastAPI
from pydantic import BaseModel
from huggingface_hub import hf_hub_download
from llama_cpp import Llama

MODEL_REPO = "Qwen/Qwen3-8B-GGUF"
MODEL_FILE = "Qwen3-8B-Q4_K_M.gguf"
MODEL_PATH = f"/models/{MODEL_FILE}"

if not os.path.exists(MODEL_PATH):
    print("Downloading Qwen3-8B...")
    hf_hub_download(
        repo_id=MODEL_REPO,
        filename=MODEL_FILE,
        local_dir="/models",
        token=os.environ.get("HF_TOKEN"),
    )

app = FastAPI(title="AABS LLM Service")

print("Loading Qwen3-8B...")

llm = Llama(
    model_path=MODEL_PATH,
    n_ctx=4096,
    n_threads=8,
    verbose=False,
    chat_format="chatml",
)

print("Qwen3-8B loaded")


class NameClusterRequest(BaseModel):
    posts: List[str]
    system_prompt: Optional[str] = None
    user_prompt: Optional[str] = None
    temperature: float = 0.2
    max_tokens: int = 64


class NameClusterResponse(BaseModel):
    name: str
    raw: str


class GenerateRequest(BaseModel):
    system_prompt: Optional[str] = None
    user_prompt: str
    temperature: float = 0.2
    max_tokens: int = 128


class GenerateResponse(BaseModel):
    text: str


@app.get("/health")
def health():
    return {
        "status": "ok",
        "model": MODEL_FILE,
    }


@app.post("/generate", response_model=GenerateResponse)
def generate(request: GenerateRequest):
    system_prompt = (
        request.system_prompt
        or "You are a helpful AI assistant."
    )

    result = llm.create_chat_completion(
        messages=[
            {
                "role": "system",
                "content": system_prompt.strip(),
            },
            {
                "role": "user",
                "content": request.user_prompt.strip(),
            },
        ],
        temperature=request.temperature,
        max_tokens=request.max_tokens,
    )

    text = result["choices"][0]["message"]["content"].strip()

    return GenerateResponse(
        text=text,
    )

@app.post("/name-cluster", response_model=NameClusterResponse)
def name_cluster(request: NameClusterRequest):
    sample_posts = "\n".join(
        f"- {post}"
        for post in request.posts[:10]
    )

    system_prompt = request.system_prompt or """
You are a semantic clustering service.

Your task:
- Read a group of social media posts.
- Produce a short neutral cluster name.

Rules:
- Return ONLY valid JSON.
- No explanations.
- No markdown.
- No reasoning.
- No extra text.

Valid response:

{"name":"Crypto Promotion"}
"""

    user_prompt = request.user_prompt or f"""
Posts:

{sample_posts}

Return only:

{{"name":"cluster name"}}
"""

    result = llm.create_chat_completion(
        messages=[
            {
                "role": "system",
                "content": system_prompt.strip(),
            },
            {
                "role": "user",
                "content": user_prompt.strip(),
            },
        ],
        temperature=request.temperature,
        max_tokens=request.max_tokens,
        response_format={
            "type": "json_object"
        },
    )

    raw = result["choices"][0]["message"]["content"].strip()

    try:
        parsed = json.loads(raw)

        name = sanitize_name(
            parsed.get(
                "name",
                "Unnamed Cluster",
            )
        )

    except Exception:
        raw = normalize_json_like_output(raw)
        name = extract_name(raw)

    return NameClusterResponse(
        name=name,
        raw=raw,
    )
    
def normalize_json_like_output(raw: str) -> str:
    cleaned = raw.strip()

    if not cleaned:
        return '{"name":"Unnamed Cluster"}'

    if "{" in cleaned:
        cleaned = cleaned[cleaned.find("{"):]

    if "}" in cleaned:
        cleaned = cleaned[: cleaned.find("}") + 1]

    if cleaned.startswith('"name"'):
        cleaned = "{" + cleaned

    if cleaned.startswith("{") and "}" not in cleaned:
        cleaned = cleaned + "}"

    return cleaned.strip()


def extract_name(raw: str) -> str:
    cleaned = raw.strip()

    try:
        parsed = json.loads(cleaned)
        name = parsed.get("name", "").strip()
        if name:
            return sanitize_name(name)
    except Exception:
        pass

    match = re.search(r'"name"\s*:\s*"([^"]+)"', cleaned)
    if match:
        return sanitize_name(match.group(1))

    first_line = cleaned.splitlines()[0] if cleaned else "Unnamed Cluster"
    first_line = first_line.replace("Cluster name:", "")
    first_line = first_line.replace("Label:", "")
    first_line = first_line.replace("Name:", "")

    return sanitize_name(first_line)


def sanitize_name(name: str) -> str:
    name = name.strip()
    name = name.replace('"', "").replace("'", "")

    banned_starts = r"^(okay|sure|here is|heres|let'?s see|the user wants|i need to|we need to|this cluster).*"
    name = re.sub(
        banned_starts,
        "Unnamed Cluster",
        name,
        flags=re.I,
    )

    name = re.sub(r"[^a-zA-Z0-9À-ÿ ,\-]", "", name)
    name = " ".join(name.split())

    if not name:
        return "Unnamed Cluster"

    words = name.split()

    if len(words) > 8:
        name = " ".join(words[:8])

    return name