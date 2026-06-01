package searches

import "github.com/google/uuid"

type result struct {
	identifier uuid.UUID
	kind       ResultKind
	title      string
	text       string
	score      float64
}

func (result *result) Identifier() uuid.UUID {
	return result.identifier
}

func (result *result) Kind() ResultKind {
	return result.kind
}

func (result *result) Title() string {
	return result.title
}

func (result *result) Text() string {
	return result.text
}

func (result *result) Score() float64 {
	return result.score
}
