package clusters

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
)

type detector struct {
	adapter     Adapter
	comparables clusterables.ComparableRepository
}

func (detector *detector) Detect(
	ctx context.Context,
	target clusterables.Clusterable,
	members []clusterables.Clusterable,
) ([]Cluster, error) {
	if target == nil {
		return nil, ErrInvalidClusterDetectorTarget
	}

	if len(members) == 0 {
		return []Cluster{}, nil
	}

	comparables := make([]clusterables.Comparable, 0, len(members))

	var memberKind clusterables.Kind

	for index, member := range members {
		if member == nil {
			return nil, ErrInvalidClusterDetectorMember
		}

		if index == 0 {
			memberKind = member.ClusterKind()
		}

		if member.ClusterKind() != memberKind {
			return nil, ErrInvalidClusterDetectorMemberKind
		}

		comparable, err := detector.comparables.FindByID(
			ctx,
			member.Identifier(),
		)
		if err != nil {
			return nil, err
		}

		if comparable == nil {
			return nil, ErrInvalidClusterDetectorComparable
		}

		comparables = append(comparables, comparable)
	}

	centroid, err := calculateCentroid(comparables)
	if err != nil {
		return nil, err
	}

	confidenceScore, err := calculateConfidenceScore(
		centroid,
		comparables,
	)
	if err != nil {
		return nil, err
	}

	memberIDs := make([]uuid.UUID, 0, len(members))

	for _, member := range members {
		memberIDs = append(memberIDs, member.Identifier())
	}

	cluster, err := detector.adapter.ToDomain(
		ClusterInput{
			Identifier: uuid.New(),
			Target: clusterables.ClusterableInput{
				Identifier:  target.Identifier(),
				ClusterKind: target.ClusterKind(),
			},
			MemberIDs:       memberIDs,
			MemberKind:      memberKind,
			ConfidenceScore: confidenceScore,
			Centroid:        centroid,
			CreatedOn:       time.Now().UTC(),
		},
	)
	if err != nil {
		return nil, err
	}

	return []Cluster{
		cluster,
	}, nil
}

func calculateCentroid(
	members []clusterables.Comparable,
) ([]float32, error) {
	if len(members) == 0 {
		return nil, ErrInvalidClusterDetectorMember
	}

	firstVector := members[0].Vector()
	if len(firstVector) == 0 {
		return nil, ErrInvalidClusterDetectorVectorSize
	}

	centroid := make([]float32, len(firstVector))

	for _, member := range members {
		vector := member.Vector()

		if len(vector) != len(centroid) {
			return nil, ErrInvalidClusterDetectorVectorSize
		}

		for index := range vector {
			centroid[index] += vector[index]
		}
	}

	for index := range centroid {
		centroid[index] = centroid[index] / float32(len(members))
	}

	return centroid, nil
}

func calculateConfidenceScore(
	centroid []float32,
	members []clusterables.Comparable,
) (float64, error) {
	if len(centroid) == 0 || len(members) == 0 {
		return 0, ErrInvalidClusterDetectorVectorSize
	}

	var total float64

	for _, member := range members {
		vector := member.Vector()

		if len(vector) != len(centroid) {
			return 0, ErrInvalidClusterDetectorVectorSize
		}

		total += cosineSimilarity(
			centroid,
			vector,
		)
	}

	return total / float64(len(members)), nil
}

func cosineSimilarity(
	source []float32,
	target []float32,
) float64 {
	var dot float64
	var sourceMagnitude float64
	var targetMagnitude float64

	for index := range source {
		sourceValue := float64(source[index])
		targetValue := float64(target[index])

		dot += sourceValue * targetValue
		sourceMagnitude += sourceValue * sourceValue
		targetMagnitude += targetValue * targetValue
	}

	if sourceMagnitude == 0 || targetMagnitude == 0 {
		return 0
	}

	return dot / (math.Sqrt(sourceMagnitude) * math.Sqrt(targetMagnitude))
}
