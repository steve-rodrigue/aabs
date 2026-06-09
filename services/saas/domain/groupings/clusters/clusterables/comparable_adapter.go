package clusterables

type comparableAdapter struct {
	clusterableAdapter Adapter
}

func (app *comparableAdapter) ToDomain(
	input ComparableInput,
) (Comparable, error) {
	clusterable, err := app.clusterableAdapter.ToDomain(input.Clusterable)
	if err != nil {
		return nil, err
	}

	if len(input.Vector) == 0 {
		return nil, ErrInvalidComparableVector
	}

	vector := make([]float32, len(input.Vector))
	copy(vector, input.Vector)

	return &comparable{
		clusterable: clusterable,
		vector:      vector,
	}, nil
}
