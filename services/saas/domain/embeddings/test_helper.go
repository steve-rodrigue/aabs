package embeddings

import "context"

func NewMockEmbedder() *MockEmbedder {
	return &MockEmbedder{}
}

type MockEmbedder struct {
	EmbedCalls int
	EmbedErr   error

	LastContext context.Context
	LastText    string

	Vector Vector
}

func (embedder *MockEmbedder) Embed(
	ctx context.Context,
	text string,
) (Vector, error) {
	embedder.EmbedCalls++

	embedder.LastContext = ctx
	embedder.LastText = text

	return embedder.Vector, embedder.EmbedErr
}
