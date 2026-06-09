package embeddings

import "context"

type Vector []float32

// Embedder represents an embedder
type Embedder interface {
	Embed(ctx context.Context, text string) (Vector, error)
}
