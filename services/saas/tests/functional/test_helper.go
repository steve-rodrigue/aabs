package functional

import "github.com/google/uuid"

func mustParseUUID(
	value string,
) uuid.UUID {
	id, err := uuid.Parse(value)
	if err != nil {
		panic(err)
	}

	return id
}
