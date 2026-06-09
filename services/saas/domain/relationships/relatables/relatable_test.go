package relatables

import (
	"testing"

	"github.com/google/uuid"
)

func TestRelatable(t *testing.T) {
	id := uuid.New()

	relatable := &relatable{
		identifier:       id,
		relationshipKind: PostKind,
	}

	if relatable.Identifier() != id {
		t.Fatalf("expected identifier %s, got %s", id, relatable.Identifier())
	}

	if relatable.RelationshipKind() != PostKind {
		t.Fatalf(
			"expected relationship kind %s, got %s",
			PostKind,
			relatable.RelationshipKind(),
		)
	}
}
