package milvus

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"

	"github.com/steve-rodrigue/aabs/services/saas/domain/groupings/clusters/clusterables"
)

type groupingsClustersClusterablesComparableRepository struct {
	client     client.Client
	adapter    clusterables.ComparableAdapter
	collection string
}

func (repository *groupingsClustersClusterablesComparableRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (clusterables.Comparable, error) {
	result, err := repository.client.Query(
		ctx,
		repository.collection,
		[]string{},
		fmt.Sprintf(
			"%s == %q",
			groupingsClustersClusterablesCandidatesComparableIdentifierField,
			id.String(),
		),
		[]string{
			groupingsClustersClusterablesCandidatesComparableIdentifierField,
			groupingsClustersClusterablesCandidatesComparableKindField,
			groupingsClustersClusterablesCandidatesComparableVectorField,
		},
	)
	if err != nil {
		return nil, err
	}

	if result.Len() == 0 {
		return nil, nil
	}

	return repository.resultToDomain(result, 0)
}

func (repository *groupingsClustersClusterablesComparableRepository) FindByKind(
	ctx context.Context,
	kind clusterables.Kind,
	index int,
	amount int,
) ([]clusterables.Comparable, error) {
	result, err := repository.client.Query(
		ctx,
		repository.collection,
		[]string{},
		fmt.Sprintf(
			"%s == %q",
			groupingsClustersClusterablesCandidatesComparableKindField,
			string(kind),
		),
		[]string{
			groupingsClustersClusterablesCandidatesComparableIdentifierField,
			groupingsClustersClusterablesCandidatesComparableKindField,
			groupingsClustersClusterablesCandidatesComparableVectorField,
		},
		client.WithOffset(int64(index)),
		client.WithLimit(int64(amount)),
	)
	if err != nil {
		return nil, err
	}

	out := make([]clusterables.Comparable, 0, result.Len())

	for row := 0; row < result.Len(); row++ {
		comparable, err := repository.resultToDomain(result, row)
		if err != nil {
			return nil, err
		}

		out = append(out, comparable)
	}

	return out, nil
}

func (repository *groupingsClustersClusterablesComparableRepository) FindNearest(
	ctx context.Context,
	target clusterables.Comparable,
	kind clusterables.Kind,
	amount int,
) ([]clusterables.Comparable, error) {
	if amount <= 0 {
		return []clusterables.Comparable{}, nil
	}

	searchParameters, err := entity.NewIndexFlatSearchParam()
	if err != nil {
		return nil, err
	}

	results, err := repository.client.Search(
		ctx,
		repository.collection,
		[]string{},
		fmt.Sprintf(
			"%s == %q",
			groupingsClustersClusterablesCandidatesComparableKindField,
			string(kind),
		),
		[]string{
			groupingsClustersClusterablesCandidatesComparableIdentifierField,
			groupingsClustersClusterablesCandidatesComparableKindField,
			groupingsClustersClusterablesCandidatesComparableVectorField,
		},
		[]entity.Vector{
			entity.FloatVector(target.Vector()),
		},
		groupingsClustersClusterablesCandidatesComparableVectorField,
		entity.COSINE,
		amount,
		searchParameters,
	)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []clusterables.Comparable{}, nil
	}

	out := make([]clusterables.Comparable, 0, results[0].ResultCount)

	for row := 0; row < results[0].ResultCount; row++ {
		comparable, err := repository.searchResultToDomain(
			results[0],
			row,
		)
		if err != nil {
			return nil, err
		}

		out = append(out, comparable)
	}

	return out, nil
}

func (repository *groupingsClustersClusterablesComparableRepository) searchResultToDomain(
	result client.SearchResult,
	row int,
) (clusterables.Comparable, error) {
	id, err := repository.searchResultStringColumnValue(
		result,
		groupingsClustersClusterablesCandidatesComparableIdentifierField,
		row,
	)
	if err != nil {
		return nil, err
	}

	identifier, err := uuid.Parse(id)
	if err != nil {
		return nil, clusterables.ErrInvalidClusterableIdentifier
	}

	kind, err := repository.searchResultStringColumnValue(
		result,
		groupingsClustersClusterablesCandidatesComparableKindField,
		row,
	)
	if err != nil {
		return nil, err
	}

	vector, err := repository.searchResultVectorColumnValue(
		result,
		groupingsClustersClusterablesCandidatesComparableVectorField,
		row,
	)
	if err != nil {
		return nil, err
	}

	return repository.adapter.ToDomain(
		clusterables.ComparableInput{
			Clusterable: clusterables.ClusterableInput{
				Identifier:  identifier,
				ClusterKind: clusterables.Kind(kind),
			},
			Vector: vector,
		},
	)
}

func (repository *groupingsClustersClusterablesComparableRepository) searchResultColumn(
	result client.SearchResult,
	field string,
) entity.Column {
	for _, column := range result.Fields {
		if column.Name() == field {
			return column
		}
	}

	return nil
}

func (repository *groupingsClustersClusterablesComparableRepository) searchResultStringColumnValue(
	result client.SearchResult,
	field string,
	row int,
) (string, error) {
	column := repository.searchResultColumn(result, field)
	if column == nil {
		return "", fmt.Errorf("missing milvus field %s", field)
	}

	value, err := column.Get(row)
	if err != nil {
		return "", err
	}

	typed, ok := value.(string)
	if !ok {
		return "", fmt.Errorf(
			"invalid milvus field %s type %T",
			field,
			value,
		)
	}

	return typed, nil
}

func (repository *groupingsClustersClusterablesComparableRepository) searchResultVectorColumnValue(
	result client.SearchResult,
	field string,
	row int,
) ([]float32, error) {
	column := repository.searchResultColumn(result, field)
	if column == nil {
		return nil, fmt.Errorf("missing milvus field %s", field)
	}

	value, err := column.Get(row)
	if err != nil {
		return nil, err
	}

	switch typed := value.(type) {
	case []float32:
		return copyVector(typed), nil

	case entity.FloatVector:
		return copyVector([]float32(typed)), nil

	default:
		return nil, fmt.Errorf(
			"invalid milvus field %s type %T",
			field,
			value,
		)
	}
}

func (repository *groupingsClustersClusterablesComparableRepository) resultToDomain(
	result client.ResultSet,
	row int,
) (clusterables.Comparable, error) {
	id, err := repository.stringColumnValue(
		result,
		groupingsClustersClusterablesCandidatesComparableIdentifierField,
		row,
	)
	if err != nil {
		return nil, err
	}

	identifier, err := uuid.Parse(id)
	if err != nil {
		return nil, clusterables.ErrInvalidClusterableIdentifier
	}

	kind, err := repository.stringColumnValue(
		result,
		groupingsClustersClusterablesCandidatesComparableKindField,
		row,
	)
	if err != nil {
		return nil, err
	}

	vector, err := repository.vectorColumnValue(
		result,
		groupingsClustersClusterablesCandidatesComparableVectorField,
		row,
	)
	if err != nil {
		return nil, err
	}

	return repository.adapter.ToDomain(
		clusterables.ComparableInput{
			Clusterable: clusterables.ClusterableInput{
				Identifier:  identifier,
				ClusterKind: clusterables.Kind(kind),
			},
			Vector: vector,
		},
	)
}

func (repository *groupingsClustersClusterablesComparableRepository) stringColumnValue(
	result client.ResultSet,
	field string,
	row int,
) (string, error) {
	column := result.GetColumn(field)
	if column == nil {
		return "", fmt.Errorf("missing milvus field %s", field)
	}

	value, err := column.Get(row)
	if err != nil {
		return "", err
	}

	switch typed := value.(type) {
	case string:
		return typed, nil
	default:
		return "", fmt.Errorf(
			"invalid milvus field %s type %T",
			field,
			value,
		)
	}
}

func (repository *groupingsClustersClusterablesComparableRepository) vectorColumnValue(
	result client.ResultSet,
	field string,
	row int,
) ([]float32, error) {
	column := result.GetColumn(field)
	if column == nil {
		return nil, fmt.Errorf("missing milvus field %s", field)
	}

	value, err := column.Get(row)
	if err != nil {
		return nil, err
	}

	switch typed := value.(type) {
	case []float32:
		return copyVector(typed), nil

	case entity.FloatVector:
		return copyVector([]float32(typed)), nil

	default:
		return nil, fmt.Errorf(
			"invalid milvus field %s type %T",
			field,
			value,
		)
	}
}

func copyVector(
	vector []float32,
) []float32 {
	out := make([]float32, len(vector))
	copy(out, vector)

	return out
}
