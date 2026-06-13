package hatchet

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	hatchetsdk "github.com/hatchet-dev/hatchet/sdks/go"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/replies"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/posts/contents/threads"
	"github.com/steve-rodrigue/aabs/services/posts/restapis/domain/users"
)

func TestNewPostService(t *testing.T) {
	client := &hatchetsdk.Client{}

	service := NewPostService(client)

	if service == nil {
		t.Fatalf("expected service")
	}
}

func TestPostServiceSavePushesPostSavedEvent(t *testing.T) {
	token := os.Getenv("HATCHET_TEST_CLIENT_TOKEN")
	hostPort := os.Getenv("HATCHET_TEST_HOST_PORT")
	tlsStrategy := os.Getenv("HATCHET_TEST_TLS_STRATEGY")

	if token == "" || token == "test-token" {
		t.Skip("HATCHET_TEST_CLIENT_TOKEN not set")
	}

	if hostPort == "" {
		t.Skip("HATCHET_TEST_HOST_PORT not set")
	}

	if tlsStrategy == "" {
		tlsStrategy = "none"
	}

	t.Setenv("HATCHET_CLIENT_TOKEN", token)
	t.Setenv("HATCHET_CLIENT_HOST_PORT", hostPort)
	t.Setenv("HATCHET_CLIENT_TLS_STRATEGY", tlsStrategy)

	client, err := hatchetsdk.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	service := NewPostService(client)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	post := newTestHatchetPost(t)

	err = service.Save(ctx, post)
	if err != nil {
		t.Fatal(err)
	}
}

func newTestHatchetPost(
	t *testing.T,
) posts.Post {
	t.Helper()

	creator := users.NewMockUser(
		"hatchet-user",
		"Hatchet User",
	)

	contentAdapter := contents.NewAdapter(
		replies.NewAdapter(),
		threads.NewAdapter(),
	)

	post, err := posts.NewAdapter(contentAdapter).ToDomain(
		posts.PostInput{
			Identifier: uuid.New(),
			CommunityIDs: []uuid.UUID{
				uuid.New(),
			},
			Creator: creator,
			Content: contents.ContentInput{
				Identifier: uuid.New(),
				Thread: &threads.ThreadInput{
					Identifier: uuid.New(),
					Creator:    creator,
					Title:      "Hatchet test thread",
					Text:       "hello from hatchet test",
				},
				CreatedAt: time.Now().UTC(),
			},
			CreatedOn: time.Now().UTC(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	return post
}
