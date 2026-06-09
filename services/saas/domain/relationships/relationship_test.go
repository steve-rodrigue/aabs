package relationships

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

func TestRelationship(t *testing.T) {
	id := uuid.New()

	source := relatables.NewMockRelatable(
		uuid.New(),
		relatables.PostKind,
	)

	target := relatables.NewMockRelatable(
		uuid.New(),
		relatables.UserKind,
	)

	createdOn := time.Now().UTC()

	relationship := &relationship{
		identifier: id,
		source:     source,
		target:     target,
		similarity: 0.75,
		createdOn:  createdOn,
	}

	if relationship.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, relationship.Identifier())
	}

	if relationship.Source() != source {
		t.Fatalf("expected source")
	}

	if relationship.Target() != target {
		t.Fatalf("expected target")
	}

	if relationship.Similarity() != 0.75 {
		t.Fatalf("expected similarity %f, got %f", 0.75, relationship.Similarity())
	}

	if !relationship.CreatedOn().Equal(createdOn) {
		t.Fatalf(
			"expected created on %s, got %s",
			createdOn,
			relationship.CreatedOn(),
		)
	}
}
