package embeddings

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func embeddingsURL() string {
	url := os.Getenv("EMBEDDINGS_TEST_URL")

	if url == "" {
		return "http://localhost:8080"
	}

	return url
}

func TestNewEmbedder(t *testing.T) {
	embedder := NewEmbedder(
		embeddingsURL(),
	)

	if embedder == nil {
		t.Fatal("expected embedder")
	}
}

func TestEmbed(t *testing.T) {
	var requestBody embedRequest

	server := httptest.NewServer(
		http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			if request.Method != http.MethodPost {
				t.Fatalf(
					"expected POST, got %s",
					request.Method,
				)
			}

			if request.URL.Path != "/embed" {
				t.Fatalf(
					"expected /embed, got %s",
					request.URL.Path,
				)
			}

			if err := json.NewDecoder(
				request.Body,
			).Decode(&requestBody); err != nil {
				t.Fatal(err)
			}

			writer.Header().Set(
				"Content-Type",
				"application/json",
			)

			_ = json.NewEncoder(writer).Encode(
				embedResponse{
					Dimensions: 3,
					Embedding: []float32{
						0.1,
						0.2,
						0.3,
					},
				},
			)
		}),
	)

	defer server.Close()

	embedder := NewEmbedder(
		server.URL,
	)

	vector, err := embedder.Embed(
		context.Background(),
		" hello world ",
	)

	if err != nil {
		t.Fatal(err)
	}

	if requestBody.Text != "hello world" {
		t.Fatalf(
			"expected trimmed text, got %q",
			requestBody.Text,
		)
	}

	if len(vector) != 3 {
		t.Fatalf(
			"expected vector length 3, got %d",
			len(vector),
		)
	}

	if vector[0] != 0.1 ||
		vector[1] != 0.2 ||
		vector[2] != 0.3 {
		t.Fatalf(
			"unexpected vector %+v",
			vector,
		)
	}
}

func TestEmbedReturnsInvalidTextError(
	t *testing.T,
) {
	embedder := NewEmbedder(
		embeddingsURL(),
	)

	_, err := embedder.Embed(
		context.Background(),
		" ",
	)

	if !errors.Is(
		err,
		ErrInvalidEmbeddingText,
	) {
		t.Fatalf(
			"expected invalid text error, got %v",
			err,
		)
	}
}

func TestEmbedReturnsInvalidEndpointError(
	t *testing.T,
) {
	embedder := NewEmbedder("")

	_, err := embedder.Embed(
		context.Background(),
		"hello",
	)

	if !errors.Is(
		err,
		ErrInvalidEmbeddingEndpoint,
	) {
		t.Fatalf(
			"expected invalid endpoint error, got %v",
			err,
		)
	}
}

func TestEmbedReturnsServerError(
	t *testing.T,
) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			writer.WriteHeader(
				http.StatusInternalServerError,
			)
		}),
	)

	defer server.Close()

	embedder := NewEmbedder(
		server.URL,
	)

	_, err := embedder.Embed(
		context.Background(),
		"hello",
	)

	if !errors.Is(
		err,
		ErrInvalidEmbeddingResponse,
	) {
		t.Fatalf(
			"expected invalid response error, got %v",
			err,
		)
	}
}

func TestEmbedReturnsInvalidJSON(
	t *testing.T,
) {
	server := httptest.NewServer(
		http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			_, _ = writer.Write(
				[]byte("invalid-json"),
			)
		}),
	)

	defer server.Close()

	embedder := NewEmbedder(
		server.URL,
	)

	_, err := embedder.Embed(
		context.Background(),
		"hello",
	)

	if err == nil {
		t.Fatal(
			"expected decode error",
		)
	}
}

func TestEmbedReturnsContextError(
	t *testing.T,
) {
	ctx, cancel := context.WithCancel(
		context.Background(),
	)
	cancel()

	embedder := NewEmbedder(
		embeddingsURL(),
	)

	_, err := embedder.Embed(
		ctx,
		"hello",
	)

	if !errors.Is(
		err,
		context.Canceled,
	) {
		t.Fatalf(
			"expected context canceled error, got %v",
			err,
		)
	}
}

func TestEmbedReturnsInvalidResponseWhenDimensionsAreZero(
	t *testing.T,
) {
	assertInvalidResponse(
		t,
		embedResponse{
			Dimensions: 0,
			Embedding: []float32{
				0.1,
			},
		},
	)
}

func TestEmbedReturnsInvalidResponseWhenEmbeddingIsEmpty(
	t *testing.T,
) {
	assertInvalidResponse(
		t,
		embedResponse{
			Dimensions: 1,
			Embedding:  []float32{},
		},
	)
}

func TestEmbedReturnsInvalidResponseWhenDimensionsMismatch(
	t *testing.T,
) {
	assertInvalidResponse(
		t,
		embedResponse{
			Dimensions: 3,
			Embedding: []float32{
				0.1,
				0.2,
			},
		},
	)
}

func TestEmbedReturnsConnectionError(
	t *testing.T,
) {
	embedder := NewEmbedder(
		"http://127.0.0.1:1",
	)

	_, err := embedder.Embed(
		context.Background(),
		"hello",
	)

	if err == nil {
		t.Fatal(
			"expected connection error",
		)
	}
}

func assertInvalidResponse(
	t *testing.T,
	response embedResponse,
) {
	t.Helper()

	server := httptest.NewServer(
		http.HandlerFunc(func(
			writer http.ResponseWriter,
			request *http.Request,
		) {
			writer.Header().Set(
				"Content-Type",
				"application/json",
			)

			_ = json.NewEncoder(
				writer,
			).Encode(response)
		}),
	)

	defer server.Close()

	embedder := NewEmbedder(
		server.URL,
	)

	_, err := embedder.Embed(
		context.Background(),
		"hello",
	)

	if !errors.Is(
		err,
		ErrInvalidEmbeddingResponse,
	) {
		t.Fatalf(
			"expected invalid response error, got %v",
			err,
		)
	}
}
