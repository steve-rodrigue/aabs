package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type EmbeddingRequest struct {
	Text string `json:"text"`
}

type EmbeddingResponse struct {
	Dimensions int       `json:"dimensions"`
	Embedding  []float32 `json:"embedding"`
}

func main() {
	request := EmbeddingRequest{
		Text: "Trump is bad. Trump is terrible. Trump is awful.",
	}

	body, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Post(
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

	fmt.Printf("Dimensions: %d\n", embedding.Dimensions)
	fmt.Printf("Embedding length: %d\n", len(embedding.Embedding))

	if len(embedding.Embedding) > 10 {
		fmt.Printf(
			"First 10 values: %v\n",
			embedding.Embedding[:10],
		)
	}
}
