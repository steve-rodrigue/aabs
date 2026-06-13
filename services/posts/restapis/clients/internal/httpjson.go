package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type apiError struct {
	Error string `json:"error"`
}

func DoJSON(
	ctx context.Context,
	client *http.Client,
	method string,
	url string,
	input any,
	output any,
) error {
	var body *bytes.Reader

	if input != nil {
		payload, err := json.Marshal(input)
		if err != nil {
			return err
		}

		body = bytes.NewReader(payload)
	} else {
		body = bytes.NewReader(nil)
	}

	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}

	if input != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		var apiErr apiError
		_ = json.NewDecoder(response.Body).Decode(&apiErr)

		if apiErr.Error != "" {
			return fmt.Errorf(apiErr.Error)
		}

		return fmt.Errorf("http error: %d", response.StatusCode)
	}

	if output == nil || response.StatusCode == http.StatusNoContent {
		return nil
	}

	return json.NewDecoder(response.Body).Decode(output)
}
