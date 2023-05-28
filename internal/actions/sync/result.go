package sync

import "regexp"

type SyncResult struct {
	Cloned  []string `json:"cloned"`
	Updated []string `json:"updated"`
	Deleted []string `json:"deleted"`
}

func NewSyncResult(output []byte) (SyncResult, error) {
	result := SyncResult{}

	regexPattern := `(.*) is (updated|deleted|cloned)`
	regex := regexp.MustCompile(regexPattern)

	matches := regex.FindAllStringSubmatch(string(output), -1)
	for _, match := range matches {
		repoName := match[1]
		action := match[2]

		switch action {
		case "cloned":
			result.Cloned = append(result.Cloned, repoName)
		case "updated":
			result.Updated = append(result.Updated, repoName)
		case "deleted":
			result.Deleted = append(result.Deleted, repoName)
		}
	}

	return result, nil
}
