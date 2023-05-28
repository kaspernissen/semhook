package repo

import "testing"

func TestRepoListOneOrg(t *testing.T) {
	output := `1 hoeg/lars-lars
	2 hoeg/config
	3 hoeg/ql-controller
	4 hoeg/psh
	5 hoeg/gh-insights
	6 hoeg/app-runners
	7 hoeg/intro-to-semgrep
	8 hoeg/codeql-uboot
	9 hoeg/aws-bbb
   10 hoeg/exitnode`

	result, err := NewListResult([]byte(output))
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 org, got %d", len(result))
	}
	if result[0].Org != "hoeg" {
		t.Fatalf("expected hoeg, got: %s", result[0].Org)
	}
	if len(result[0].RepoNames) != 10 {
		t.Fatalf("expected 10 repos, got %d", len(result[0].RepoNames))
	}
}

func TestRepoListTwoOrgs(t *testing.T) {
	output := `1 hoeg/lars-lars
	2 hoeg/config
	3 hello/some-repo
	4 hoeg/psh
	5 hoeg/gh-insights
	6 hoeg/app-runners
	7 hello/another-repo
	8 hoeg/codeql-uboot
	9 hoeg/aws-bbb
   10 hoeg/exitnode`

	result, err := NewListResult([]byte(output))
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 org, got %d", len(result))
	}
}
