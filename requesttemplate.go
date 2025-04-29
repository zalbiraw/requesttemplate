// Package requesttemplate a request template plugin.
package requesttemplate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
)

// Config the plugin configuration.
type Config struct {
	Template string `json:"template,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Template: "",
	}
}

// RequestTemplate a Request Template plugin.
type RequestTemplate struct {
	next     http.Handler
	template string
}

// New created a new RequestTemplate plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.Template == "" {
		return nil, fmt.Errorf("template cannot be empty")
	}

	return &RequestTemplate{
		template: config.Template,
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
	var jsonData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &jsonData); err != nil {
		http.Error(rw, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.New("request").Parse(a.template)
	if err != nil {
		http.Error(rw, "Invalid template: "+err.Error(), http.StatusBadRequest)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, jsonData)
	if err != nil {
		http.Error(rw, "Template execution error: "+err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	_, err = rw.Write(buf.Bytes())
	if err != nil {
		return
	}
}
