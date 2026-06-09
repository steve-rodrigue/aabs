package llms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
)

type groupingsCampaignsNamer struct {
	endpoint string
	client   *http.Client
}

type nameClusterRequest struct {
	Posts        []string `json:"posts"`
	SystemPrompt string   `json:"system_prompt,omitempty"`
	UserPrompt   string   `json:"user_prompt,omitempty"`
	Temperature  float64  `json:"temperature"`
	MaxTokens    int      `json:"max_tokens"`
}

type nameClusterResponse struct {
	Name string `json:"name"`
	Raw  string `json:"raw"`
}

func (namer *groupingsCampaignsNamer) Name(
	ctx context.Context,
	posts []posts.Post,
) (string, error) {
	if namer.endpoint == "" {
		return "", ErrInvalidGroupingsCampaignsNamerEndpoint
	}

	texts := postTexts(posts)
	if len(texts) == 0 {
		return "", ErrInvalidGroupingsCampaignsNamerPosts
	}

	payload, err := json.Marshal(
		nameClusterRequest{
			Posts:       texts,
			Temperature: 0.2,
			MaxTokens:   64,
		},
	)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		namer.endpoint+"/name-cluster",
		bytes.NewReader(payload),
	)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := namer.client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode < http.StatusOK ||
		response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf(
			"%w: status %d",
			ErrInvalidGroupingsCampaignsNamerResponse,
			response.StatusCode,
		)
	}

	var result nameClusterResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", err
	}

	result.Name = strings.TrimSpace(result.Name)
	if result.Name == "" {
		return "", ErrInvalidGroupingsCampaignsNamerResponse
	}

	return result.Name, nil
}

func postTexts(
	posts []posts.Post,
) []string {
	out := []string{}

	for _, post := range posts {
		if post == nil ||
			post.Content() == nil {
			continue
		}

		text := strings.TrimSpace(post.Content().Text())
		if text == "" {
			continue
		}

		out = append(out, text)
	}

	return out
}
