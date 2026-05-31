package namers

import "github.com/steve-rodrigue/aabs/services/saas/domain/posts"

// Namer represents a campaign namer
type Namer interface {
	Name(posts []posts.Post) (string, error)
}
