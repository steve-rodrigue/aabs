package embeddings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	domain_embeddings "github.com/steve-rodrigue/aabs/services/saas/domain/embeddings"
)

type embedder struct {
	endpoint string
	client   *http.Client
}

type embedRequest struct {
	Text string `json:"text"`
}

type embedResponse struct {
	Dimensions int       `json:"dimensions"`
	Embedding  []float32 `json:"embedding"`
}

func (embedder *embedder) Embed(
	ctx context.Context,
	text string,
) (domain_embeddings.Vector, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, ErrInvalidEmbeddingText
	}

	if embedder.endpoint == "" {
		return nil, ErrInvalidEmbeddingEndpoint
	}

	payload, err := json.Marshal(
		embedRequest{
			Text: text,
		},
	)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		embedder.endpoint+"/embed",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := embedder.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf(
			"%w: status %d",
			ErrInvalidEmbeddingResponse,
			response.StatusCode,
		)
	}

	var result embedResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Dimensions <= 0 ||
		len(result.Embedding) == 0 ||
		result.Dimensions != len(result.Embedding) {
		return nil, ErrInvalidEmbeddingResponse
	}

	return domain_embeddings.Vector(result.Embedding), nil
}
