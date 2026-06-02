package relationships

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	domain_relationships "github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

var errTest = errors.New("test error")

func TestBuild(t *testing.T) {
	fixture := newApplicationFixture()

	source := relatables.NewMockRelatable(uuid.New(), relatables.UserKind)
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
		t.Fatalf("expected 1 build call, got %d", fixture.builder.BuildCalls)
	}

	if len(result) != 1 || result[0] != relationship {
		t.Fatalf("expected built relationship")
	}
}

func TestBuildReturnsError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.builder.BuildErr = errTest

	source := relatables.NewMockRelatable(uuid.New(), relatables.UserKind)
	target := relatables.NewMockRelatable(uuid.New(), relatables.PostKind)

	_, err := fixture.application.Build(
		source,
		[]relatables.Relatable{target},
	)

	if !errors.Is(err, errTest) {
		t.Fatalf("expected build error, got %v", err)
	}
}

func TestSync(t *testing.T) {
	fixture := newApplicationFixture()

	err := fixture.application.Sync([]domain_relationships.Relationship{
		domain_relationships.NewMockRelationship(),
		domain_relationships.NewMockRelationship(),
	})

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.SaveCalls != 2 {
		t.Fatalf("expected 2 save calls, got %d", fixture.repository.SaveCalls)
	}
}

func TestSyncReturnsSaveError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.repository.SaveErr = errTest

	err := fixture.application.Sync([]domain_relationships.Relationship{
		domain_relationships.NewMockRelationship(),
	})

	if !errors.Is(err, errTest) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func TestFindAll(t *testing.T) {
	fixture := newApplicationFixture()

	relationship := domain_relationships.NewMockRelationship()
	fixture.repository.FindAllValue = []domain_relationships.Relationship{
		relationship,
	}

	result, err := fixture.application.FindAll()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindAllCalls != 1 {
		t.Fatalf("expected 1 find all call")
	}

	if len(result) != 1 || result[0] != relationship {
		t.Fatalf("expected relationship result")
	}
}

func TestRelationshipsBySource(t *testing.T) {
	fixture := newApplicationFixture()

	source := uuid.New()
	relationship := domain_relationships.NewMockRelationship()
	fixture.repository.FindBySourceIDValue = []domain_relationships.Relationship{
		relationship,
	}

	result, err := fixture.application.RelationshipsBySource(source)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindBySourceIDCalls != 1 {
		t.Fatalf("expected 1 find by source id call")
	}

	if len(result) != 1 || result[0] != relationship {
		t.Fatalf("expected relationship result")
	}
}

func TestRelationshipsByTarget(t *testing.T) {
	fixture := newApplicationFixture()

	target := uuid.New()
	relationship := domain_relationships.NewMockRelationship()
	fixture.repository.FindByTargetIDValue = []domain_relationships.Relationship{
		relationship,
	}

	result, err := fixture.application.RelationshipsByTarget(target)

	if err != nil {
		t.Fatal(err)
	}

	if fixture.repository.FindByTargetIDCalls != 1 {
		t.Fatalf("expected 1 find by target id call")
	}

	if len(result) != 1 || result[0] != relationship {
		t.Fatalf("expected relationship result")
	}
}

func TestRebuildRelationships(t *testing.T) {
	fixture := newApplicationFixture()

	first := relatables.NewMockRelatable(uuid.New(), relatables.UserKind)
	second := relatables.NewMockRelatable(uuid.New(), relatables.PostKind)
	third := relatables.NewMockRelatable(uuid.New(), relatables.CampaignKind)

	fixture.relatableRepository.Items = []relatables.Relatable{
		first,
		second,
		third,
	}

	fixture.builder.BuildValue = []domain_relationships.Relationship{
		domain_relationships.NewMockRelationship(),
	}

	err := fixture.application.RebuildRelationships()

	if err != nil {
		t.Fatal(err)
	}

	if fixture.relatableRepository.FindAllCalls != 1 {
		t.Fatalf("expected 1 relatable find all call")
	}

	if fixture.builder.BuildCalls != 3 {
		t.Fatalf("expected 3 build calls, got %d", fixture.builder.BuildCalls)
	}

	if fixture.repository.SaveCalls != 3 {
		t.Fatalf("expected 3 save calls, got %d", fixture.repository.SaveCalls)
	}
}

func TestRebuildRelationshipsReturnsRelatableRepositoryError(t *testing.T) {
	fixture := newApplicationFixture()
	fixture.relatableRepository.FindAllErr = errTest

	err := fixture.application.RebuildRelationships()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected relatable repository error, got %v", err)
	}

	if fixture.builder.BuildCalls != 0 {
		t.Fatalf("expected builder not to be called")
	}
}

func TestRebuildRelationshipsReturnsBuilderError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.relatableRepository.Items = []relatables.Relatable{
		relatables.NewMockRelatable(uuid.New(), relatables.UserKind),
		relatables.NewMockRelatable(uuid.New(), relatables.PostKind),
	}

	fixture.builder.BuildErr = errTest

	err := fixture.application.RebuildRelationships()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected builder error, got %v", err)
	}

	if fixture.repository.SaveCalls != 0 {
		t.Fatalf("expected save not to be called")
	}
}

func TestRebuildRelationshipsReturnsSyncError(t *testing.T) {
	fixture := newApplicationFixture()

	fixture.relatableRepository.Items = []relatables.Relatable{
		relatables.NewMockRelatable(uuid.New(), relatables.UserKind),
		relatables.NewMockRelatable(uuid.New(), relatables.PostKind),
	}

	fixture.builder.BuildValue = []domain_relationships.Relationship{
		domain_relationships.NewMockRelationship(),
	}

	fixture.repository.SaveErr = errTest

	err := fixture.application.RebuildRelationships()

	if !errors.Is(err, errTest) {
		t.Fatalf("expected sync error, got %v", err)
	}
}
