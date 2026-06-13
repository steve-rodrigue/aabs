package dtos

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
	domain_posts "github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/threads"
)

type PostResponse struct {
	Identifier   uuid.UUID       `json:"identifier"`
	CommunityIDs []uuid.UUID     `json:"community_ids"`
	Creator      UserResponse    `json:"creator"`
	Content      ContentResponse `json:"content"`
	CreatedOn    time.Time       `json:"created_on"`
}

type ContentResponse struct {
	Identifier uuid.UUID       `json:"identifier"`
	Kind       string          `json:"kind"`
	Thread     *ThreadResponse `json:"thread,omitempty"`
	Reply      *ReplyResponse  `json:"reply,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

type ThreadResponse struct {
	Identifier uuid.UUID    `json:"identifier"`
	Creator    UserResponse `json:"creator"`
	Title      string       `json:"title"`
	Text       string       `json:"text"`
}

type ReplyResponse struct {
	Identifier uuid.UUID           `json:"identifier"`
	Target     ReplyTargetResponse `json:"target"`
	Text       string              `json:"text"`
}

type ReplyTargetResponse struct {
	Thread *ThreadResponse `json:"thread,omitempty"`
	Reply  *ReplyResponse  `json:"reply,omitempty"`
}

func PostListDTO(items []posts.Post) []PostResponse {
	out := make([]PostResponse, 0, len(items))

	for _, item := range items {
		out = append(out, PostDTO(item))
	}

	return out
}

func PostDTO(
	post domain_posts.Post,
) PostResponse {
	return PostResponse{
		Identifier:   post.Identifier(),
		CommunityIDs: post.CommunityIDs(),
		Creator:      UserDTO(post.Creator()),
		Content:      ContentDTO(post.Content()),
		CreatedOn:    post.CreatedOn(),
	}
}

func ContentDTO(
	content contents.Content,
) ContentResponse {
	response := ContentResponse{
		Identifier: content.Identifier(),
		CreatedAt:  content.CreatedAt(),
	}

	if content.IsThread() {
		response.Kind = "thread"
		response.Thread = ThreadDTO(
			content.Thread(),
		)
	}

	if content.IsReply() {
		response.Kind = "reply"
		response.Reply = ReplyDTO(
			content.Reply(),
		)
	}

	return response
}

func ThreadDTO(
	thread threads.Thread,
) *ThreadResponse {
	return &ThreadResponse{
		Identifier: thread.Identifier(),
		Creator:    UserDTO(thread.Creator()),
		Title:      thread.Title(),
		Text:       thread.Text(),
	}
}

func ReplyDTO(
	reply replies.Reply,
) *ReplyResponse {
	return &ReplyResponse{
		Identifier: reply.Identifier(),
		Target:     ReplyTargetDTO(reply.Target()),
		Text:       reply.Text(),
	}
}

func ReplyTargetDTO(
	target replies.Target,
) ReplyTargetResponse {
	response := ReplyTargetResponse{}

	if target.IsThread() {
		response.Thread = ThreadDTO(
			target.Thread(),
		)
	}

	if target.IsReply() {
		response.Reply = ReplyDTO(
			target.Reply(),
		)
	}

	return response
}
