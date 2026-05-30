FROM python:3.11-slim

WORKDIR /app

ARG HF_TOKEN

ENV PYTHONUNBUFFERED=1
ENV HF_HOME=/models
ENV TRANSFORMERS_CACHE=/models
ENV HF_TOKEN=${HF_TOKEN}

RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    build-essential \
    cmake \
    git \
    python3-dev \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .

RUN pip install --upgrade pip setuptools wheel
RUN pip install --no-cache-dir -r requirements.txt

RUN python -c "\
from huggingface_hub import hf_hub_download; \
hf_hub_download( \
    repo_id='Qwen/Qwen3-8B-GGUF', \
    filename='Qwen3-8B-Q4_K_M.gguf', \
    local_dir='/models', \
    token='${HF_TOKEN}' \
)"

COPY . .

EXPOSE 8100

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8100"]