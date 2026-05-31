package pipelines

import domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/posts"

// Pipeline represents the full post ingestion and analysis pipeline
type Pipeline interface {
	Process(post domain_posts.Post) error
}
