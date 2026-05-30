package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	client "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	collectionName   = "comments"
	milvusAddress    = "localhost:19530"
	embeddingService = "http://localhost:8080/embed"
	vectorDimension  = 1024
)

type EmbeddingRequest struct {
	Text string `json:"text"`
}

type EmbeddingResponse struct {
	Dimensions int       `json:"dimensions"`
	Embedding  []float32 `json:"embedding"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	milvus, err := client.NewGrpcClient(ctx, milvusAddress)
	if err != nil {
		panic(err)
	}
	defer milvus.Close()

	command := os.Args[1]

	switch command {
	case "store":
		text := getTextArg()
		mustEnsureCollection(ctx, milvus)
		mustStoreEmbedding(ctx, milvus, text)

	case "verify":
		text := getTextArg()
		mustEnsureCollection(ctx, milvus)
		mustVerifyMeaning(ctx, milvus, text)

	case "all":
		text := getTextArg()
		mustEnsureCollection(ctx, milvus)
		mustStoreEmbedding(ctx, milvus, text)
		mustVerifyMeaning(ctx, milvus, text)

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run main.go store \"text to embed\"")
	fmt.Println("  go run main.go verify \"text to search\"")
	fmt.Println("  go run main.go all \"text to embed and search\"")
}

func getTextArg() string {
	if len(os.Args) >= 3 {
		return os.Args[2]
	}

	return "Trump is bad. Trump is terrible. Trump is awful."
}

func mustEnsureCollection(ctx context.Context, milvus client.Client) {
	exists, err := milvus.HasCollection(ctx, collectionName)
	if err != nil {
		panic(err)
	}

	if exists {
		mustEnsureIndex(ctx, milvus)
		return
	}

	schema := &entity.Schema{
		CollectionName: collectionName,
		Description:    "AABS comment embeddings",
		AutoID:         true,
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     true,
			},
			{
				Name:     "text",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "65535",
				},
			},
			{
				Name:     "embedding",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": fmt.Sprintf("%d", vectorDimension),
				},
			},
		},
	}

	if err := milvus.CreateCollection(ctx, schema, 2); err != nil {
		panic(err)
	}

	fmt.Println("Collection created:", collectionName)

	mustEnsureIndex(ctx, milvus)
}

func mustEnsureIndex(ctx context.Context, milvus client.Client) {
	index, err := entity.NewIndexFlat(entity.COSINE)
	if err != nil {
		panic(err)
	}

	err = milvus.CreateIndex(
		ctx,
		collectionName,
		"embedding",
		index,
		false,
	)

	if err != nil {
		fmt.Println("CreateIndex:", err)
		return
	}

	fmt.Println("Index created")
}

func mustFetchEmbedding(text string) EmbeddingResponse {
	request := EmbeddingRequest{
		Text: text,
	}

	body, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := httpClient.Post(
		embeddingService,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("unexpected status code %d: %s", resp.StatusCode, string(responseBody)))
	}

	var embedding EmbeddingResponse

	if err := json.NewDecoder(resp.Body).Decode(&embedding); err != nil {
		panic(err)
	}

	if embedding.Dimensions != vectorDimension {
		panic(fmt.Sprintf("unexpected embedding dimensions: got %d, expected %d", embedding.Dimensions, vectorDimension))
	}

	return embedding
}

func mustStoreEmbedding(ctx context.Context, milvus client.Client, text string) {
	embedding := mustFetchEmbedding(text)

	textColumn := entity.NewColumnVarChar(
		"text",
		[]string{text},
	)

	vectorColumn := entity.NewColumnFloatVector(
		"embedding",
		vectorDimension,
		[][]float32{embedding.Embedding},
	)

	_, err := milvus.Insert(
		ctx,
		collectionName,
		"",
		textColumn,
		vectorColumn,
	)
	if err != nil {
		panic(err)
	}

	if err := milvus.Flush(ctx, collectionName, false); err != nil {
		panic(err)
	}

	fmt.Println("Embedding stored in Milvus")
	fmt.Println("Text:", text)
	fmt.Println("Dimensions:", embedding.Dimensions)
}

func mustVerifyMeaning(ctx context.Context, milvus client.Client, text string) {
	if err := milvus.LoadCollection(ctx, collectionName, false); err != nil {
		panic(err)
	}

	embedding := mustFetchEmbedding(text)

	searchParam, err := entity.NewIndexFlatSearchParam()
	if err != nil {
		panic(err)
	}

	results, err := milvus.Search(
		ctx,
		collectionName,
		[]string{},
		"",
		[]string{"text"},
		[]entity.Vector{
			entity.FloatVector(embedding.Embedding),
		},
		"embedding",
		entity.COSINE,
		10,
		searchParam,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Input text:")
	fmt.Println(text)
	fmt.Println()
	fmt.Println("Closest stored meanings:")

	for _, result := range results {
		for i, score := range result.Scores {
			fmt.Printf("Score: %.4f\n", score)

			for _, field := range result.Fields {
				if field.Name() == "text" {
					texts := field.(*entity.ColumnVarChar).Data()
					fmt.Println("Text:", texts[i])
				}
			}

			fmt.Println()
		}
	}
}
