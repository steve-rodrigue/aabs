package clusters

import (
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
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
		memberIDs:       memberIDs,
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
	return cluster.memberIDs
}

func (cluster *MockCluster) MemberKind() clusterables.Kind {
	return cluster.memberKind
}

func (cluster *MockCluster) ConfidenceScore() float64 {
	return cluster.confidenceScore
}

func (cluster *MockCluster) Centroid() []float32 {
	return cluster.centroid
}

func (cluster *MockCluster) CreatedOn() time.Time {
	return time.Time{}
}

type MockClusterRepository struct {
	SaveCalls int
	SaveErr   error

	Items map[uuid.UUID]Cluster

	FindByIDCalls int
	FindByIDErr   error

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
}

func (repository *MockClusterRepository) Save(
	cluster Cluster,
) error {
	repository.SaveCalls++

	return repository.SaveErr
}

func (repository *MockClusterRepository) FindByID(
	id uuid.UUID,
) (Cluster, error) {
	repository.FindByIDCalls++

	if repository.FindByIDErr != nil {
		return nil, repository.FindByIDErr
	}

	return repository.Items[id], nil
}

func (repository *MockClusterRepository) FindByTarget(
	target uuid.UUID,
) ([]Cluster, error) {
	repository.FindByTargetCalls++

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
	member uuid.UUID,
) ([]Cluster, error) {
	repository.FindByMemberCalls++

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
	index int,
	amount int,
) ([]Cluster, error) {
	repository.FindCalls++

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
	cursor uuid.UUID,
	amount int,
) ([]Cluster, error) {
	repository.FindAfterCalls++

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

func (repository *MockClusterRepository) Count() (int64, error) {
	repository.CountCalls++

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

	LastTarget  clusterables.Clusterable
	LastMembers []clusterables.Clusterable
}

func (detector *MockClusterDetector) Detect(
	target clusterables.Clusterable,
	members []clusterables.Clusterable,
) ([]Cluster, error) {
	detector.DetectCalls++
	detector.LastTarget = target
	detector.LastMembers = members

	return detector.DetectValue, detector.DetectErr
}
