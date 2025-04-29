// Package requesttemplate a request template plugin.
package requesttemplate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/itchyny/gojq"
)

// Config the plugin configuration.
type Config struct {
	Commands []string `json:"commands,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Commands: make([]string, 0),
	}
}

// RequestTemplate a Request Template plugin.
type RequestTemplate struct {
	next     http.Handler
	commands []string
}

// New created a new RequestTemplate plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Commands) == 0 {
		return nil, fmt.Errorf("commands cannot be empty")
	}

	return &RequestTemplate{
		commands: config.Commands,
		next:     next,
	}, nil
}

func (a *RequestTemplate) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Read the request body
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
	}
	_ = req.Body.Close()

	if len(bodyBytes) == 0 {
		a.next.ServeHTTP(rw, req)
		return
	}
	var jsonData interface{}
	if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
		http.Error(rw, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	for _, command := range a.commands {
		query, err := gojq.Parse(command)
		if err != nil {
			http.Error(rw, "Invalid jq filter: "+err.Error(), http.StatusBadRequest)
			return
		}
		iter := query.RunWithContext(req.Context(), jsonData)
		// Only take the first result
		v, ok := iter.Next()
		if !ok {
			http.Error(rw, "jq produced no output", http.StatusBadRequest)
			return
		}
		if err, ok := v.(error); ok {
			http.Error(rw, "jq error: "+err.Error(), http.StatusBadRequest)
			return
		}
		jsonData = v
	}
	// Marshal the result back to JSON
	bodyBytes, err := json.Marshal(jsonData)
	if err != nil {
		http.Error(rw, "Failed to marshal JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := rw.Write(bodyBytes); err != nil {
		// Optionally log the error, but cannot send another HTTP error at this point
	}
	return
}
