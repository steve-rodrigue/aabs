package relationships

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/relatables"
	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
)

var errTest = errors.New("test error")

func TestBuild(t *testing.T) {
	fixture := newApplicationFixture()

	source := relatables.NewMockRelatable(uuid.New(), relatables.PostKind)
	target := relatables.NewMockRelatable(uuid.New(), relatables.PostKind)
	relationship := domain_relationships.NewMockRelationship()

	fixture.builder.BuildValue = []domain_relationships.Relationship{
		relationship,
	}

	result, err := fixture.application.Build(
		source,
		[]relatables.Relatable{target},
	)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.builder.BuildCalls != 1 {
		t.Fatalf("expected 1 build call")
	}

	if fixture.builder.LastSource != source {
		t.Fatalf("expected source to be passed to builder")
	}

	if len(fixture.builder.LastTargets) != 1 ||
		fixture.builder.LastTargets[0] != target {
		t.Fatalf("expected target to be passed to builder")
	}

	if len(result) != 1 || result[0] != relationship {
		t.Fatalf("expected relationship result")
	}
}

func TestBuildReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.builder.BuildErr = errTest

	_, err := fixture.application.Build(
		relatables.NewMockRelatable(uuid.New(), relatables.PostKind),
		nil,
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected build error, got %v", err)
	}
}

func TestSync(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	err := fixture.application.Sync(ctx, []domain_relationships.Relationship{
		domain_relationships.NewMockRelationship(),
		domain_relationships.NewMockRelationship(),
	})

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.SaveCalls != 2 {
		t.Fatalf("expected 2 save calls, got %d", fixture.repository.SaveCalls)
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}
}

func TestSyncReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.SaveErr = errTest

	err := fixture.application.Sync(ctx, []domain_relationships.Relationship{
		domain_relationships.NewMockRelationship(),
	})

	if !errors.Is(err, errTest) {
		t.Fatalf("expected sync error, got %v", err)
	}
}

func TestFind(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	relationship := domain_relationships.NewMockRelationship()
	fixture.repository.FindValue = []domain_relationships.Relationship{
		relationship,
	}

	result, err := fixture.application.Find(ctx, 0, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindCalls != 1 {
		t.Fatalf("expected 1 find call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastIndex != 0 {
		t.Fatalf("expected index 0")
	}

	if fixture.repository.LastAmount != 25 {
		t.Fatalf("expected amount 25")
	}

	if len(result) != 1 || result[0] != relationship {
		t.Fatalf("expected relationship result")
	}
}

func TestFindReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindErr = errTest

	_, err := fixture.application.Find(ctx, 0, 25)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find error, got %v", err)
	}
}

func TestFindByID(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	relationship := domain_relationships.NewMockRelationship()
	fixture.repository.Items[relationship.Identifier()] = relationship

	result, err := fixture.application.FindByID(ctx, relationship.Identifier())

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByIDCalls != 1 {
		t.Fatalf("expected 1 find by id call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastID != relationship.Identifier() {
		t.Fatalf("expected id to be passed")
	}

	if result != relationship {
		t.Fatalf("expected relationship result")
	}
}

func TestFindByIDReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindByIDErr = errTest

	_, err := fixture.application.FindByID(ctx, uuid.New())

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find by id error, got %v", err)
	}
}

func TestFindAfter(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	cursor := uuid.New()
	relationship := domain_relationships.NewMockRelationship()

	fixture.repository.FindAfterValue = []domain_relationships.Relationship{
		relationship,
	}

	result, err := fixture.application.FindAfter(ctx, cursor, 25)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindAfterCalls != 1 {
		t.Fatalf("expected 1 find after call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastCursor != cursor {
		t.Fatalf("expected cursor to be passed")
	}

	if fixture.repository.LastAmount != 25 {
		t.Fatalf("expected amount 25")
	}

	if len(result) != 1 || result[0] != relationship {
		t.Fatalf("expected relationship result")
	}
}

func TestFindAfterReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindAfterErr = errTest

	_, err := fixture.application.FindAfter(ctx, uuid.New(), 25)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected find after error, got %v", err)
	}
}

func TestCount(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.CountValue = 123

	result, err := fixture.application.Count(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.CountCalls != 1 {
		t.Fatalf("expected 1 count call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if result != 123 {
		t.Fatalf("expected count 123, got %d", result)
	}
}

func TestCountReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.CountErr = errTest

	_, err := fixture.application.Count(ctx)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected count error, got %v", err)
	}
}

func TestRelationshipsBySource(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	id := uuid.New()
	relationship := domain_relationships.NewMockRelationship()

	fixture.repository.FindBySourceIDValue = []domain_relationships.Relationship{
		relationship,
	}

	result, err := fixture.application.RelationshipsBySource(ctx, id)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindBySourceIDCalls != 1 {
		t.Fatalf("expected 1 find by source id call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastSourceID != id {
		t.Fatalf("expected source id to be passed")
	}

	if len(result) != 1 || result[0] != relationship {
		t.Fatalf("expected relationship result")
	}
}

func TestRelationshipsBySourceReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindBySourceIDErr = errTest

	_, err := fixture.application.RelationshipsBySource(ctx, uuid.New())

	if !errors.Is(err, errTest) {
		t.Fatalf("expected source error, got %v", err)
	}
}

func TestRelationshipsByTarget(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	id := uuid.New()
	relationship := domain_relationships.NewMockRelationship()

	fixture.repository.FindByTargetIDValue = []domain_relationships.Relationship{
		relationship,
	}

	result, err := fixture.application.RelationshipsByTarget(ctx, id)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByTargetIDCalls != 1 {
		t.Fatalf("expected 1 find by target id call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed")
	}

	if fixture.repository.LastTargetID != id {
		t.Fatalf("expected target id to be passed")
	}

	if len(result) != 1 || result[0] != relationship {
		t.Fatalf("expected relationship result")
	}
}

func TestRelationshipsByTargetReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.repository.FindByTargetIDErr = errTest

	_, err := fixture.application.RelationshipsByTarget(ctx, uuid.New())

	if !errors.Is(err, errTest) {
		t.Fatalf("expected target error, got %v", err)
	}
}

func TestRebuildRelationships(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	source := relatables.NewMockRelatable(uuid.New(), relatables.PostKind)
	target := relatables.NewMockRelatable(uuid.New(), relatables.TopicKind)
	relationship := domain_relationships.NewMockRelationship()

	fixture.relatables.FindAfterValue = []relatables.Relatable{
		source,
	}
	fixture.candidates.FindCandidatesValue = []relatables.Relatable{
		target,
	}
	fixture.builder.BuildValue = []domain_relationships.Relationship{
		relationship,
	}

	err := fixture.application.RebuildRelationships(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.relatables.FindAfterCalls != 2 {
		t.Fatalf(
			"expected 2 find after calls, got %d",
			fixture.relatables.FindAfterCalls,
		)
	}

	if fixture.relatables.LastContext != ctx {
		t.Fatalf("expected context to be passed to relatables")
	}

	if fixture.candidates.FindCandidatesCalls != 1 {
		t.Fatalf("expected 1 find candidates call")
	}

	if fixture.candidates.LastContext != ctx {
		t.Fatalf("expected context to be passed to candidates")
	}

	if fixture.candidates.LastSource != source {
		t.Fatalf("expected source to be passed to candidates")
	}

	if fixture.builder.BuildCalls != 1 {
		t.Fatalf("expected 1 build call")
	}

	if fixture.builder.LastSource != source {
		t.Fatalf("expected source to be passed to builder")
	}

	if len(fixture.builder.LastTargets) != 1 ||
		fixture.builder.LastTargets[0] != target {
		t.Fatalf("expected target to be passed to builder")
	}

	if fixture.repository.SaveCalls != 1 {
		t.Fatalf("expected 1 save call")
	}

	if fixture.repository.LastContext != ctx {
		t.Fatalf("expected context to be passed to repository")
	}
}

func TestRebuildRelationshipsReturnsRelatableError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.relatables.FindAfterErr = errTest

	err := fixture.application.RebuildRelationships(ctx)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected relatable error, got %v", err)
	}
}

func TestRebuildRelationshipsReturnsCandidateError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.relatables.FindAfterValue = []relatables.Relatable{
		relatables.NewMockRelatable(uuid.New(), relatables.PostKind),
	}
	fixture.candidates.FindCandidatesErr = errTest

	err := fixture.application.RebuildRelationships(ctx)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected candidate error, got %v", err)
	}
}

func TestRebuildRelationshipsReturnsBuilderError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.relatables.FindAfterValue = []relatables.Relatable{
		relatables.NewMockRelatable(uuid.New(), relatables.PostKind),
	}
	fixture.builder.BuildErr = errTest

	err := fixture.application.RebuildRelationships(ctx)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected builder error, got %v", err)
	}
}

func TestRebuildRelationshipsReturnsSyncError(t *testing.T) {
	fixture := newApplicationFixture()
	ctx := context.Background()

	fixture.relatables.FindAfterValue = []relatables.Relatable{
		relatables.NewMockRelatable(uuid.New(), relatables.PostKind),
	}
	fixture.builder.BuildValue = []domain_relationships.Relationship{
		domain_relationships.NewMockRelationship(),
	}
	fixture.repository.SaveErr = errTest

	err := fixture.application.RebuildRelationships(ctx)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected sync error, got %v", err)
	}
}
