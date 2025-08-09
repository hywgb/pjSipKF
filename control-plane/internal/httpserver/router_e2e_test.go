package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/hywgb/pjSipKF/control-plane/internal/config"
)

type respMap map[string]any

func TestEndToEnd_RegisterLookupSession(t *testing.T) {
	cfg := config.Load()
	logger, _ := zap.NewDevelopment()
	mux := chi.NewRouter()
	mux.Use(WithMiddlewares(logger))
	RegisterRoutes(mux, logger, cfg)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	// register
	body := []byte(`{"user":"alice","contact":"sip:alice@test","expires":5}`)
	res, err := http.Post(ts.URL+"/v1/register", "application/json", bytes.NewReader(body))
	if err != nil { t.Fatal(err) }
	if res.StatusCode != 200 { t.Fatalf("register status=%d", res.StatusCode) }

	// lookup
	res, err = http.Get(ts.URL+"/v1/lookup?user=alice")
	if err != nil { t.Fatal(err) }
	if res.StatusCode != 200 { t.Fatalf("lookup status=%d", res.StatusCode) }
	var out respMap
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil { t.Fatal(err) }
	cts, _ := out["contacts"].([]any)
	if len(cts) != 1 { t.Fatalf("expected 1 contact, got %v", cts) }

	// create session
	req := map[string]any{"sdp_offer": "v=0\n", "metadata": map[string]string{"caller":"alice"}}
	buf, _ := json.Marshal(req)
	res, err = http.Post(ts.URL+"/v1/sessions", "application/json", bytes.NewReader(buf))
	if err != nil { t.Fatal(err) }
	if res.StatusCode != 200 { t.Fatalf("sessions status=%d", res.StatusCode) }
	out = respMap{}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil { t.Fatal(err) }
	if out["session_id"] == "" { t.Fatalf("missing session_id") }
	if out["sdp_answer"] == "" { t.Fatalf("missing sdp_answer") }
}