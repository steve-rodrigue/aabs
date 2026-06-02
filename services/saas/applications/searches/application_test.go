package searches

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	domain_searches "github.com/steve-rodrigue/aabs/services/saas/domain/searches"
)

var errTest = errors.New("test error")

func TestIndex(t *testing.T) {
	fixture := newFixture()

	searchable := domain_searches.NewMockSearchable(
		uuid.New(),
		domain_searches.PostKind,
		"",
		"hello world",
	)

	err := fixture.app.Index(searchable)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.embedder.EmbedCalls != 1 {
		t.Fatalf("expected 1 embed call")
	}

	if fixture.searchRepository.StoreCalls != 1 {
		t.Fatalf("expected 1 store call")
	}

	if fixture.searchRepository.LastStoredTarget != searchable.Identifier() {
		t.Fatalf("expected stored target")
	}

	if fixture.searchRepository.LastStoredKind != domain_searches.PostKind {
		t.Fatalf("expected stored kind")
	}

	if len(fixture.searchRepository.LastStoredVector) != 3 {
		t.Fatalf("expected stored vector")
	}
}

func TestIndexReturnsEmbedderError(t *testing.T) {
	fixture := newFixture()
	fixture.embedder.EmbedErr = errTest

	searchable := domain_searches.NewMockSearchable(
		uuid.New(),
		domain_searches.PostKind,
		"",
		"hello world",
	)

	err := fixture.app.Index(searchable)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected embedder error, got %v", err)
	}

	if fixture.searchRepository.StoreCalls != 0 {
		t.Fatalf("expected store not to be called")
	}
}

func TestIndexReturnsStoreError(t *testing.T) {
	fixture := newFixture()
	fixture.searchRepository.StoreErr = errTest

	searchable := domain_searches.NewMockSearchable(
		uuid.New(),
		domain_searches.PostKind,
		"",
		"hello world",
	)

	err := fixture.app.Index(searchable)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected store error, got %v", err)
	}
}

func TestSearch(t *testing.T) {
	fixture := newFixture()

	id := uuid.New()

	searchable := domain_searches.NewMockSearchable(
		id,
		domain_searches.PostKind,
		"",
		"reply text",
	)

	fixture.searchableRepository.Items[id] = searchable

	fixture.searchRepository.Matches = []domain_searches.Match{
		domain_searches.NewMockMatch(
			id,
			domain_searches.PostKind,
			0.95,
		),
	}

	result, err := fixture.app.Search("hello", 10)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.embedder.EmbedCalls != 1 {
		t.Fatalf("expected 1 embed call")
	}

	if fixture.searchRepository.SearchCalls != 1 {
		t.Fatalf("expected 1 search call")
	}

	if fixture.searchableRepository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 searchable lookup")
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result")
	}

	if result[0].Identifier() != id {
		t.Fatalf("expected result identifier")
	}

	if result[0].Kind() != PostKind {
		t.Fatalf("expected post result kind")
	}

	if result[0].HasTitle() {
		t.Fatalf("expected result without title")
	}

	if result[0].Title() != "" {
		t.Fatalf("expected empty title")
	}

	if result[0].Text() != "reply text" {
		t.Fatalf("expected result text")
	}

	if result[0].Score() != 0.95 {
		t.Fatalf("expected result score")
	}
}

func TestSearchWithTitle(t *testing.T) {
	fixture := newFixture()

	id := uuid.New()

	searchable := domain_searches.NewMockSearchable(
		id,
		domain_searches.PostKind,
		"Thread title",
		"thread text",
	)

	fixture.searchableRepository.Items[id] = searchable

	fixture.searchRepository.Matches = []domain_searches.Match{
		domain_searches.NewMockMatch(
			id,
			domain_searches.PostKind,
			0.90,
		),
	}

	result, err := fixture.app.Search("thread", 10)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result")
	}

	if !result[0].HasTitle() {
		t.Fatalf("expected result with title")
	}

	if result[0].Title() != "Thread title" {
		t.Fatalf("expected thread title")
	}

	if result[0].Text() != "thread text" {
		t.Fatalf("expected thread text")
	}
}

func TestSearchReturnsEmbedderError(t *testing.T) {
	fixture := newFixture()
	fixture.embedder.EmbedErr = errTest

	_, err := fixture.app.Search("hello", 10)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected embedder error, got %v", err)
	}

	if fixture.searchRepository.SearchCalls != 0 {
		t.Fatalf("expected search not to be called")
	}
}

func TestSearchReturnsSearchRepositoryError(t *testing.T) {
	fixture := newFixture()
	fixture.searchRepository.SearchErr = errTest

	_, err := fixture.app.Search("hello", 10)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected search error, got %v", err)
	}
}

func TestSearchReturnsSearchableRepositoryError(t *testing.T) {
	fixture := newFixture()

	id := uuid.New()

	fixture.searchRepository.Matches = []domain_searches.Match{
		domain_searches.NewMockMatch(
			id,
			domain_searches.PostKind,
			0.95,
		),
	}

	fixture.searchableRepository.FindByIDErr = errTest

	_, err := fixture.app.Search("hello", 10)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected searchable repository error, got %v", err)
	}
}

func TestSearchEmpty(t *testing.T) {
	fixture := newFixture()

	result, err := fixture.app.Search("nothing", 10)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 0 {
		t.Fatalf("expected no results")
	}
}
