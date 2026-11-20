package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/AmadeusITGroup/gh-app-auth/pkg/config"
)

type sample struct {
	line    int
	body    string
	invalid bool
}

const unknownFieldFixture = `version: "1.0"
pats:
  - name: Example
    token_source: keyring
    patterns:
      - github.com/example-org/
`

func fenceStart(trimmed string, number int) *sample {
	info := strings.Fields(trimmed[3:])
	if len(info) == 0 || (info[0] != "yaml" && info[0] != "yml") {
		return &sample{line: -1}
	}
	return &sample{line: number + 1, invalid: len(info) > 1 && info[1] == "invalid"}
}

func extractSamples(markdown string) ([]sample, error) {
	var samples []sample
	var current *sample
	for number, line := range strings.Split(markdown, "\n") {
		trimmed := strings.TrimSpace(line)
		isFence := strings.HasPrefix(trimmed, "```")
		switch {
		case current == nil && isFence:
			current = fenceStart(trimmed, number)
		case current != nil && isFence:
			if current.line > 0 {
				samples = append(samples, *current)
			}
			current = nil
		case current != nil:
			current.body += line + "\n"
		}
	}
	if err := unterminated(current); err != nil {
		return nil, err
	}
	return samples, nil
}

func unterminated(current *sample) error {
	if current == nil {
		return nil
	}
	if current.line > 0 {
		return fmt.Errorf("line %d: unterminated yaml code block", current.line)
	}
	return errors.New("unterminated code block")
}

func validate(body string) (*config.Config, error) {
	decoder := yaml.NewDecoder(bytes.NewReader([]byte(body)))
	decoder.KnownFields(true)
	var cfg config.Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return nil, errors.New("decode: sample contains more than one YAML document")
		}
		return nil, fmt.Errorf("decode: trailing content after the first document: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	return &cfg, nil
}

type tally struct {
	appOnly, patOnly, mixed int
}

func (t *tally) count(cfg *config.Config) {
	switch {
	case len(cfg.GitHubApps) > 0 && len(cfg.PATs) > 0:
		t.mixed++
	case len(cfg.GitHubApps) > 0:
		t.appOnly++
	default:
		t.patOnly++
	}
}

func checkSamples(samples []sample) (tally, []string) {
	var counts tally
	var problems []string
	for _, s := range samples {
		cfg, err := validate(s.body)
		switch {
		case s.invalid && err == nil:
			problems = append(problems, fmt.Sprintf("line %d: sample marked invalid was accepted", s.line))
		case !s.invalid && err != nil:
			problems = append(problems, fmt.Sprintf("line %d: %v", s.line, err))
		case !s.invalid:
			counts.count(cfg)
		}
	}
	if counts.appOnly == 0 || counts.patOnly == 0 || counts.mixed == 0 {
		problems = append(problems, fmt.Sprintf("need valid App-only, PAT-only and mixed samples, found %+v", counts))
	}
	return counts, problems
}

func checkUnknownFieldFixture() []string {
	_, err := validate(unknownFieldFixture)
	if err == nil {
		return []string{"unknown token_source field fixture was accepted"}
	}
	if !strings.Contains(err.Error(), "token_source") {
		return []string{fmt.Sprintf("unknown field fixture failed for another reason: %v", err)}
	}
	return nil
}

func check(markdown string) (tally, []string) {
	samples, err := extractSamples(markdown)
	if err != nil {
		return tally{}, []string{err.Error()}
	}
	if len(samples) == 0 {
		return tally{}, []string{"no configuration samples found"}
	}
	counts, problems := checkSamples(samples)
	return counts, append(problems, checkUnknownFieldFixture()...)
}

const configurationReference = "docs/reference/configuration.md"

func run() error {
	content, err := os.ReadFile(configurationReference)
	if err != nil {
		return fmt.Errorf("read %s (run from the repository root): %w", configurationReference, err)
	}
	counts, problems := check(string(content))
	if len(problems) > 0 {
		return fmt.Errorf("configuration samples are invalid:\n- %s", strings.Join(problems, "\n- "))
	}
	fmt.Printf("Validated configuration samples: %d App, %d PAT, %d mixed\n", counts.appOnly, counts.patOnly, counts.mixed)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
