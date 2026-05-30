package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	client "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type EmbeddingRequest struct {
	Text string `json:"text"`
}

type EmbeddingResponse struct {
	Dimensions int       `json:"dimensions"`
	Embedding  []float32 `json:"embedding"`
}

func main() {
	ctx := context.Background()

	milvus, err := client.NewGrpcClient(
		ctx,
		"localhost:19530",
	)
	if err != nil {
		panic(err)
	}
	defer milvus.Close()

	collectionName := "comments"

	exists, err := milvus.HasCollection(
		ctx,
		collectionName,
	)
	if err != nil {
		panic(err)
	}

	if !exists {
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
						"dim": "1024",
					},
				},
			},
		}

		err = milvus.CreateCollection(
			ctx,
			schema,
			2,
		)
		if err != nil {
			panic(err)
		}

		fmt.Println("Collection created")
	}

	request := EmbeddingRequest{
		Text: "Trump is bad. Trump is terrible. Trump is awful.",
	}

	body, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := httpClient.Post(
		"http://localhost:8080/embed",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)

		panic(fmt.Sprintf(
			"unexpected status code %d: %s",
			resp.StatusCode,
			string(responseBody),
		))
	}

	var embedding EmbeddingResponse

	if err := json.NewDecoder(resp.Body).Decode(&embedding); err != nil {
		panic(err)
	}

	textColumn := entity.NewColumnVarChar(
		"text",
		[]string{
			request.Text,
		},
	)

	vectorColumn := entity.NewColumnFloatVector(
		"embedding",
		1024,
		[][]float32{
			embedding.Embedding,
		},
	)

	_, err = milvus.Insert(
		ctx,
		collectionName,
		"",
		textColumn,
		vectorColumn,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Embedding stored in Milvus")
}
