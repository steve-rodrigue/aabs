package clusterables

import (
	"context"
)

type candidateRepository struct {
	comparableRepository ComparableRepository
}

func (app *candidateRepository) FindCandidates(
	ctx context.Context,
	target Clusterable,
	kind Kind,
	amount int,
) ([]Clusterable, error) {
	if amount <= 0 {
		return []Clusterable{}, nil
	}

	targetComparable, err := app.comparableRepository.FindByID(
		ctx,
		target.Identifier(),
	)
	if err != nil {
		return nil, err
	}

	if targetComparable == nil {
		return []Clusterable{}, nil
	}

	candidates, err := app.comparableRepository.FindNearest(
		ctx,
		targetComparable,
		kind,
		amount+1,
	)
	if err != nil {
		return nil, err
	}

	out := make([]Clusterable, 0, amount)

	for _, candidate := range candidates {
		if candidate.Identifier() == target.Identifier() &&
			candidate.ClusterKind() == target.ClusterKind() {
			continue
		}

		out = append(out, candidate)

		if len(out) >= amount {
			break
		}
	}

	return out, nil
}
