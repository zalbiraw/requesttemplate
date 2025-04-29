package requesttemplate_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zalbiraw/requesttemplate"
)

func TestServeHTTP_TemplatePrependStaticValue(t *testing.T) {
	cfg := requesttemplate.CreateConfig()
	cfg.Template = `{"message": "hello, {{ .user.message }}"}`

	ctx := context.Background()
	var mutatedReq *http.Request
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		mutatedReq = req
	})

	handler, err := requesttemplate.New(ctx, next, cfg, "template-prepend-plugin")
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

	// Assert mutated request if relevant
	if mutatedReq == nil {
		t.Fatalf("next handler was not called; recorder code: %d, body: %s", recorder.Code, recorder.Body.String())
	}

	// Assert on the response body
	gotBody, _ := io.ReadAll(mutatedReq.Body)
	var got map[string]any
	if err := json.Unmarshal(gotBody, &got); err != nil {
		t.Fatalf("failed to unmarshal body: %v. Raw body: %s", err, string(gotBody))
	}
	msg, ok := got["message"].(string)
	if !ok {
		t.Fatalf("expected message to be a string")
	}
	if msg != "hello, world" {
		t.Errorf("expected message to be 'hello, world', got '%s'", msg)
	}
}
