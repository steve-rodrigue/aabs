package posts

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"

	client_internal "github.com/steve-rodrigue/aabs/services/posts/restapis/clients/internal"
	domain_platforms "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/platforms"
	domain_posts "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/threads"
	domain_users "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/dtos"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/inputs"
)

type application struct {
	baseURL string
	client  *http.Client
}

func (app *application) Save(
	ctx context.Context,
	post domain_posts.Post,
) error {
	return client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodPost,
		app.baseURL+"/posts",
		savePostRequest(post),
		nil,
	)
}

func (app *application) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (domain_posts.Post, error) {
	var response dtos.PostResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/posts/"+id.String(),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return postResponseToDomain(response)
}

func (app *application) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]domain_posts.Post, error) {
	values := url.Values{}
	values.Set("index", strconv.Itoa(index))
	values.Set("amount", strconv.Itoa(amount))

	var response []dtos.PostResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/posts?"+values.Encode(),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return postsResponseToDomain(response)
}

func (app *application) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]domain_posts.Post, error) {
	values := url.Values{}
	values.Set("amount", strconv.Itoa(amount))

	if cursor != uuid.Nil {
		values.Set("cursor", cursor.String())
	}

	var response []dtos.PostResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodGet,
		app.baseURL+"/posts?"+values.Encode(),
		nil,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return postsResponseToDomain(response)
}

func (app *application) FindByCriteria(
	ctx context.Context,
	criteria domain_posts.Criteria,
	index int,
	amount int,
) ([]domain_posts.Post, error) {
	values := url.Values{}
	values.Set("index", strconv.Itoa(index))
	values.Set("amount", strconv.Itoa(amount))

	var response []dtos.PostResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodPost,
		app.baseURL+"/posts/search?"+values.Encode(),
		criteria,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return postsResponseToDomain(response)
}

func (app *application) FindByCriteriaAfter(
	ctx context.Context,
	criteria domain_posts.Criteria,
	cursor uuid.UUID,
	amount int,
) ([]domain_posts.Post, error) {
	values := url.Values{}
	values.Set("amount", strconv.Itoa(amount))

	if cursor != uuid.Nil {
		values.Set("cursor", cursor.String())
	}

	var response []dtos.PostResponse

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodPost,
		app.baseURL+"/posts/search?"+values.Encode(),
		criteria,
		&response,
	)
	if err != nil {
		return nil, err
	}

	return postsResponseToDomain(response)
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
		app.baseURL+"/posts/count",
		nil,
		&response,
	)
	if err != nil {
		return 0, err
	}

	return response.Count, nil
}

func (app *application) CountByCriteria(
	ctx context.Context,
	criteria domain_posts.Criteria,
) (int64, error) {
	var response struct {
		Count int64 `json:"count"`
	}

	err := client_internal.DoJSON(
		ctx,
		app.client,
		http.MethodPost,
		app.baseURL+"/posts/search/count",
		criteria,
		&response,
	)
	if err != nil {
		return 0, err
	}

	return response.Count, nil
}

func savePostRequest(
	post domain_posts.Post,
) inputs.SavePostRequest {
	return inputs.SavePostRequest{
		Identifier:   post.Identifier(),
		CommunityIDs: post.CommunityIDs(),
		CreatorID:    post.Creator().Identifier(),
		Content:      saveContentInput(post.Content()),
		CreatedOn:    post.CreatedOn(),
	}
}

func saveContentInput(
	content contents.Content,
) inputs.ContentInput {
	input := inputs.ContentInput{
		Identifier: content.Identifier(),
		CreatedAt:  content.CreatedAt(),
	}

	if content.IsThread() {
		input.Kind = "thread"
		input.Thread = &inputs.ThreadInput{
			Identifier: content.Thread().Identifier(),
			Title:      content.Thread().Title(),
			Text:       content.Thread().Text(),
		}
	}

	if content.IsReply() {
		input.Kind = "reply"

		reply := content.Reply()

		input.Reply = &inputs.ReplyInput{
			Identifier: reply.Identifier(),
			Text:       reply.Text(),
		}

		if reply.Target().IsReply() {
			id := reply.Target().Reply().Identifier()
			input.Reply.TargetReplyID = &id
		}

		if reply.Target().IsThread() {
			id := reply.Target().Thread().Identifier()
			input.Reply.TargetThreadID = &id
		}
	}

	return input
}

func postResponseToDomain(
	response dtos.PostResponse,
) (domain_posts.Post, error) {
	creator, err := userResponseToDomain(response.Creator)
	if err != nil {
		return nil, err
	}

	content, err := contentResponseToInput(response.Content)
	if err != nil {
		return nil, err
	}

	return domain_posts.NewAdapter(contentAdapter()).ToDomain(
		domain_posts.PostInput{
			Identifier:   response.Identifier,
			CommunityIDs: response.CommunityIDs,
			Creator:      creator,
			Content:      content,
			CreatedOn:    response.CreatedOn,
		},
	)
}

func postsResponseToDomain(
	input []dtos.PostResponse,
) ([]domain_posts.Post, error) {
	out := make([]domain_posts.Post, 0, len(input))

	for _, response := range input {
		post, err := postResponseToDomain(response)
		if err != nil {
			return nil, err
		}

		out = append(out, post)
	}

	return out, nil
}

func contentResponseToInput(
	response dtos.ContentResponse,
) (contents.ContentInput, error) {
	input := contents.ContentInput{
		Identifier: response.Identifier,
		CreatedAt:  response.CreatedAt,
	}

	if response.Thread != nil {
		thread, err := threadResponseToInput(*response.Thread)
		if err != nil {
			return contents.ContentInput{}, err
		}

		input.Thread = &thread
	}

	if response.Reply != nil {
		reply, err := replyResponseToInput(*response.Reply)
		if err != nil {
			return contents.ContentInput{}, err
		}

		input.Reply = &reply
	}

	return input, nil
}

func threadResponseToInput(
	response dtos.ThreadResponse,
) (threads.ThreadInput, error) {
	creator, err := userResponseToDomain(response.Creator)
	if err != nil {
		return threads.ThreadInput{}, err
	}

	return threads.ThreadInput{
		Identifier: response.Identifier,
		Creator:    creator,
		Title:      response.Title,
		Text:       response.Text,
	}, nil
}

func replyResponseToInput(
	response dtos.ReplyResponse,
) (replies.ReplyInput, error) {
	target, err := replyTargetResponseToInput(response.Target)
	if err != nil {
		return replies.ReplyInput{}, err
	}

	return replies.ReplyInput{
		Identifier: response.Identifier,
		Target:     target,
		Text:       response.Text,
	}, nil
}

func replyTargetResponseToInput(
	response dtos.ReplyTargetResponse,
) (replies.TargetInput, error) {
	if response.Thread != nil {
		threadInput, err := threadResponseToInput(*response.Thread)
		if err != nil {
			return replies.TargetInput{}, err
		}

		thread, err := threads.NewAdapter().ToDomain(threadInput)
		if err != nil {
			return replies.TargetInput{}, err
		}

		return replies.TargetInput{
			Thread: thread,
		}, nil
	}

	if response.Reply != nil {
		replyInput, err := replyResponseToInput(*response.Reply)
		if err != nil {
			return replies.TargetInput{}, err
		}

		reply, err := replies.NewAdapter().ToDomain(replyInput)
		if err != nil {
			return replies.TargetInput{}, err
		}

		return replies.TargetInput{
			Reply: reply,
		}, nil
	}

	return replies.TargetInput{}, nil
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

func contentAdapter() contents.Adapter {
	return contents.NewAdapter(
		replies.NewAdapter(),
		threads.NewAdapter(),
	)
}
