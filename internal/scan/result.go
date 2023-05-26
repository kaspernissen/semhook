package scan

type result struct {
	ruleName string
	repoResults []SemgrepResult
}

type SemgrepResult struct {
	CheckID string   `json:"check_id"`
	Path    string   `json:"path"`
	Start   Position `json:"start"`
	End     Position `json:"end"`
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

func newResult(rule string) result {
	return result{
		ruleName: rule
	}
}

func (r *result) addRepoResult(semgrepOutput string) error { 
	var output SemgrepOutput
	err := json.Unmarshal([]byte(jsonStr), &output)
	if err != nil {
		return err
	}
	r.repoResult = append(r.repoResult, output)
	return nil
}