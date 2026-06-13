package users

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications/platforms"
	client_internal "github.com/steve-rodrigue/aabs/services/posts/restapis/clients/internal"
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
	user domain_users.User,
) error {
	return client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodPost,
		app.baseURL+"/users",
		inputs.SaveUserRequest{
			Identifier:  user.Identifier(),
			PlatformID:  user.Platform().Identifier(),
			ExternalID:  user.ExternalID(),
			Handle:      user.Handle(),
			DisplayName: user.DisplayName(),
			ProfileURL:  user.ProfileURL(),
			CreatedOn:   user.CreatedOn(),
		},
		nil,
	)
}

func (app *application) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_users.User, error) {
	var response dtos.UserResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/users/"+id.String(),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return userResponseToDomain(response)
}

func (app *application) FindByExternalID(
	ctx context.Context,
	platform domain_platforms.Platform,
	externalID string,
) (domain_users.User, error) {
	var response dtos.UserResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/users/platform/"+platform.Identifier().String()+"/external/"+url.PathEscape(externalID),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return userResponseToDomain(response)
}

func (app *application) FindByHandle(
	ctx context.Context,
	platform domain_platforms.Platform,
	handle string,
) (domain_users.User, error) {
	var response dtos.UserResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/users/platform/"+platform.Identifier().String()+"/handle/"+url.PathEscape(handle),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return userResponseToDomain(response)
}

func (app *application) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_users.User, error) {
	values := url.Values{}
	values.Set("index", strconv.Itoa(index))
	values.Set("amount", strconv.Itoa(amount))

	var response []dtos.UserResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/users?"+values.Encode(),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return usersResponseToDomain(response)
}

func (app *application) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_users.User, error) {
	values := url.Values{}
	values.Set("amount", strconv.Itoa(amount))

	if cursor != uuid.Nil {
		values.Set("cursor", cursor.String())
	}

	var response []dtos.UserResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/users?"+values.Encode(),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return usersResponseToDomain(response)
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
		app.baseURL+"/users/count",
		nil,
		&response,
	)
	if err != nil {
		return 0, err
	}

	return response.Count, nil
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
