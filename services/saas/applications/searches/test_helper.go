package searches

import (
	"github.com/steve-rodrigue/aabs/services/saas/domain/communities"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/campaigns"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/narratives"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/topics"
	"github.com/steve-rodrigue/aabs/services/saas/domain/posts"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/users"
)

func NewMockSearchApplication() *MockSearchApplication {
	return &MockSearchApplication{}
}

type MockSearchApplication struct {
	IndexPostCalls int
	IndexPostErr   error
	LastPost       posts.Post

	SearchCalls int
	SearchErr   error
	SearchValue []Result
	LastQuery   string
	LastLimit   int

	SearchPostsCalls int
	SearchPostsErr   error
	SearchPostsValue []posts.Post

	SearchCampaignsCalls int
	SearchCampaignsErr   error
	SearchCampaignsValue []campaigns.Campaign

	SearchTopicsCalls int
	SearchTopicsErr   error
	SearchTopicsValue []topics.Topic

	SearchNarrativesCalls int
	SearchNarrativesErr   error
	SearchNarrativesValue []narratives.Narrative

	SearchUsersCalls int
	SearchUsersErr   error
	SearchUsersValue []users.User

	SearchCommunitiesCalls int
	SearchCommunitiesErr   error
	SearchCommunitiesValue []communities.Community

	SearchRelationshipsCalls int
	SearchRelationshipsErr   error
	SearchRelationshipsValue []relationships.Relationship
}

func (application *MockSearchApplication) IndexPost(
	post posts.Post,
) error {
	application.IndexPostCalls++
	application.LastPost = post

	return application.IndexPostErr
}

func (application *MockSearchApplication) Search(
	query string,
	limit int,
) ([]Result, error) {
	application.SearchCalls++
	application.LastQuery = query
	application.LastLimit = limit

	return application.SearchValue, application.SearchErr
}

func (application *MockSearchApplication) SearchPosts(
	query string,
	limit int,
) ([]posts.Post, error) {
	application.SearchPostsCalls++
	application.LastQuery = query
	application.LastLimit = limit

	return application.SearchPostsValue, application.SearchPostsErr
}

func (application *MockSearchApplication) SearchCampaigns(
	query string,
	limit int,
) ([]campaigns.Campaign, error) {
	application.SearchCampaignsCalls++
	application.LastQuery = query
	application.LastLimit = limit

	return application.SearchCampaignsValue, application.SearchCampaignsErr
}

func (application *MockSearchApplication) SearchTopics(
	query string,
	limit int,
) ([]topics.Topic, error) {
	application.SearchTopicsCalls++
	application.LastQuery = query
	application.LastLimit = limit

	return application.SearchTopicsValue, application.SearchTopicsErr
}

func (application *MockSearchApplication) SearchNarratives(
	query string,
	limit int,
) ([]narratives.Narrative, error) {
	application.SearchNarrativesCalls++
	application.LastQuery = query
	application.LastLimit = limit

	return application.SearchNarrativesValue, application.SearchNarrativesErr
}

func (application *MockSearchApplication) SearchUsers(
	query string,
	limit int,
) ([]users.User, error) {
	application.SearchUsersCalls++
	application.LastQuery = query
	application.LastLimit = limit

	return application.SearchUsersValue, application.SearchUsersErr
}

func (application *MockSearchApplication) SearchCommunities(
	query string,
	limit int,
) ([]communities.Community, error) {
	application.SearchCommunitiesCalls++
	application.LastQuery = query
	application.LastLimit = limit

	return application.SearchCommunitiesValue, application.SearchCommunitiesErr
}

func (application *MockSearchApplication) SearchRelationships(
	query string,
	limit int,
) ([]relationships.Relationship, error) {
	application.SearchRelationshipsCalls++
	application.LastQuery = query
	application.LastLimit = limit

	return application.SearchRelationshipsValue, application.SearchRelationshipsErr
}
