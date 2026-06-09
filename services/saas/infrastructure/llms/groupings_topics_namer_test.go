package llms

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
)

func TestNewGroupingsTopicsNamer(t *testing.T) {
	endpoint := os.Getenv("LLM_TEST_URL")
	if endpoint == "" {
		t.Skip("LLM_TEST_URL is not set")
	}

	namer := NewGroupingsTopicsNamer(endpoint)

	if namer == nil {
		t.Fatalf("expected namer")
	}
}

func TestGroupingsTopicsNamerName(t *testing.T) {
	var requestBody nameClusterRequest

	server := httptest.NewServer(
		http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
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
					Name: " Electric Vehicles ",
					Raw:  `{"name":"Electric Vehicles"}`,
				},
			)
		}),
	)

	defer server.Close()

	namer := NewGroupingsTopicsNamer(server.URL)

	result, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{
			domain_posts.NewMockPost("EV subsidies are changing the auto market"),
			domain_posts.NewMockPost("Electric cars are getting cheaper"),
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if result != "Electric Vehicles" {
		t.Fatalf("expected Electric Vehicles, got %s", result)
	}

	if len(requestBody.Posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(requestBody.Posts))
	}

	if requestBody.Posts[0] != "EV subsidies are changing the auto market" {
		t.Fatalf("unexpected first post %s", requestBody.Posts[0])
	}

	if !strings.Contains(requestBody.SystemPrompt, "semantic topic naming service") {
		t.Fatalf("expected topic system prompt")
	}

	if !strings.Contains(requestBody.UserPrompt, "EV subsidies are changing the auto market") {
		t.Fatalf("expected user prompt to contain post text")
	}

	if requestBody.Temperature != 0.2 {
		t.Fatalf("expected temperature 0.2, got %f", requestBody.Temperature)
	}

	if requestBody.MaxTokens != 64 {
		t.Fatalf("expected max tokens 64, got %d", requestBody.MaxTokens)
	}
}

func TestGroupingsTopicsNamerNameReturnsInvalidEndpointError(t *testing.T) {
	namer := NewGroupingsTopicsNamer("")

	_, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{
			domain_posts.NewMockPost("hello"),
		},
	)

	if !errors.Is(err, ErrInvalidGroupingsTopicsNamerEndpoint) {
		t.Fatalf("expected invalid endpoint error, got %v", err)
	}
}

func TestGroupingsTopicsNamerNameReturnsInvalidPostsErrorWhenPostsAreEmpty(t *testing.T) {
	namer := NewGroupingsTopicsNamer("http://localhost")

	_, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{},
	)

	if !errors.Is(err, ErrInvalidGroupingsTopicsNamerPosts) {
		t.Fatalf("expected invalid posts error, got %v", err)
	}
}

func TestGroupingsTopicsNamerNameReturnsInvalidPostsErrorWhenPostsContainNoText(t *testing.T) {
	namer := NewGroupingsTopicsNamer("http://localhost")

	_, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{
			domain_posts.NewMockPost(""),
			domain_posts.NewMockPost("   "),
		},
	)

	if !errors.Is(err, ErrInvalidGroupingsTopicsNamerPosts) {
		t.Fatalf("expected invalid posts error, got %v", err)
	}
}

func TestGroupingsTopicsNamerNameReturnsServerError(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			writer.WriteHeader(http.StatusInternalServerError)
		}),
	)

	defer server.Close()

	namer := NewGroupingsTopicsNamer(server.URL)

	_, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{
			domain_posts.NewMockPost("hello"),
		},
	)

	if !errors.Is(err, ErrInvalidGroupingsTopicsNamerResponse) {
		t.Fatalf("expected invalid response error, got %v", err)
	}
}

func TestGroupingsTopicsNamerNameReturnsDecodeError(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			_, _ = writer.Write([]byte("not-json"))
		}),
	)

	defer server.Close()

	namer := NewGroupingsTopicsNamer(server.URL)

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

func TestGroupingsTopicsNamerNameReturnsInvalidResponseWhenNameIsEmpty(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			_ = json.NewEncoder(writer).Encode(
				nameClusterResponse{
					Name: "   ",
				},
			)
		}),
	)

	defer server.Close()

	namer := NewGroupingsTopicsNamer(server.URL)

	_, err := namer.Name(
		context.Background(),
		[]domain_posts.Post{
			domain_posts.NewMockPost("hello"),
		},
	)

	if !errors.Is(err, ErrInvalidGroupingsTopicsNamerResponse) {
		t.Fatalf("expected invalid response error, got %v", err)
	}
}

func TestGroupingsTopicsUserPromptLimitsToTenPosts(t *testing.T) {
	texts := []string{
		"1", "2", "3", "4", "5",
		"6", "7", "8", "9", "10",
		"11", "12",
	}

	prompt := groupingsTopicsUserPrompt(texts)

	if strings.Contains(prompt, "- 11") {
		t.Fatalf("expected prompt to limit posts to 10")
	}

	if strings.Contains(prompt, "- 12") {
		t.Fatalf("expected prompt to limit posts to 10")
	}
}
