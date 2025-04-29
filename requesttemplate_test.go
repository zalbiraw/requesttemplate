package requesttemplate

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP_JQPrependStaticValue(t *testing.T) {
	// This jq command prepends "hello, " to the nested .user.message string
	jqCmd := `.user.message |= "hello, " + .`

	cfg := CreateConfig()
	cfg.Commands = []string{jqCmd}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// No-op: nothing needed for the next handler in this test
	})

	handler, err := New(ctx, next, cfg, "jq-prepend-plugin")
	if err != nil {
		t.Fatal(err)
	}

	// Input JSON body
	input := `{"user": {"message": "world"}}`
	recorder := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost", bytes.NewBufferString(input))
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(recorder, req)

	gotBody := recorder.Body.Bytes()

	// Expect the body to have the prepended value
	var got map[string]any
	if err := json.Unmarshal(gotBody, &got); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}
	userMap, ok := got["user"].(map[string]any)
	if !ok {
		t.Fatalf("expected user to be a map[string]any")
	}
	msg, ok := userMap["message"].(string)
	if !ok {
		t.Fatalf("expected message to be a string")
	}
	if msg != "hello, world" {
		t.Errorf("expected message to be 'hello, world', got '%s'", msg)
	}
}
