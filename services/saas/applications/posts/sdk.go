package posts

import domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"

// Processor represents a post processing application service
type Processor interface {
	Process(post domain_posts.Post) error
}
