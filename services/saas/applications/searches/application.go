package searches

import (
	"github.com/steve-rodrigue/aabs/services/saas/domain/embeddings"
	domain_searches "github.com/steve-rodrigue/aabs/services/saas/domain/searches"
)

type application struct {
	embedder             embeddings.Embedder
	searchRepository     domain_searches.Repository
	searchableRepository domain_searches.SearchableRepository
}

func createApplication(
	embedder embeddings.Embedder,
	searchRepository domain_searches.Repository,
	searchableRepository domain_searches.SearchableRepository,
) Application {
	return &application{
		embedder:             embedder,
		searchRepository:     searchRepository,
		searchableRepository: searchableRepository,
	}
}

func (app *application) Index(
	searchable domain_searches.Searchable,
) error {
	vector, err := app.embedder.Embed(searchable.SearchText())
	if err != nil {
		return err
	}

	return app.searchRepository.Store(
		searchable.Identifier(),
		searchable.SearchKind(),
		vector,
	)
}

func (app *application) Search(
	query string,
	limit int,
) ([]Result, error) {
	matches, err := app.search(query, limit)
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(matches))

	for _, match := range matches {
		searchable, err := app.searchableRepository.FindByID(
			match.Target(),
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &result{
			identifier: searchable.Identifier(),
			kind:       ResultKind(searchable.SearchKind()),
			title:      searchable.SearchTitle(),
			text:       searchable.SearchText(),
			score:      match.Similarity(),
		})
	}

	return results, nil
}

func (app *application) search(
	query string,
	limit int,
) ([]domain_searches.Match, error) {
	vector, err := app.embedder.Embed(query)
	if err != nil {
		return nil, err
	}

	return app.searchRepository.Search(vector, limit)
}
