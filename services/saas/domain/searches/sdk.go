package searches

import "github.com/google/uuid"

// MatchInput represents a match input
type MatchInput struct {
	Target     uuid.UUID
	Similarity float64
}

// Adapter represents a match adapter
type Adapter interface {
	ToDomain(input MatchInput) (Match, error)
}

// Match represents a match
type Match interface {
	Target() uuid.UUID
	Similarity() float64
}

// Repository represents a search repository
type Repository interface {
	Store(target uuid.UUID, vector []float32) error
	Search(vector []float32, limit int) ([]Match, error)
}
