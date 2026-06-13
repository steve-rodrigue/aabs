package servers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/steve-rodrigue/aabs/services/posts/restapis/applications"
)

type api struct {
	application applications.Application
}

func decodeJSON(
	request *http.Request,
	out any,
) error {
	defer request.Body.Close()

	return json.NewDecoder(request.Body).Decode(out)
}

func respondJSON(
	writer http.ResponseWriter,
	status int,
	value any,
) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	if value == nil {
		return
	}

	_ = json.NewEncoder(writer).Encode(value)
}

func respondError(
	writer http.ResponseWriter,
	status int,
	err error,
) {
	respondJSON(
		writer,
		status,
		map[string]string{
			"error": err.Error(),
		},
	)
}

func routeUUID(
	request *http.Request,
	name string,
) (uuid.UUID, error) {
	return uuid.Parse(mux.Vars(request)[name])
}

func queryInt(
	request *http.Request,
	name string,
	defaultValue int,
) int {
	value := request.URL.Query().Get(name)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

func queryUUID(
	request *http.Request,
	name string,
) uuid.UUID {
	value := request.URL.Query().Get(name)
	if value == "" {
		return uuid.Nil
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil
	}

	return parsed
}
