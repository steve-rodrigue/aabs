package namers

import (
	"context"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
)

// Namer represents a grouping namer
type Namer interface {
	Name(ctx context.Context, posts []posts.Post) (string, error)
}
