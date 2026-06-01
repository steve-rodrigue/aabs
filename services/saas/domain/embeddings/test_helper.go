package embeddings

type MockEmbedder struct {
	EmbedCalls int
	LastText   string
	Vector     Vector
	EmbedErr   error
}

func (embedder *MockEmbedder) Embed(
	text string,
) (Vector, error) {
	embedder.EmbedCalls++
	embedder.LastText = text

	return embedder.Vector, embedder.EmbedErr
}
