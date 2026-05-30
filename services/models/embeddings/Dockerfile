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
    git \
    zlib1g-dev \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .

RUN pip install --upgrade pip setuptools wheel

RUN pip install --no-cache-dir torch==2.7.1

RUN pip install --no-cache-dir -r requirements.txt

RUN python -c "from FlagEmbedding import BGEM3FlagModel; BGEM3FlagModel('BAAI/bge-m3')"

COPY . .

EXPOSE 8080

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]