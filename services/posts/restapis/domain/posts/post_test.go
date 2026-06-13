package posts

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

func TestPost(t *testing.T) {
	id := uuid.New()
	communityID := uuid.New()
	creator := users.NewMockUser("@user", "User")
	content := &contents.MockContent{
		TextValue: "Post text",
	}
	createdOn := time.Now().UTC()

	post := &post{
		identifier: id,
		communityIDs: []uuid.UUID{
			communityID,
		},
		creator:   creator,
		content:   content,
		createdOn: createdOn,
	}

	if post.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, post.Identifier())
	}

	if len(post.CommunityIDs()) != 1 || post.CommunityIDs()[0] != communityID {
		t.Fatalf("expected community id")
	}

	if post.Creator() != creator {
		t.Fatalf("expected creator")
	}

	if post.Content() != content {
		t.Fatalf("expected content")
	}

	if !post.CreatedOn().Equal(createdOn) {
		t.Fatalf("expected created on %s, got %s", createdOn, post.CreatedOn())
	}
}

func TestPostCommunityIDsReturnsCopy(t *testing.T) {
	communityID := uuid.New()

	post := &post{
		communityIDs: []uuid.UUID{
			communityID,
		},
	}

	ids := post.CommunityIDs()
	ids[0] = uuid.New()

	if post.CommunityIDs()[0] != communityID {
		t.Fatalf("expected community ids copy not to mutate original")
	}
}
