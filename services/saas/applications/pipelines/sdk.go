package pipelines

import domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"

// Application represents the post pipeline application
type Application interface {
	ProcessPost(post domain_posts.Post) error
	ProcessPosts(posts []domain_posts.Post) error
	Rebuild() error
}
