# AABS LLM Service

## Health Check

Verify that the service is running and the model has loaded.

curl http://localhost:8100/health 

Expected response:

json {   "status": "ok",   "model": "Qwen3-8B-Q4_K_M.gguf" } 

---

## Name a Cluster

Ask the LLM to generate a semantic label for a group of related posts.

curl -X POST http://localhost:8100/name-cluster \   -H "Content-Type: application/json" \   -d '{     "posts": [       "Trump is bad",       "Trump is terrible",       "Trump is awful"     ]   }' 

Example response:

json {   "name": "Anti-Trump Sentiment" } 

---

## Crypto Example

curl -X POST http://localhost:8100/name-cluster \   -H "Content-Type: application/json" \   -d '{     "posts": [       "Buy Bitcoin now",       "BTC is going to the moon",       "Best crypto investment ever"     ]   }' 

Example response:

json {   "name": "Crypto Promotion" } 

---

## OnlyFans Example

curl -X POST http://localhost:8100/name-cluster \   -H "Content-Type: application/json" \   -d '{     "posts": [       "Check my profile for exclusive content",       "Subscribe to my page",       "New content available now"     ]   }' 

Example response:

json {   "name": "OnlyFans Promotion" } 

---

## View Logs

bash docker compose logs -f llm 

---

## Rebuild

bash docker compose build --no-cache llm docker compose up -d llm 

---

## Stop

bash docker compose stop llm 

---

## Restart

bash docker compose restart llm 