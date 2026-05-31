package factors

type Type string

const (
	SemanticRepetitionType        Type = "semantic_repetition"
	CampaignParticipationType     Type = "campaign_participation"
	UserConcentrationType         Type = "user_concentration"
	AccountAgeType                Type = "account_age"
	PostingVelocityType           Type = "posting_velocity"
	RelationshipRiskType          Type = "relationship_risk"
	CommunitySpreadType           Type = "community_spread"
	LLMQualitySignalType          Type = "llm_quality_signal"
	TopicParticipationType        Type = "topic_participation"
	NarrativeParticipationType    Type = "narrative_participation"
	CommunityConcentrationType    Type = "community_concentration"
	CrossCommunitySpreadType      Type = "cross_community_spread"
	ContentDiversityType          Type = "content_diversity"
	TemporalSynchronizationType   Type = "temporal_synchronization"
	NetworkCoordinationType       Type = "network_coordination"
	SimilarityToKnownCampaignType Type = "similarity_to_known_campaign"
)

// Factor represents a trust score factor
type Factor interface {
	Name() Type
	Value() float64
	Weight() float64
	Reason() string
}
