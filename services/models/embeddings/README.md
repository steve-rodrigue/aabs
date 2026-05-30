# CURL Example

Generate an embedding vector for a piece of text using the BAAI/bge-m3 model.

curl -s -X POST http://localhost:8080/embed \   -H "Content-Type: application/json" \   -d '{     "text": "Trump is bad. Trump is terrible. Trump is awful."   }' | jq 

### Request:

json {   "text": "Trump is bad. Trump is terrible. Trump is awful." } 

### Response:

json {   "dimensions": 1024,   "embedding": [     0.01234,     -0.05432,     ...   ] } 

## Health Check

curl http://localhost:8080/health 

### Response:

json {   "status": "ok",   "model": "BAAI/bge-m3" } 