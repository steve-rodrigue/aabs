package servers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/dtos"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/inputs"
)

func (api *api) createPlatform(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var input inputs.SavePlatformRequest

	if err := decodeJSON(request, &input); err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	platform, err := platforms.NewAdapter().ToDomain(
		platforms.PlatformInput{
			Identifier: input.Identifier,
			Name:       input.Name,
			Handle:     input.Handle,
			BaseURL:    input.BaseURL,
			CreatedOn:  input.CreatedOn,
		},
	)
	if err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	if err := api.application.Platforms().Save(request.Context(), platform); err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	respondJSON(writer, http.StatusCreated, dtos.PlatformDTO(platform))
}

func (api *api) findPlatformByID(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id, err := routeUUID(request, "id")
	if err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	platform, err := api.application.Platforms().FindByID(request.Context(), id)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if platform == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	respondJSON(writer, http.StatusOK, dtos.PlatformDTO(platform))
}

func (api *api) findPlatformByHandle(
	writer http.ResponseWriter,
	request *http.Request,
) {
	handle := routeValue(request, "handle")

	platform, err := api.application.Platforms().FindByHandle(request.Context(), handle)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if platform == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	respondJSON(writer, http.StatusOK, dtos.PlatformDTO(platform))
}

func (api *api) findPlatforms(
	writer http.ResponseWriter,
	request *http.Request,
) {
	index := queryInt(request, "index", 0)
	amount := queryInt(request, "amount", 25)
	cursor := queryUUID(request, "cursor")

	var (
		result []platforms.Platform
		err    error
	)

	if cursor != uuid.Nil {
		result, err = api.application.Platforms().FindAfter(request.Context(), cursor, amount)
	} else {
		result, err = api.application.Platforms().Find(request.Context(), index, amount)
	}

	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	out := []dtos.PlatformResponse{}
	for _, platform := range result {
		out = append(out, dtos.PlatformDTO(platform))
	}

	respondJSON(writer, http.StatusOK, out)
}

func (api *api) countPlatforms(
	writer http.ResponseWriter,
	request *http.Request,
) {
	count, err := api.application.Platforms().Count(request.Context())
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	respondJSON(writer, http.StatusOK, map[string]int64{"count": count})
}
