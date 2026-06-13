package posts

import (
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

type post struct {
	identifier uuid.UUID

	communityIDs []uuid.UUID
	creator      users.User
	content      contents.Content

	createdOn time.Time
}

func (post *post) Identifier() uuid.UUID {
	return post.identifier
}

func (post *post) CommunityIDs() []uuid.UUID {
	out := make([]uuid.UUID, len(post.communityIDs))
	copy(out, post.communityIDs)

	return out
}

func (post *post) Creator() users.User {
	return post.creator
}

func (post *post) Content() contents.Content {
	return post.content
}

func (post *post) CreatedOn() time.Time {
	return post.createdOn
}
