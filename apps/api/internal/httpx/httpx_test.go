package httpx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestJSON(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   any
		want   map[string]any
	}{
		{
			name:   "object body",
			status: http.StatusOK,
			body:   map[string]any{"id": "abc", "count": float64(3)},
			want:   map[string]any{"id": "abc", "count": float64(3)},
		},
		{
			name:   "empty object",
			status: http.StatusCreated,
			body:   map[string]any{},
			want:   map[string]any{},
		},
		{
			name:   "non-2xx status still encodes body",
			status: http.StatusTeapot,
			body:   map[string]any{"ok": false},
			want:   map[string]any{"ok": false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			JSON(rec, tt.status, tt.body)

			if rec.Code != tt.status {
				t.Errorf("status = %d, want %d", rec.Code, tt.status)
			}

			if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
				t.Errorf("Content-Type = %q, want %q", got, "application/json; charset=utf-8")
			}

			var got map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("decoding body: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("body = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestJSONStruct(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	rec := httptest.NewRecorder()
	in := payload{Name: "ada", Age: 36}

	JSON(rec, http.StatusOK, in)

	var out payload
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decoding body: %v", err)
	}
	if out != in {
		t.Errorf("body = %#v, want %#v", out, in)
	}
}

func TestJSONError(t *testing.T) {
	tests := []struct {
		name   string
		status int
		msg    string
	}{
		{name: "bad request", status: http.StatusBadRequest, msg: "invalid payload"},
		{name: "not found", status: http.StatusNotFound, msg: "missing"},
		{name: "empty message", status: http.StatusInternalServerError, msg: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			JSONError(rec, tt.status, tt.msg)

			if rec.Code != tt.status {
				t.Errorf("status = %d, want %d", rec.Code, tt.status)
			}

			if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
				t.Errorf("Content-Type = %q, want %q", got, "application/json; charset=utf-8")
			}

			var got map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("decoding body: %v", err)
			}
			if len(got) != 1 {
				t.Fatalf("body has %d fields, want exactly 1: %#v", len(got), got)
			}
			if got["error"] != tt.msg {
				t.Errorf("error = %q, want %q", got["error"], tt.msg)
			}
		})
	}
}
