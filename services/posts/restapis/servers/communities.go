package servers

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/communities"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/dtos"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/inputs"
)

func (api *api) createCommunity(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var input inputs.SaveCommunityRequest

	if err := decodeJSON(request, &input); err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	platform, err := api.application.Platforms().FindByID(
		request.Context(),
		input.PlatformID,
	)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if platform == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	moderators := []users.User{}

	for _, moderatorID := range input.ModeratorIDs {
		moderator, err := api.application.Users().FindByID(
			request.Context(),
			moderatorID,
		)
		if err != nil {
			respondError(writer, http.StatusInternalServerError, err)
			return
		}

		if moderator == nil {
			respondJSON(writer, http.StatusNotFound, nil)
			return
		}

		moderators = append(moderators, moderator)
	}

	community, err := communities.NewAdapter().ToDomain(
		communities.CommunityInput{
			Identifier: input.Identifier,
			Platform:   platform,
			Handle:     input.Handle,
			Title:      input.Title,
			Text:       input.Text,
			CreatedOn:  input.CreatedOn,
			Moderators: moderators,
		},
	)
	if err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	if err := api.application.Communities().Save(
		request.Context(),
		community,
	); err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	respondJSON(writer, http.StatusCreated, dtos.CommunityDTO(community))
}

func (api *api) findCommunityByID(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id, err := routeUUID(request, "id")
	if err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	community, err := api.application.Communities().FindByID(
		request.Context(),
		id,
	)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if community == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	respondJSON(writer, http.StatusOK, dtos.CommunityDTO(community))
}

func (api *api) findCommunityByHandle(
	writer http.ResponseWriter,
	request *http.Request,
) {
	platformID, err := routeUUID(request, "platform_id")
	if err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	platform, err := api.application.Platforms().FindByID(
		request.Context(),
		platformID,
	)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if platform == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	community, err := api.application.Communities().FindByHandle(
		request.Context(),
		platform,
		routeValue(request, "handle"),
	)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if community == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	respondJSON(writer, http.StatusOK, dtos.CommunityDTO(community))
}

func (api *api) findCommunities(
	writer http.ResponseWriter,
	request *http.Request,
) {
	index := queryInt(request, "index", 0)
	amount := queryInt(request, "amount", 25)
	cursor := queryUUID(request, "cursor")

	var (
		result []communities.Community
		err    error
	)

	if cursor != uuid.Nil {
		result, err = api.application.Communities().FindAfter(
			request.Context(),
			cursor,
			amount,
		)
	} else {
		result, err = api.application.Communities().Find(
			request.Context(),
			index,
			amount,
		)
	}

	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	out := []dtos.CommunityResponse{}
	for _, community := range result {
		out = append(out, dtos.CommunityDTO(community))
	}

	respondJSON(writer, http.StatusOK, out)
}

func (api *api) findCommunitiesByPlatform(
	writer http.ResponseWriter,
	request *http.Request,
) {
	platformID, err := routeUUID(request, "platform_id")
	if err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	platform, err := api.application.Platforms().FindByID(
		request.Context(),
		platformID,
	)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if platform == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	result, err := api.application.Communities().FindByPlatform(
		request.Context(),
		platform,
	)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	out := []dtos.CommunityResponse{}
	for _, community := range result {
		out = append(out, dtos.CommunityDTO(community))
	}

	respondJSON(writer, http.StatusOK, out)
}

func (api *api) countCommunities(
	writer http.ResponseWriter,
	request *http.Request,
) {
	count, err := api.application.Communities().Count(request.Context())
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	respondJSON(writer, http.StatusOK, map[string]int64{
		"count": count,
	})
}
