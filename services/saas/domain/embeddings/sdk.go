package embeddings

type Vector []float32

// Embedder represents an embedder
type Embedder interface {
	Embed(text string) (Vector, error)
}
