package scan

import (
	"encoding/json"
	"errors"

	"github.com/Microsoft/go-winio/pkg/guid"
)

type Result struct {
	RuleName    string
	RepoResults []SemgrepResult
}

type SemgrepResult struct {
	CheckID string    `json:"check_id"`
	Path    string    `json:"path"`
	Start   Position  `json:"start"`
	End     Position  `json:"end"`
	Extra   ExtraData `json:"extra"`
}

type Position struct {
	Line   int `json:"line"`
	Col    int `json:"col"`
	Offset int `json:"offset"`
}

type Metavar struct {
	AbstractContent string   `json:"abstract_content"`
	Start           Position `json:"start"`
	End             Position `json:"end"`
}

type ExtraData struct {
	Fingerprint string                 `json:"fingerprint"`
	Fix         string                 `json:"fix"`
	IsIgnored   bool                   `json:"is_ignored"`
	Lines       string                 `json:"lines"`
	Message     string                 `json:"message"`
	Metadata    map[string]interface{} `json:"metadata"`
	Metavars    map[string]Metavar     `json:"metavars"`
	Severity    string                 `json:"severity"`
}

type Paths struct {
	Comment string   `json:"_comment"`
	Scanned []string `json:"scanned"`
}

type SemgrepOutput struct {
	Errors  []string        `json:"errors"`
	Paths   Paths           `json:"paths"`
	Results []SemgrepResult `json:"results"`
	Version string          `json:"version"`
}

func NewResult(rule string) Result {
	return Result{
		RuleName: rule,
	}
}

func (r *Result) addRepoResult(semgrepOutput []byte) error {
	var output SemgrepOutput
	err := json.Unmarshal(semgrepOutput, &output)
	if err != nil {
		return err
	}
	if len(output.Results) != 0 {
		r.RepoResults = append(r.RepoResults, output.Results...)
	}
	return nil
}

type ScanID string

func NewScanID() ScanID {
	ID, err := guid.NewV4()
	if err != nil {
		panic("rand failed - seek help!")
	}
	return ScanID(ID.String())
}

func NewScanIDFrom(ID string) (ScanID, error) {
	if ID == "" {
		return "", errors.New("missing scanid")
	}
	return ScanID(ID), nil
}
