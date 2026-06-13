package communities

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications/platforms"
	client_internal "github.com/steve-rodrigue/aabs/services/posts/restapis/clients/internal"
	domain_communities "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/communities"
	domain_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
	domain_users "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/dtos"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/inputs"
)

type application struct {
	baseURL   string
	client    *http.Client
	platforms platforms.Application
}

func (app *application) Save(
	ctx context.Context,
	community domain_communities.Community,
) error {
	moderatorIDs := make(
		[]uuid.UUID,
		0,
		len(community.Moderators()),
	)

	for _, moderator := range community.Moderators() {
		moderatorIDs = append(
			moderatorIDs,
			moderator.Identifier(),
		)
	}

	return client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodPost,
		app.baseURL+"/communities",
		inputs.SaveCommunityRequest{
			Identifier:   community.Identifier(),
			PlatformID:   community.Platform().Identifier(),
			Handle:       community.Handle(),
			Title:        community.Title(),
			Text:         community.Text(),
			CreatedOn:    community.CreatedOn(),
			ModeratorIDs: moderatorIDs,
		},
		nil,
	)
}

func (app *application) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_communities.Community, error) {
	var response dtos.CommunityResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/communities/"+id.String(),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return communityResponseToDomain(response)
}

func (app *application) FindByHandle(
	ctx context.Context,
	platform domain_platforms.Platform,
	handle string,
) (domain_communities.Community, error) {
	var response dtos.CommunityResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/communities/platform/"+platform.Identifier().String()+"/handle/"+url.PathEscape(handle),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return communityResponseToDomain(response)
}

func (app *application) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_communities.Community, error) {
	values := url.Values{}
	values.Set("index", strconv.Itoa(index))
	values.Set("amount", strconv.Itoa(amount))

	var response []dtos.CommunityResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/communities?"+values.Encode(),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return communitiesResponseToDomain(response)
}

func (app *application) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_communities.Community, error) {
	values := url.Values{}
	values.Set("amount", strconv.Itoa(amount))

	if cursor != uuid.Nil {
		values.Set("cursor", cursor.String())
	}

	var response []dtos.CommunityResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/communities?"+values.Encode(),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return communitiesResponseToDomain(response)
}

func (app *application) FindByPlatform(
	ctx context.Context,
	platform domain_platforms.Platform,
) ([]domain_communities.Community, error) {
	var response []dtos.CommunityResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/communities/platform/"+platform.Identifier().String(),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return communitiesResponseToDomain(response)
}

func (app *application) Count(
	ctx context.Context,
) (int64, error) {
	var response struct {
		Count int64 `json:"count"`
	}

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/communities/count",
		nil,
		&response,
	)
	if err != nil {
		return 0, err
	}

	return response.Count, nil
}

func communityResponseToDomain(
	response dtos.CommunityResponse,
) (domain_communities.Community, error) {
	platform, err := platformResponseToDomain(response.Platform)
	if err != nil {
		return nil, err
	}

	moderators, err := usersResponseToDomain(response.Moderators)
	if err != nil {
		return nil, err
	}

	return domain_communities.NewAdapter().ToDomain(
		domain_communities.CommunityInput{
			Identifier: response.Identifier,
			Platform:   platform,
			Handle:     response.Handle,
			Title:      response.Title,
			Text:       response.Text,
			CreatedOn:  response.CreatedOn,
			Moderators: moderators,
		},
	)
}

func communitiesResponseToDomain(
	input []dtos.CommunityResponse,
) ([]domain_communities.Community, error) {
	out := make([]domain_communities.Community, 0, len(input))

	for _, response := range input {
		community, err := communityResponseToDomain(response)
		if err != nil {
			return nil, err
		}

		out = append(out, community)
	}

	return out, nil
}

func usersResponseToDomain(
	input []dtos.UserResponse,
) ([]domain_users.User, error) {
	out := make([]domain_users.User, 0, len(input))

	for _, response := range input {
		user, err := userResponseToDomain(response)
		if err != nil {
			return nil, err
		}

		out = append(out, user)
	}

	return out, nil
}

func userResponseToDomain(
	response dtos.UserResponse,
) (domain_users.User, error) {
	platform, err := platformResponseToDomain(response.Platform)
	if err != nil {
		return nil, err
	}

	return domain_users.NewAdapter().ToDomain(
		domain_users.UserInput{
			Identifier:  response.Identifier,
			Platform:    platform,
			ExternalID:  response.ExternalID,
			Handle:      response.Handle,
			DisplayName: response.DisplayName,
			ProfileURL:  response.ProfileURL,
			CreatedOn:   response.CreatedOn,
		},
	)
}

func platformResponseToDomain(
	response dtos.PlatformResponse,
) (domain_platforms.Platform, error) {
	return domain_platforms.NewAdapter().ToDomain(
		domain_platforms.PlatformInput{
			Identifier: response.Identifier,
			Name:       response.Name,
			Handle:     response.Handle,
			BaseURL:    response.BaseURL,
			CreatedOn:  response.CreatedOn,
		},
	)
}
