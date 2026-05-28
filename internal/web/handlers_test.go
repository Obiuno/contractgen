package web

import (
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateHandler_HappyPath(t *testing.T) {
	// arrange - mock form, create tests cases

	cases := []struct {
		name         string
		fixture      string
		wantStatus   int
		wantContains []string
	}{
		{
			name:         "happy path",
			fixture:      "valid_simple.yml",
			wantStatus:   200,
			wantContains: []string{"CREATE TABLE", "users"},
		},
		{
			name:         "valid with FKs",
			fixture:      "valid_with_fks.yml",
			wantStatus:   200,
			wantContains: []string{"CREATE TABLE", "users", "accounts"},
		},
		{
			name:         "invalid identifier",
			fixture:      "invalid_identifier.yml",
			wantStatus:   400,
			wantContains: []string{"validation failed", "identifier"},
		},
		{
			name:         "invalid yaml syntax",
			fixture:      "invalid_yaml_syntax.yml",
			wantStatus:   400,
			wantContains: []string{"invalid YAML"},
		},
		{
			name:         "invalid fk target",
			fixture:      "invalid_fk_target.yml",
			wantStatus:   400,
			wantContains: []string{"validation failed", "unknown table"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			yamlData, err := os.ReadFile(filepath.Join("testdata", "handler", tc.fixture))
			if err != nil {
				t.Fatalf("read fixture: %v", err)
			}

			form := url.Values{"yaml": {string(yamlData)}}
			req := httptest.NewRequest("POST", "/generate", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// act and capture result
			rec := httptest.NewRecorder()

			GenerateHandler(rec, req)

			// assert - on the response
			if rec.Code != tc.wantStatus {
				t.Errorf("expected status %d, got %d. body: %s", tc.wantStatus, rec.Code, rec.Body.String())
			}
			body := rec.Body.String()

			for _, exp := range tc.wantContains {
				if !strings.Contains(body, exp) {
					t.Errorf("response missing %q, body was: \n%s", exp, body)
				}
			}
		})
	}

}
