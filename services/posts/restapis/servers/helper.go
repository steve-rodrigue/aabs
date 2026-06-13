package servers

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func routeValue(
	request *http.Request,
	name string,
) string {
	return mux.Vars(request)[name]
}

func queryUUIDs(request *http.Request, key string) []uuid.UUID {
	values := request.URL.Query()[key]
	out := []uuid.UUID{}

	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			id, err := uuid.Parse(part)
			if err == nil {
				out = append(out, id)
			}
		}
	}

	return out
}
