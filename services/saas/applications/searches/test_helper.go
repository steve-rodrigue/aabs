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

type MockSearchApplication struct {
	SearchPostsCalls int
	SearchPostsErr   error
	LastQuery        string
	LastLimit        int
}

func (application *MockSearchApplication) Search(query string, limit int) ([]Result, error) {
	return nil, nil
}

func (application *MockSearchApplication) SearchPosts(query string, limit int) ([]posts.Post, error) {
	application.SearchPostsCalls++
	application.LastQuery = query
	application.LastLimit = limit

	return nil, application.SearchPostsErr
}

func (application *MockSearchApplication) SearchCampaigns(query string, limit int) ([]campaigns.Campaign, error) {
	return nil, nil
}

func (application *MockSearchApplication) SearchTopics(query string, limit int) ([]topics.Topic, error) {
	return nil, nil
}

func (application *MockSearchApplication) SearchNarratives(query string, limit int) ([]narratives.Narrative, error) {
	return nil, nil
}

func (application *MockSearchApplication) SearchUsers(query string, limit int) ([]users.User, error) {
	return nil, nil
}

func (application *MockSearchApplication) SearchCommunities(query string, limit int) ([]communities.Community, error) {
	return nil, nil
}

func (application *MockSearchApplication) SearchRelationships(query string, limit int) ([]relationships.Relationship, error) {
	return nil, nil
}
