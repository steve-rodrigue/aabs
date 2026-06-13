package inputs

import (
	"time"

	"github.com/google/uuid"
)

// SavePostRequest represents a save post request input
type SavePostRequest struct {
	Identifier   uuid.UUID    `json:"identifier"`
	CommunityIDs []uuid.UUID  `json:"community_ids"`
	CreatorID    uuid.UUID    `json:"creator_id"`
	Content      ContentInput `json:"content"`
	CreatedOn    time.Time    `json:"created_on"`
}

// ContentInput represents a post's content request input
type ContentInput struct {
	Identifier uuid.UUID    `json:"identifier"`
	Kind       string       `json:"kind"`
	Thread     *ThreadInput `json:"thread,omitempty"`
	Reply      *ReplyInput  `json:"reply,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
}

// ThreadInput represents a post's thread request input
type ThreadInput struct {
	Identifier uuid.UUID `json:"identifier"`
	Title      string    `json:"title"`
	Text       string    `json:"text"`
}

// ThreadReplyInputInput represents a reply's thread request input
type ReplyInput struct {
	Identifier     uuid.UUID  `json:"identifier"`
	TargetReplyID  *uuid.UUID `json:"target_reply_id,omitempty"`
	TargetThreadID *uuid.UUID `json:"target_thread_id,omitempty"`
	Text           string     `json:"text"`
}
