package platforms

import "time"

func parseTime(value string) (time.Time, error) {

	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err == nil {
		return parsed, nil
	}
	return time.Parse(time.RFC3339, value)

}
