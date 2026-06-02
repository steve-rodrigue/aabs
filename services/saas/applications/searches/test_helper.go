package searches

import (
	domain_searches "github.com/steve-rodrigue/aabs/services/saas/domain/searches"
)

func NewMockSearchApplication() *MockSearchApplication {
	return &MockSearchApplication{}
}

type MockSearchApplication struct {
	IndexCalls int
	IndexErr   error

	LastIndexed domain_searches.Searchable

	SearchCalls int
	SearchErr   error
	SearchValue []Result

	LastQuery string
	LastLimit int
}

func (application *MockSearchApplication) Index(
	searchable domain_searches.Searchable,
) error {
	application.IndexCalls++
	application.LastIndexed = searchable

	return application.IndexErr
}

func (application *MockSearchApplication) Search(
	query string,
	limit int,
) ([]Result, error) {
	application.SearchCalls++
	application.LastQuery = query
	application.LastLimit = limit

	return application.SearchValue, application.SearchErr
}
