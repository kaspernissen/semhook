package repo

type ListResult struct {
	Org       string   `json:"org"`
	RepoNames []string `json:"repoNames"`
}

func NewListResult(output []byte) (ListResult, error) {
	result := ListResult{}

	return result, nil
}
