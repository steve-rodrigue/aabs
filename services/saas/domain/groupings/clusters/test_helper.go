package clusters

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/steve-rodrigue/aabs/services/saas/domain/concepts/clusterables"
)

func NewMockCluster(
	target clusterables.Clusterable,
	memberKind clusterables.Kind,
	memberIDs []uuid.UUID,
) *MockCluster {
	return &MockCluster{
		id:              uuid.New(),
		target:          target,
		memberKind:      memberKind,
		memberIDs:       copyMockUUIDs(memberIDs),
		confidenceScore: 1,
		centroid:        []float32{},
	}
}

func NewMockClusterRepository() *MockClusterRepository {
	return &MockClusterRepository{
		Items: map[uuid.UUID]Cluster{},
	}
}

func NewMockClusterDetector() *MockClusterDetector {
	return &MockClusterDetector{}
}

func NewMockClusterAdapter() *MockClusterAdapter {
	return &MockClusterAdapter{}
}

type MockCluster struct {
	id uuid.UUID

	target     clusterables.Clusterable
	memberIDs  []uuid.UUID
	memberKind clusterables.Kind

	confidenceScore float64
	centroid        []float32
}

func (cluster *MockCluster) Identifier() uuid.UUID {
	return cluster.id
}

func (cluster *MockCluster) Target() clusterables.Clusterable {
	return cluster.target
}

func (cluster *MockCluster) MemberIDs() []uuid.UUID {
	return copyMockUUIDs(cluster.memberIDs)
}

func (cluster *MockCluster) MemberKind() clusterables.Kind {
	return cluster.memberKind
}

func (cluster *MockCluster) ConfidenceScore() float64 {
	return cluster.confidenceScore
}

func (cluster *MockCluster) Centroid() []float32 {
	return copyMockFloat32s(cluster.centroid)
}

func (cluster *MockCluster) CreatedOn() time.Time {
	return time.Time{}
}

type MockClusterAdapter struct {
	ToDomainCalls int
	ToDomainErr   error
	ToDomainValue Cluster

	LastInput ClusterInput
}

func (adapter *MockClusterAdapter) ToDomain(
	input ClusterInput,
) (Cluster, error) {
	adapter.ToDomainCalls++
	adapter.LastInput = input

	if adapter.ToDomainErr != nil {
		return nil, adapter.ToDomainErr
	}

	if adapter.ToDomainValue != nil {
		return adapter.ToDomainValue, nil
	}

	return &MockCluster{
		id: input.Identifier,
		target: clusterables.NewMockClusterableWithID(
			input.Target.Identifier,
			input.Target.ClusterKind,
		),
		memberKind:      input.MemberKind,
		memberIDs:       copyMockUUIDs(input.MemberIDs),
		confidenceScore: input.ConfidenceScore,
		centroid:        copyMockFloat32s(input.Centroid),
	}, nil
}

type MockClusterRepository struct {
	SaveCalls int
	SaveErr   error

	LastContext context.Context
	LastCluster Cluster

	Items map[uuid.UUID]Cluster

	FindByIDCalls int
	FindByIDErr   error
	FindByIDValue Cluster

	FindByTargetCalls int
	FindByTargetErr   error
	FindByTargetValue []Cluster

	FindByMemberCalls int
	FindByMemberErr   error
	FindByMemberValue []Cluster

	FindCalls int
	FindErr   error
	FindValue []Cluster

	FindAfterCalls int
	FindAfterErr   error
	FindAfterValue []Cluster

	CountCalls int
	CountErr   error
	CountValue int64

	LastID     uuid.UUID
	LastTarget uuid.UUID
	LastMember uuid.UUID
	LastIndex  int
	LastAmount int
	LastCursor uuid.UUID
}

func (repository *MockClusterRepository) Save(
	ctx context.Context,
	cluster Cluster,
) error {
	repository.SaveCalls++
	repository.LastContext = ctx
	repository.LastCluster = cluster

	if repository.Items != nil && cluster != nil {
		repository.Items[cluster.Identifier()] = cluster
	}

	return repository.SaveErr
}

func (repository *MockClusterRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (Cluster, error) {
	repository.FindByIDCalls++
	repository.LastContext = ctx
	repository.LastID = id

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	if repository.FindByIDValue != nil {
		return repository.FindByIDValue, nil
	}

	if repository.Items == nil {
		return nil, nil
	}

	return repository.Items[id], nil
}

func (repository *MockClusterRepository) FindByTarget(
	ctx context.Context,
	target uuid.UUID,
) ([]Cluster, error) {
	repository.FindByTargetCalls++
	repository.LastContext = ctx
	repository.LastTarget = target

	if repository.FindByTargetErr != nil {
		return nil, repository.FindByTargetErr
	}

	if repository.FindByTargetValue != nil {
		return repository.FindByTargetValue, nil
	}

	out := []Cluster{}

	for _, cluster := range repository.Items {
		if cluster.Target() == nil {
			continue
		}

		if cluster.Target().Identifier() == target {
			out = append(out, cluster)
		}
	}

	return out, nil
}

func (repository *MockClusterRepository) FindByMember(
	ctx context.Context,
	member uuid.UUID,
) ([]Cluster, error) {
	repository.FindByMemberCalls++
	repository.LastContext = ctx
	repository.LastMember = member

	if repository.FindByMemberErr != nil {
		return nil, repository.FindByMemberErr
	}

	if repository.FindByMemberValue != nil {
		return repository.FindByMemberValue, nil
	}

	out := []Cluster{}

	for _, cluster := range repository.Items {
		for _, memberID := range cluster.MemberIDs() {
			if memberID == member {
				out = append(out, cluster)
				break
			}
		}
	}

	return out, nil
}

func (repository *MockClusterRepository) Find(
	ctx context.Context,
	index int,
	amount int,
) ([]Cluster, error) {
	repository.FindCalls++
	repository.LastContext = ctx
	repository.LastIndex = index
	repository.LastAmount = amount

	if repository.FindErr != nil {
		return nil, repository.FindErr
	}

	if repository.FindValue != nil {
		return repository.FindValue, nil
	}

	clusters := repository.sortedClusters()

	if index >= len(clusters) {
		return []Cluster{}, nil
	}

	end := index + amount
	if end > len(clusters) {
		end = len(clusters)
	}

	return clusters[index:end], nil
}

func (repository *MockClusterRepository) FindAfter(
	ctx context.Context,
	cursor uuid.UUID,
	amount int,
) ([]Cluster, error) {
	repository.FindAfterCalls++
	repository.LastContext = ctx
	repository.LastCursor = cursor
	repository.LastAmount = amount

	if repository.FindAfterErr != nil {
		return nil, repository.FindAfterErr
	}

	if repository.FindAfterValue != nil {
		return repository.FindAfterValue, nil
	}

	clusters := repository.sortedClusters()

	start := 0

	if cursor != uuid.Nil {
		for index, cluster := range clusters {
			if cluster.Identifier() == cursor {
				start = index + 1
				break
			}
		}
	}

	if start >= len(clusters) {
		return []Cluster{}, nil
	}

	end := start + amount
	if end > len(clusters) {
		end = len(clusters)
	}

	return clusters[start:end], nil
}

func (repository *MockClusterRepository) Count(
	ctx context.Context,
) (int64, error) {
	repository.CountCalls++
	repository.LastContext = ctx

	if repository.CountErr != nil {
		return 0, repository.CountErr
	}

	if repository.CountValue != 0 {
		return repository.CountValue, nil
	}

	return int64(len(repository.Items)), nil
}

func (repository *MockClusterRepository) sortedClusters() []Cluster {
	out := make([]Cluster, 0, len(repository.Items))

	for _, cluster := range repository.Items {
		out = append(out, cluster)
	}

	sort.Slice(out, func(left int, right int) bool {
		return out[left].Identifier().String() <
			out[right].Identifier().String()
	})

	return out
}

type MockClusterDetector struct {
	DetectCalls int
	DetectErr   error
	DetectValue []Cluster

	LastContext context.Context
	LastTarget  clusterables.Clusterable
	LastMembers []clusterables.Clusterable
}

func (detector *MockClusterDetector) Detect(
	ctx context.Context,
	target clusterables.Clusterable,
	members []clusterables.Clusterable,
) ([]Cluster, error) {
	detector.DetectCalls++
	detector.LastContext = ctx
	detector.LastTarget = target
	detector.LastMembers = members

	return detector.DetectValue, detector.DetectErr
}

func copyMockUUIDs(
	values []uuid.UUID,
) []uuid.UUID {
	out := make([]uuid.UUID, len(values))
	copy(out, values)

	return out
}

func copyMockFloat32s(
	values []float32,
) []float32 {
	out := make([]float32, len(values))
	copy(out, values)

	return out
}
