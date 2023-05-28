package sync

import (
	"context"
	"os/exec"
)

func sync(ctx context.Context) (SyncResult, error) {
	cmd := exec.CommandContext(ctx, "starhook", "sync", "--force=true")
	output, err := cmd.Output()
	if err != nil {
		return SyncResult{}, err
	}

	return NewSyncResult(output)
}
