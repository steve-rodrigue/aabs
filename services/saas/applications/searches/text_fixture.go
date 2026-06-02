package searches

import (
	"github.com/steve-rodrigue/aabs/services/saas/domain/embeddings"
	domain_searches "github.com/steve-rodrigue/aabs/services/saas/domain/searches"
)

func newFixture() *fixture {
	embedder := &embeddings.MockEmbedder{
		Vector: embeddings.Vector{1, 2, 3},
	}

	searchRepository := domain_searches.NewMockSearchRepository()
	searchableRepository := domain_searches.NewMockSearchableRepository()

	return &fixture{
		app: New(
			embedder,
			searchRepository,
			searchableRepository,
		),

		embedder:             embedder,
		searchRepository:     searchRepository,
		searchableRepository: searchableRepository,
	}
}

type fixture struct {
	app Application

	embedder             *embeddings.MockEmbedder
	searchRepository     *domain_searches.MockSearchRepository
	searchableRepository *domain_searches.MockSearchableRepository
}
