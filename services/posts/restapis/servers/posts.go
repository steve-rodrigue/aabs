package servers

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/threads"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/dtos"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/inputs"
)

func (api *api) createPost(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var input inputs.SavePostRequest

	if err := decodeJSON(request, &input); err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	creator, err := api.application.Users().FindByID(
		request.Context(),
		input.CreatorID,
	)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if creator == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	contentInput, err := api.postContentInput(
		request,
		input.Content,
		creator,
	)
	if err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	post, err := posts.NewAdapter(
		contents.NewAdapter(
			replies.NewAdapter(),
			threads.NewAdapter(),
		),
	).ToDomain(
		posts.PostInput{
			Identifier:   input.Identifier,
			CommunityIDs: input.CommunityIDs,
			Creator:      creator,
			Content:      contentInput,
			CreatedOn:    input.CreatedOn,
		},
	)
	if err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	if err := api.application.Posts().Save(request.Context(), post); err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	respondJSON(writer, http.StatusCreated, dtos.PostDTO(post))
}

func (api *api) findPostByID(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id, err := routeUUID(request, "id")
	if err != nil {
		respondError(writer, http.StatusBadRequest, err)
		return
	}

	post, err := api.application.Posts().FindByID(
		request.Context(),
		id,
	)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	if post == nil {
		respondJSON(writer, http.StatusNotFound, nil)
		return
	}

	respondJSON(writer, http.StatusOK, dtos.PostDTO(post))
}

func (api *api) findPosts(
	writer http.ResponseWriter,
	request *http.Request,
) {
	index := queryInt(request, "index", 0)
	amount := queryInt(request, "amount", 25)
	cursor := queryUUID(request, "cursor")

	var (
		result []posts.Post
		err    error
	)

	if cursor != uuid.Nil {
		result, err = api.application.Posts().FindAfter(
			request.Context(),
			cursor,
			amount,
		)
	} else {
		result, err = api.application.Posts().Find(
			request.Context(),
			index,
			amount,
		)
	}

	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	respondJSON(writer, http.StatusOK, dtos.PostListDTO(result))
}

func (api *api) findPostsByCriteria(
	writer http.ResponseWriter,
	request *http.Request,
) {
	index := queryInt(request, "index", 0)
	amount := queryInt(request, "amount", 25)
	cursor := queryUUID(request, "cursor")

	criteria := posts.Criteria{
		UserIDs:      queryUUIDs(request, "user_ids"),
		CommunityIDs: queryUUIDs(request, "community_ids"),
		PlatformIDs:  queryUUIDs(request, "platform_ids"),
	}

	var (
		result []posts.Post
		err    error
	)

	if cursor != uuid.Nil {
		result, err = api.application.Posts().FindByCriteriaAfter(
			request.Context(),
			criteria,
			cursor,
			amount,
		)
	} else {
		result, err = api.application.Posts().FindByCriteria(
			request.Context(),
			criteria,
			index,
			amount,
		)
	}

	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	respondJSON(writer, http.StatusOK, dtos.PostListDTO(result))
}

func (api *api) countPosts(
	writer http.ResponseWriter,
	request *http.Request,
) {
	count, err := api.application.Posts().Count(request.Context())
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	respondJSON(writer, http.StatusOK, map[string]int64{
		"count": count,
	})
}

func (api *api) countPostsByCriteria(
	writer http.ResponseWriter,
	request *http.Request,
) {
	count, err := api.application.Posts().CountByCriteria(
		request.Context(),
		posts.Criteria{
			UserIDs:      queryUUIDs(request, "user_ids"),
			CommunityIDs: queryUUIDs(request, "community_ids"),
			PlatformIDs:  queryUUIDs(request, "platform_ids"),
		},
	)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, err)
		return
	}

	respondJSON(writer, http.StatusOK, map[string]int64{
		"count": count,
	})
}

func (api *api) postContentInput(
	request *http.Request,
	input inputs.ContentInput,
	creator users.User,
) (contents.ContentInput, error) {
	out := contents.ContentInput{
		Identifier: input.Identifier,
		CreatedAt:  input.CreatedAt,
	}

	switch input.Kind {
	case "thread":
		if input.Thread == nil {
			return contents.ContentInput{}, posts.ErrInvalidPostContent
		}

		out.Thread = &threads.ThreadInput{
			Identifier: input.Thread.Identifier,
			Creator:    creator,
			Title:      input.Thread.Title,
			Text:       input.Thread.Text,
		}

	case "reply":
		if input.Reply == nil {
			return contents.ContentInput{}, posts.ErrInvalidPostContent
		}

		target, err := api.replyTargetInput(request, input.Reply)
		if err != nil {
			return contents.ContentInput{}, err
		}

		out.Reply = &replies.ReplyInput{
			Identifier: input.Reply.Identifier,
			Target:     target,
			Text:       input.Reply.Text,
		}

	default:
		return contents.ContentInput{}, posts.ErrInvalidPostContent
	}

	return out, nil
}

func (api *api) replyTargetInput(
	request *http.Request,
	input *inputs.ReplyInput,
) (replies.TargetInput, error) {
	if input.TargetThreadID != nil {
		target, err := api.application.Posts().FindByID(
			request.Context(),
			*input.TargetThreadID,
		)
		if err != nil {
			return replies.TargetInput{}, err
		}
		if target == nil || !target.Content().IsThread() {
			return replies.TargetInput{}, replies.ErrInvalidReplyTarget
		}

		return replies.TargetInput{
			Thread: target.Content().Thread(),
		}, nil
	}

	if input.TargetReplyID != nil {
		target, err := api.application.Posts().FindByID(
			request.Context(),
			*input.TargetReplyID,
		)
		if err != nil {
			return replies.TargetInput{}, err
		}
		if target == nil || !target.Content().IsReply() {
			return replies.TargetInput{}, replies.ErrInvalidReplyTarget
		}

		return replies.TargetInput{
			Reply: target.Content().Reply(),
		}, nil
	}

	return replies.TargetInput{}, replies.ErrInvalidReplyTarget
}
