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

type groupingsTopicsNamer struct {
	endpoint string
	client   *http.Client
}

func (namer *groupingsTopicsNamer) Name(
	ctx context.Context,
	posts []posts.Post,
) (string, error) {
	if namer.endpoint == "" {
		return "", ErrInvalidGroupingsTopicsNamerEndpoint
	}

	texts := postTexts(posts)
	if len(texts) == 0 {
		return "", ErrInvalidGroupingsTopicsNamerPosts
	}

	payload, err := json.Marshal(
		nameClusterRequest{
			Posts:        texts,
			SystemPrompt: groupingsTopicsSystemPrompt,
			UserPrompt:   groupingsTopicsUserPrompt(texts),
			Temperature:  0.2,
			MaxTokens:    64,
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
			ErrInvalidGroupingsTopicsNamerResponse,
			response.StatusCode,
		)
	}

	var result nameClusterResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", err
	}

	result.Name = strings.TrimSpace(result.Name)
	if result.Name == "" {
		return "", ErrInvalidGroupingsTopicsNamerResponse
	}

	return result.Name, nil
}

const groupingsTopicsSystemPrompt = `
You are a semantic topic naming service.

Your task:
- Read a group of social media posts.
- Produce a short neutral topic name.
- The name must describe the subject being discussed.
- Do not describe persuasion, coordination, intent, propaganda, or campaign behavior.
- Prefer concrete topics over abstract narratives.

Rules:
- Return ONLY valid JSON.
- No explanations.
- No markdown.
- No reasoning.
- No extra text.

Valid response:

{"name":"Electric Vehicles"}
`

func groupingsTopicsUserPrompt(
	texts []string,
) string {
	sample := strings.Builder{}

	for index, text := range texts {
		if index >= 10 {
			break
		}

		sample.WriteString("- ")
		sample.WriteString(text)
		sample.WriteString("\n")
	}

	return fmt.Sprintf(
		`
Posts:

%s

Return only:

{"name":"topic name"}
`,
		sample.String(),
	)
}
