# AABS Milvus Test CLI

## Store an Embedding

go run main.go store "Trump is bad. Trump is terrible. Trump is awful." 

## Verify Similar Meanings

go run main.go verify "Trump is horrible" 

## Store and Verify

go run main.go all "Trump is bad. Trump is terrible. Trump is awful." 