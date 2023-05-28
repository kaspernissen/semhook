package repo

import (
	"context"
	"os/exec"
)

func ListRepos(ctx context.Context) ([]ListResult, error) {
	cmd := exec.CommandContext(ctx, "starhook", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return NewListResult(output)
}
