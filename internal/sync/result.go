package sync

type SyncResult struct {
	Cloned  []string `json:"cloned"`
	Updated []string `json:"updated"`
	Deleted []string `json:"deleted"`
}

func NewSyncResult(output []byte) (SyncResult, error) {
	return SyncResult{}, nil
}

/*

{"message":"Sync completed: querying for latest repositories ...
last synced: 2 days ago
updates found:
  clone  :   0
  update :   1
  delete :   0
cloned: 0 repositories (elapsed time: 250ns)
updated: 1 repositories (elapsed time: 2.3528115s)
deleted: 0 repositories (elapsed time: 166ns)
  \"semhook\" is updated (last updated: 2 days ago)
"}
*/
