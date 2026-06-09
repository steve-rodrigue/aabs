package builders

import (
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships"
	"github.com/steve-rodrigue/aabs/services/saas/domain/relationships/relatables"
)

func NewMockBuilder() *MockBuilder {
	return &MockBuilder{}
}

type MockBuilder struct {
	BuildCalls int
	BuildErr   error
	BuildValue []relationships.Relationship

	LastSource  relatables.Relatable
	LastTargets []relatables.Relatable
}

func (builder *MockBuilder) Build(
	source relatables.Relatable,
	targets []relatables.Relatable,
) ([]relationships.Relationship, error) {
	builder.BuildCalls++
	builder.LastSource = source
	builder.LastTargets = targets

	return builder.BuildValue, builder.BuildErr
}
