package platforms

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"

	client_internal "github.com/steve-rodrigue/aabs/services/posts/restapis/clients/internal"
	domain_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/dtos"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/inputs"
)

type application struct {
	baseURL string
	client  *http.Client
}

func (app *application) Save(
	ctx context.Context,
	platform domain_platforms.Platform,
) error {
	return client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodPost,
		app.baseURL+"/platforms",
		inputs.SavePlatformRequest{
			Identifier: platform.Identifier(),
			Name:       platform.Name(),
			Handle:     platform.Handle(),
			BaseURL:    platform.BaseURL(),
			CreatedOn:  platform.CreatedOn(),
		},
		nil,
	)
}

func (app *application) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_platforms.Platform, error) {
	var response dtos.PlatformResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/platforms/"+id.String(),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return platformResponseToDomain(response)
}

func (app *application) FindByHandle(
	ctx context.Context,
	handle string,
) (domain_platforms.Platform, error) {
	var response dtos.PlatformResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/platforms/handle/"+url.PathEscape(handle),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return platformResponseToDomain(response)
}

func (app *application) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_platforms.Platform, error) {
	values := url.Values{}
	values.Set("index", strconv.Itoa(index))
	values.Set("amount", strconv.Itoa(amount))

	var response []dtos.PlatformResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/platforms?"+values.Encode(),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return platformsResponseToDomain(response)
}

func (app *application) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_platforms.Platform, error) {
	values := url.Values{}
	values.Set("amount", strconv.Itoa(amount))

	if cursor != uuid.Nil {
		values.Set("cursor", cursor.String())
	}

	var response []dtos.PlatformResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/platforms?"+values.Encode(),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return platformsResponseToDomain(response)
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
		app.baseURL+"/platforms/count",
		nil,
		&response,
	)
	if err != nil {
		return 0, err
	}

	return response.Count, nil
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

func platformsResponseToDomain(
	input []dtos.PlatformResponse,
) ([]domain_platforms.Platform, error) {
	out := make([]domain_platforms.Platform, 0, len(input))

	for _, response := range input {
		platform, err := platformResponseToDomain(response)
		if err != nil {
			return nil, err
		}

		out = append(out, platform)
	}

	return out, nil
}
