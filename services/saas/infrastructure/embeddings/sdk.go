package embeddings

import (
	"errors"
	"net/http"
	"strings"
	"time"

	domain_embeddings "github.com/steve-rodrigue/aabs/services/saas/domain/embeddings"
)

var (
	ErrInvalidEmbeddingText     = errors.New("invalid embedding text")
	ErrInvalidEmbeddingEndpoint = errors.New("invalid embedding endpoint")
	ErrInvalidEmbeddingResponse = errors.New("invalid embedding response")
)

// NewEmbedder creates a new embedding service embedder
func NewEmbedder(
	endpoint string,
) domain_embeddings.Embedder {
	return &embedder{
		endpoint: strings.TrimRight(endpoint, "/"),
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}
