package posts

import (
	"context"

	"github.com/google/uuid"

	domain_posts "github.com/steve-rodrigue/aabs/services/saas/domain/entities/posts"
)

func NewMockPostsApplication() *MockPostsApplication {
	return &MockPostsApplication{}
}

type MockPostsApplication struct {
	SaveCalls int
	SaveErr   error
	LastPost  domain_posts.Post

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue domain_posts.Post

	FindCalls int
	FindErr   error
	FindValue []domain_posts.Post

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []domain_posts.Post

	FindByCriteriaCalls int
	FindByCriteriaErr   error
	FindByCriteriaValue []domain_posts.Post

	FindByCriteriaAfterCalls int
	FindByCriteriaAfterErr   error
	FindByCriteriaAfterValue []domain_posts.Post

	CountCalls int
	CountErr   error
	CountValue int64

	CountByCriteriaCalls int
	CountByCriteriaErr   error
	CountByCriteriaValue int64

	LastContext  context.Context
	LastID       uuid.UUID
	LastIndex    int
	LastAmount   int
	LastCursor   uuid.UUID
	LastCriteria domain_posts.Criteria
}

func (application *MockPostsApplication) Save(
	ctx context.Context,
	post domain_posts.Post,
) error {
	application.SaveCalls++
	application.LastContext = ctx
	application.LastPost = post

	return application.SaveErr
}

func (application *MockPostsApplication) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_posts.Post, error) {
	application.FindByIDCalls++
	application.LastContext = ctx
	application.LastID = id

	return application.FindByIDValue, application.FindByIDErr
}

func (application *MockPostsApplication) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_posts.Post, error) {
	application.FindCalls++
	application.LastContext = ctx
	application.LastIndex = index
	application.LastAmount = amount

	return application.FindValue, application.FindErr
}

func (application *MockPostsApplication) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_posts.Post, error) {
	application.FindAfterCalls++
	application.LastContext = ctx
	application.LastCursor = cursor
	application.LastAmount = amount

	if application.FindAfterErr != nil {
		return nil, application.FindAfterErr
	}

	if application.FindAfterValue != nil {
		if application.FindAfterCalls == 1 {
			return application.FindAfterValue, nil
		}

		return []domain_posts.Post{}, nil
	}

	return []domain_posts.Post{}, nil
}

func (application *MockPostsApplication) FindByCriteria(
	ctx context.Context,
	criteria domain_posts.Criteria,
	index int,
	amount int,
) ([]domain_posts.Post, error) {
	application.FindByCriteriaCalls++
	application.LastContext = ctx
	application.LastCriteria = criteria
	application.LastIndex = index
	application.LastAmount = amount

	return application.FindByCriteriaValue,
		application.FindByCriteriaErr
}

func (application *MockPostsApplication) FindByCriteriaAfter(
	ctx context.Context,
	criteria domain_posts.Criteria,
	cursor uuid.UUID,
	amount int,
) ([]domain_posts.Post, error) {
	application.FindByCriteriaAfterCalls++
	application.LastContext = ctx
	application.LastCriteria = criteria
	application.LastCursor = cursor
	application.LastAmount = amount

	if application.FindByCriteriaAfterErr != nil {
		return nil, application.FindByCriteriaAfterErr
	}

	if application.FindByCriteriaAfterValue != nil {
		if application.FindByCriteriaAfterCalls == 1 {
			return application.FindByCriteriaAfterValue, nil
		}

		return []domain_posts.Post{}, nil
	}

	return []domain_posts.Post{}, nil
}

func (application *MockPostsApplication) Count(
	ctx context.Context,
) (int64, error) {
	application.CountCalls++
	application.LastContext = ctx

	return application.CountValue, application.CountErr
}

func (application *MockPostsApplication) CountByCriteria(
	ctx context.Context,
	criteria domain_posts.Criteria,
) (int64, error) {
	application.CountByCriteriaCalls++
	application.LastContext = ctx
	application.LastCriteria = criteria

	return application.CountByCriteriaValue,
		application.CountByCriteriaErr
}
