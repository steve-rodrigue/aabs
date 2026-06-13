package servers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/dtos"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/inputs"
)

func (api *api) createUser(writer http.ResponseWriter, request *http.Request) {
	var input inputs.SaveUserRequest

	if err := decodeJSON(request, &input); err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	platform, err := api.application.Platforms().FindByID(request.Context(), input.PlatformID)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if platform == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	user, err := users.NewAdapter().ToDomain(
		users.UserInput{
			Identifier:  input.Identifier,
			Platform:    platform,
			ExternalID:  input.ExternalID,
			Handle:      input.Handle,
			DisplayName: input.DisplayName,
			ProfileURL:  input.ProfileURL,
			CreatedOn:   input.CreatedOn,
		},
	)
	if err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	if err := api.application.Users().Save(request.Context(), user); err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	respondJSON(writer, http.StatusCreated, dtos.UserDTO(user))
}

func (api *api) findUserByID(writer http.ResponseWriter, request *http.Request) {
	id, err := routeUUID(request, "id")
	if err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	user, err := api.application.Users().FindByID(request.Context(), id)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	respondJSON(writer, http.StatusOK, dtos.UserDTO(user))
}

func (api *api) findUserByExternalID(writer http.ResponseWriter, request *http.Request) {
	platformID, err := routeUUID(request, "platform_id")
	if err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	platform, err := api.application.Platforms().FindByID(request.Context(), platformID)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if platform == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	user, err := api.application.Users().FindByExternalID(
		request.Context(),
		platform,
		routeValue(request, "external_id"),
	)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	respondJSON(writer, http.StatusOK, dtos.UserDTO(user))
}

func (api *api) findUserByHandle(writer http.ResponseWriter, request *http.Request) {
	platformID, err := routeUUID(request, "platform_id")
	if err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	platform, err := api.application.Platforms().FindByID(request.Context(), platformID)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if platform == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	user, err := api.application.Users().FindByHandle(
		request.Context(),
		platform,
		routeValue(request, "handle"),
	)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	respondJSON(writer, http.StatusOK, dtos.UserDTO(user))
}

func (api *api) findUsers(writer http.ResponseWriter, request *http.Request) {
	index := queryInt(request, "index", 0)
	amount := queryInt(request, "amount", 25)
	cursor := queryUUID(request, "cursor")

	var (
		result []users.User
		err    error
	)

	if cursor != uuid.Nil {
		result, err = api.application.Users().FindAfter(request.Context(), cursor, amount)
	} else {
		result, err = api.application.Users().Find(request.Context(), index, amount)
	}

	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	out := []dtos.UserResponse{}
	for _, user := range result {
		out = append(out, dtos.UserDTO(user))
	}

	respondJSON(writer, http.StatusOK, out)
}

func (api *api) countUsers(writer http.ResponseWriter, request *http.Request) {
	count, err := api.application.Users().Count(request.Context())
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	respondJSON(writer, http.StatusOK, map[string]int64{"count": count})
}
