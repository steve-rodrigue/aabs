package topics

import (
	"strings"

	"github.com/google/uuid"
)

type adapter struct{}

func (adapter *adapter) ToDomain(
	input TopicInput,
) (Topic, error) {
	if input.Identifier == uuid.Nil {
		return nil, ErrInvalidTopicIdentifier
	}

	if input.Cluster == nil {
		return nil, ErrInvalidTopicCluster
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidTopicName
	}

	if input.CreatedOn.IsZero() {
		return nil, ErrInvalidTopicCreatedOn
	}

	return &topic{
		identifier:  input.Identifier,
		cluster:     input.Cluster,
		name:        name,
		description: strings.TrimSpace(input.Description),
		parent:      input.Parent,
		createdOn:   input.CreatedOn.UTC(),
	}, nil
}
