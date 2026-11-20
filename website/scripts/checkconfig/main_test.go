package main

import (
	"strings"
	"testing"
)

const validDocument = "```yaml\nversion: \"1.0\"\ngithub_apps:\n  - name: App\n    client_id: Iv1.example\n    private_key_source: keyring\n    patterns:\n      - github.com/example-org/\n```\n\n" +
	"```yaml\nversion: \"1.0\"\npats:\n  - name: PAT\n    private_key_source: keyring\n    patterns:\n      - github.com/example-org/\n```\n\n" +
	"```yaml\nversion: \"1.0\"\ngithub_apps:\n  - name: App\n    app_id: 1\n    private_key_source: keyring\n    patterns:\n      - github.com/a/\npats:\n  - name: PAT\n    patterns:\n      - github.com/b/\n```\n"

func TestExtractSamples(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     int
		wantErr  string
	}{
		{"no blocks", "text", 0, ""},
		{"fragment without version is still a sample", "```yaml\nname: x\n```\n", 1, ""},
		{"config block", "```yaml\nversion: \"1\"\n```\n", 1, ""},
		{"other languages ignored", "```json\nversion: 1\n```\n", 0, ""},
		{"text fences ignored", "```text\nversion: 1\n```\n", 0, ""},
		{"three blocks", validDocument, 3, ""},
		{"unterminated yaml block", "```yaml\nversion: \"1\"\n", 0, "unterminated yaml"},
		{"unterminated other block", "```text\nx\n", 0, "unterminated"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractSamples(tt.markdown)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("extractSamples() error = %v, want containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("extractSamples() error = %v", err)
			}
			if len(got) != tt.want {
				t.Fatalf("extractSamples() = %d samples, want %d", len(got), tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr string
	}{
		{"valid app", "version: \"1.0\"\ngithub_apps:\n  - name: A\n    client_id: Iv1.x\n    private_key_source: keyring\n    patterns: [github.com/o/]\n", ""},
		{"unknown field", unknownFieldFixture, "token_source"},
		{"filesystem source without path", "version: \"1.0\"\ngithub_apps:\n  - name: A\n    app_id: 1\n    private_key_source: filesystem\n    patterns: [github.com/o/]\n", "private_key_path"},
		{"no credentials", "version: \"1.0\"\n", "at least one"},
		{"missing version", "pats:\n  - name: P\n    patterns: [x]\n", "version"},
		{"second document", "version: \"1.0\"\npats:\n  - name: P\n    patterns: [x]\n---\nversion: \"1.0\"\n", "more than one YAML document"},
		{"trailing garbage", "version: \"1.0\"\npats:\n  - name: P\n    patterns: [x]\n---\n: : :\n", "trailing content"},
		{"empty", "", "decode"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validate(tt.body)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validate() error = %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("validate() error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	if _, problems := check(validDocument); len(problems) != 0 {
		t.Fatalf("check(valid) problems = %v", problems)
	}
	if _, problems := check("nothing here"); len(problems) == 0 {
		t.Fatal("check(empty) should fail when no samples exist")
	}
	invalidMarked := validDocument + "\n```yaml invalid\nversion: \"1.0\"\ngithub_apps:\n  - name: A\n    app_id: 1\n    private_key_source: keyring\n    patterns: [x]\n```\n"
	if _, problems := check(invalidMarked); len(problems) == 0 {
		t.Fatal("check should flag a sample marked invalid that validates")
	}
	missingVersion := validDocument + "\n```yaml\npats:\n  - name: P\n    patterns: [x]\n```\n"
	if _, problems := check(missingVersion); len(problems) == 0 {
		t.Fatal("check should flag a yaml block without a version")
	}
	if _, problems := check(validDocument + "\n```yaml\nversion: \"1.0\"\n"); len(problems) == 0 {
		t.Fatal("check should flag an unterminated yaml block")
	}
	rejected := validDocument + "\n```yaml invalid\nversion: \"1.0\"\npats:\n  - name: P\n    token_source: keyring\n    patterns: [x]\n```\n"
	if _, problems := check(rejected); len(problems) != 0 {
		t.Fatalf("check(rejected) problems = %v", problems)
	}
}
