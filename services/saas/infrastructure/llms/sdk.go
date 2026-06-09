package llms

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/namers"
)

var (
	// campaigns
	ErrInvalidGroupingsCampaignsNamerEndpoint = errors.New("invalid groupings campaigns namer endpoint")
	ErrInvalidGroupingsCampaignsNamerPosts    = errors.New("invalid groupings campaigns namer posts")
	ErrInvalidGroupingsCampaignsNamerResponse = errors.New("invalid groupings campaigns namer response")

	// topics
	ErrInvalidGroupingsTopicsNamerEndpoint = errors.New("invalid groupings topics namer endpoint")
	ErrInvalidGroupingsTopicsNamerPosts    = errors.New("invalid groupings topics namer posts")
	ErrInvalidGroupingsTopicsNamerResponse = errors.New("invalid groupings topics namer response")
)

// NewGroupingsCampaignsNamer creates a new groupings campaigns namer
func NewGroupingsCampaignsNamer(
	endpoint string,
) namers.Namer {
	return &groupingsCampaignsNamer{
		endpoint: strings.TrimRight(endpoint, "/"),
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// NewGroupingsTopicsNamer creates a new groupings topics namer
func NewGroupingsTopicsNamer(
	endpoint string,
) namers.Namer {

	return &groupingsTopicsNamer{
		endpoint: strings.TrimRight(endpoint, "/"),
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}

}
