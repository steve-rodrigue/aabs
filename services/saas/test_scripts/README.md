# AABS Campaign Detection CLI

## Store a Post

Generate an embedding and store it in Milvus.

go run main.go store "Trump is bad. Trump is terrible. Trump is awful." 

---

## Verify Similar Meanings

Search Milvus for semantically similar posts and identify the closest campaign.

go run main.go verify "Trump is horrible" 

---

## Store and Verify

Store a post and immediately search for similar meanings.

go run main.go all "Trump is bad. Trump is terrible. Trump is awful." 

---

## Build Campaign Clusters

Run HDBSCAN on all stored embeddings, generate campaign names using the LLM, and store campaign clusters in PostgreSQL.

go run main.go cluster 

---

## Index Multiple Posts

Create multiple posts belonging to different campaigns.

go run main.go store "Trump is bad"

go run main.go store "Trump is terrible"

go run main.go store "Trump is awful"

go run main.go store "Buy Bitcoin now"

go run main.go store "BTC is going to the moon"

go run main.go store "Best crypto investment ever"

go run main.go store "Subscribe to my page"

go run main.go store "Check my profile for exclusive content"

go run main.go store "New content available now"

---

## Generate Campaign Clusters

go run main.go cluster 

Example output:

text Campaign cluster stored  Cluster ID: 0 Name: Anti-Trump Sentiment  Campaign cluster stored  Cluster ID: 1 Name: Crypto Promotion  Campaign cluster stored  Cluster ID: 2 Name: OnlyFans Promotion 

---

## Search for a Campaign Related to a Post

Given a new post, find the closest campaign.

go run main.go verify "Trump is horrible" 

Example output:

text Input text: Trump is horrible  Closest stored meanings:  Score: 0.97 Text: Trump is terrible  Score: 0.96 Text: Trump is awful  Score: 0.95 Text: Trump is bad 

The application can then determine that the post most likely belongs to the campaign:

text Anti-Trump Sentiment 

---

## Complete Workflow

go run main.go store "Trump is bad"

go run main.go store "Trump is terrible"

go run main.go store "Trump is awful"

go run main.go store "Buy Bitcoin now"

go run main.go store "BTC is going to the moon"

go run main.go store "Best crypto investment ever"

go run main.go store "Subscribe to my page"

go run main.go store "Check my profile for exclusive content"

go run main.go store "New content available now"

go run main.go cluster

go run main.go verify "Trump is horrible"

go run main.go verify "Invest in Bitcoin today"

go run main.go verify "Exclusive content in my profile"

---

## Architecture

- Posts
- Embeddings Service (BGE-M3)
- Milvus Vector Database
- HDBSCAN Clustering
- Representative Post Selection
- Campaign Clusters (PostgreSQL)
- LLM Naming (Qwen3)
- Generated Campaign Names
  - Anti-Trump Sentiment
  - Crypto Promotion
  - OnlyFans Promotion
  - AI-Generated Political Talking Points

---

## Goal

The objective is to automatically discover coordinated campaigns by:

1. Converting posts into embeddings.
2. Storing embeddings in Milvus.
3. Grouping semantically similar posts with HDBSCAN.
4. Identifying representative posts for each cluster.
5. Naming each cluster with an LLM.
6. Storing campaign metadata in PostgreSQL.
7. Associating new posts with existing campaigns through semantic similarity search.
8. Continuously refining campaign definitions as new posts are discovered.

---

## Example

Input posts:

- Trump is bad
- Trump is terrible
- Trump is awful

- Buy Bitcoin now
- BTC is going to the moon
- Best crypto investment ever

- Subscribe to my page
- Check my profile for exclusive content
- New content available now

Discovered campaigns:

- Anti-Trump Sentiment
- Crypto Promotion
- OnlyFans Promotion

New post:

- Trump is horrible

Matched campaign:

- Anti-Trump Sentiment
- Confidence: 98.7%

---

## Long-Term Vision

The system automatically builds a continuously evolving map of online narratives, promotional campaigns, coordinated messaging, bot activity, spam networks, and influence operations.

Rather than relying on predefined keywords or rules, campaigns emerge naturally from semantic similarity, allowing previously unknown narratives and coordinated behaviors to be detected automatically.

As the dataset grows, campaign clusters become increasingly accurate and can be tracked over time, enabling the detection of emerging narratives, coordinated influence operations, marketing campaigns, spam networks, and AI-generated content at scale.
```
