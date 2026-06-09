package llms

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
)

func llmURL() string {
	url := os.Getenv("LLM_TEST_URL")
	if url == "" {
		return "http://localhost:8100"
	}

	return url
}

func TestNewGroupingsCampaignsNamer(t *testing.T) {
	namer := NewGroupingsCampaignsNamer(llmURL())

	if namer == nil {
		t.Fatalf("expected namer")
	}
}

func TestGroupingsCampaignsNamerName(t *testing.T) {
	var requestBody nameClusterRequest

	server := httptest.NewServer(
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", request.Method)
			}

			if request.URL.Path != "/name-cluster" {
				t.Fatalf("expected /name-cluster, got %s", request.URL.Path)
			}

			if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
				t.Fatal(err)
			}

			writer.Header().Set("Content-Type", "application/json")

			_ = json.NewEncoder(writer).Encode(
				nameClusterResponse{
					Name: " Coordinated Crypto Promotion ",
					Raw:  `{"name":"Coordinated Crypto Promotion"}`,
				},
			)
		}),
	)
	defer server.Close()

	namer := NewGroupingsCampaignsNamer(server.URL)

	result, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{
			domain_posts.NewMockPost("Buy this coin now"),
			domain_posts.NewMockPost("This coin is going to moon"),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != "Coordinated Crypto Promotion" {
		t.Fatalf("expected campaign name, got %q", result)
	}

	if len(requestBody.Posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(requestBody.Posts))
	}

	if requestBody.Posts[0] != "Buy this coin now" {
		t.Fatalf("expected first post text, got %q", requestBody.Posts[0])
	}

	if requestBody.Posts[1] != "This coin is going to moon" {
		t.Fatalf("expected second post text, got %q", requestBody.Posts[1])
	}

	if requestBody.Temperature != 0.2 {
		t.Fatalf("expected temperature 0.2, got %f", requestBody.Temperature)
	}

	if requestBody.MaxTokens != 64 {
		t.Fatalf("expected max tokens 64, got %d", requestBody.MaxTokens)
	}
}

func TestGroupingsCampaignsNamerNameTrimsEndpoint(t *testing.T) {
	var called bool

	server := httptest.NewServer(
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			called = true

			if request.URL.Path != "/name-cluster" {
				t.Fatalf("expected /name-cluster, got %s", request.URL.Path)
			}

			writer.Header().Set("Content-Type", "application/json")

			_ = json.NewEncoder(writer).Encode(
				nameClusterResponse{
					Name: "Crypto Promotion",
					Raw:  `{"name":"Crypto Promotion"}`,
				},
			)
		}),
	)
	defer server.Close()

	namer := NewGroupingsCampaignsNamer(server.URL + "/")

	_, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{
			domain_posts.NewMockPost("Buy crypto"),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Fatalf("expected server to be called")
	}
}

func TestGroupingsCampaignsNamerNameReturnsInvalidEndpointError(t *testing.T) {
	namer := NewGroupingsCampaignsNamer("")

	_, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{
			domain_posts.NewMockPost("hello"),
		},
	)

	if !errors.Is(err, ErrInvalidGroupingsCampaignsNamerEndpoint) {
		t.Fatalf("expected invalid endpoint error, got %v", err)
	}
}

func TestGroupingsCampaignsNamerNameReturnsInvalidPostsErrorWhenPostsAreEmpty(t *testing.T) {
	namer := NewGroupingsCampaignsNamer(llmURL())

	_, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{},
	)

	if !errors.Is(err, ErrInvalidGroupingsCampaignsNamerPosts) {
		t.Fatalf("expected invalid posts error, got %v", err)
	}
}

func TestGroupingsCampaignsNamerNameReturnsInvalidPostsErrorWhenPostTextIsEmpty(t *testing.T) {
	namer := NewGroupingsCampaignsNamer(llmURL())

	_, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{
			domain_posts.NewMockPost(" "),
		},
	)

	if !errors.Is(err, ErrInvalidGroupingsCampaignsNamerPosts) {
		t.Fatalf("expected invalid posts error, got %v", err)
	}
}

func TestGroupingsCampaignsNamerNameReturnsServerError(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer server.Close()

	namer := NewGroupingsCampaignsNamer(server.URL)

	_, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{
			domain_posts.NewMockPost("hello"),
		},
	)

	if !errors.Is(err, ErrInvalidGroupingsCampaignsNamerResponse) {
		t.Fatalf("expected invalid response error, got %v", err)
	}
}

func TestGroupingsCampaignsNamerNameReturnsDecodeError(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte("invalid-json"))
		}),
	)
	defer server.Close()

	namer := NewGroupingsCampaignsNamer(server.URL)

	_, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{
			domain_posts.NewMockPost("hello"),
		},
	)

	if err == nil {
		t.Fatalf("expected decode error")
	}
}

func TestGroupingsCampaignsNamerNameReturnsInvalidResponseWhenNameIsEmpty(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "application/json")

			_ = json.NewEncoder(writer).Encode(
				nameClusterResponse{
					Name: " ",
					Raw:  `{"name":""}`,
				},
			)
		}),
	)
	defer server.Close()

	namer := NewGroupingsCampaignsNamer(server.URL)

	_, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{
			domain_posts.NewMockPost("hello"),
		},
	)

	if !errors.Is(err, ErrInvalidGroupingsCampaignsNamerResponse) {
		t.Fatalf("expected invalid response error, got %v", err)
	}
}

func TestGroupingsCampaignsNamerNameReturnsContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	namer := NewGroupingsCampaignsNamer(llmURL())

	_, err := namer.Name(
		ctx,
		[]domain_posts.Post{
			domain_posts.NewMockPost("hello"),
		},
	)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled error, got %v", err)
	}
}

func TestPostTextsSkipsNilAndEmptyPosts(t *testing.T) {
	result := postTexts(
		[]domain_posts.Post{
			nil,
			domain_posts.NewMockPost(" "),
			domain_posts.NewMockPost("hello"),
		},
	)

	if len(result) != 1 {
		t.Fatalf("expected 1 post text, got %d", len(result))
	}

	if result[0] != "hello" {
		t.Fatalf("expected hello, got %q", result[0])
	}
}
