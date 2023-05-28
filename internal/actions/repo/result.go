package repo

import (
	"fmt"
	"strings"
)

type ListResult struct {
	Org       string   `json:"org"`
	RepoNames []string `json:"repoNames"`
}

func NewListResult(output []byte) ([]ListResult, error) {
	lines := strings.Split(string(output), "\n")

	orgMap := make(map[string][]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "==> ") {
			break
		}
		if line != "" {
			parts := strings.SplitN(line, "/", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("list output in unexpected format: %s", line)
			}

			firstPart := strings.Split(parts[0], " ")
			org := firstPart[1]
			repoName := parts[1]
			if _, ok := orgMap[org]; !ok {
				orgMap[org] = make([]string, 0)
			}

			orgMap[org] = append(orgMap[org], repoName)
		}
	}

	results := make([]ListResult, 0)
	for org, repoNames := range orgMap {
		result := ListResult{
			Org:       org,
			RepoNames: repoNames,
		}
		results = append(results, result)
	}

	return results, nil
}
