package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
)

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query     string          `json:"query"`
	Variables json.RawMessage `json:"variables,omitempty"`
}

// CallGraphQL sends a GraphQL query to the specified endpoint and returns the response
func CallGraphQL(ctx context.Context, endpoint string, query string, variables *string, headers map[string]string) (string, error) {
	// Create HTTP client with reasonable timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Parse variables if provided
	var variablesJSON json.RawMessage
	if variables != nil && *variables != "" {
		if err := json.Unmarshal([]byte(*variables), &variablesJSON); err != nil {
			return "", errors.Wrap(err, "failed to parse variables as JSON")
		}
	}

	// Create GraphQL request
	gqlRequest := GraphQLRequest{
		Query:     query,
		Variables: variablesJSON,
	}

	// Marshal request to JSON
	requestBody, err := json.Marshal(gqlRequest)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal GraphQL request")
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", errors.Wrap(err, "failed to create HTTP request")
	}

	// Set default Content-Type
	req.Header.Set("Content-Type", "application/json")

	// Set custom headers if provided
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute GraphQL request")
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}

	// Check for HTTP error
	if resp.StatusCode >= 400 {
		return "", errors.Newf("GraphQL request failed with status code %d: %s", resp.StatusCode, string(respBody))
	}

	// Pretty-print JSON for better readability in response
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, respBody, "", "  "); err != nil {
		// If we can't pretty-print, return the raw JSON
		return string(respBody), nil
	}

	// JSON を圧縮する
	compressedJSON := bytes.NewBuffer(nil)
	if err := json.Compact(compressedJSON, respBody); err != nil {
		return "", errors.Wrap(err, "failed to compress JSON")
	}

	return compressedJSON.String(), nil
}
