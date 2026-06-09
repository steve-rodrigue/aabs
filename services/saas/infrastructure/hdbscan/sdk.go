package hdbscan

import (
	"errors"
	"net/http"
	"strings"
	"time"

	domain_campaigns "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	domain_clusters "github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

var (
	ErrInvalidGroupingsCampaignDetectorEndpoint   = errors.New("invalid groupings campaigns detector endpoint")
	ErrInvalidGroupingsCampaignDetectorCandidates = errors.New("invalid groupings campaigns detector candidates")
	ErrInvalidGroupingsCampaignDetectorComparable = errors.New("invalid groupings campaigns detector comparable")
	ErrInvalidGroupingsCampaignDetectorVector     = errors.New("invalid groupings campaigns detector vector")
	ErrInvalidGroupingsCampaignDetectorKind       = errors.New("invalid groupings campaigns detector kind")
	ErrInvalidGroupingsCampaignDetectorResponse   = errors.New("invalid groupings campaigns detector response")
)

// NewGroupingsCampaignDetector creates a new grouping campaign detector
func NewGroupingsCampaignDetector(
	endpoint string,
	campaignAdapter domain_campaigns.Adapter,
	clusterAdapter domain_clusters.Adapter,
	comparables clusterables.ComparableRepository,
	minClusterSize int,
	minSamples *int,
) domain_campaigns.Detector {
	return &groupingsCampaignDetector{
		endpoint:        strings.TrimRight(endpoint, "/"),
		client:          &http.Client{Timeout: 120 * time.Second},
		campaignAdapter: campaignAdapter,
		clusterAdapter:  clusterAdapter,
		comparables:     comparables,
		minClusterSize:  minClusterSize,
		minSamples:      minSamples,
	}
}
