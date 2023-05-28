package scan

import "testing"

func TestNothingToScan(t *testing.T) {
	input := `{"errors": [], "paths": {"_comment": "<add --verbose for a list of skipped paths>", "scanned": []}, "results": [], "version": "1.22.0"}`

	r := result{}
	err := r.addRepoResult([]byte(input))
	if err != nil {
		t.Fatalf("expected to unmarshal result: %v", err)
	}

	if len(r.repoResults) != 0 {
		t.Fatal("expected no result")
	}
}

func TestWithResult(t *testing.T) {
	input := `{"errors": [], "paths": {"_comment": "<add --verbose for a list of skipped paths>", "scanned": ["/Users/hoeg/Documents/security-research/hoeg.github/terraperm/app/terraperm/main.go", "/Users/hoeg/Documents/security-research/hoeg.github/terraperm/module/policy/parser.go", "/Users/hoeg/Documents/security-research/hoeg.github/terraperm/module/policy/printer.go", "/Users/hoeg/Documents/security-research/hoeg.github/terraperm/module/policy/service/s3.go", "/Users/hoeg/Documents/security-research/hoeg.github/terraperm/module/policy/statement.go", "/Users/hoeg/Documents/security-research/hoeg.github/terraperm/module/terraform/cmd.go", "/Users/hoeg/Documents/security-research/hoeg.github/terraperm/module/terraform/trace.go"]}, "results": [{"check_id": "var.folders.p6.h26fyb256tlfxlc2jlm_c5s40000gn.T.temp331558593.go-exec-command", "end": {"col": 43, "line": 14, "offset": 175}, "extra": {"engine_kind": "OSS", "fingerprint": "5c5b5bc77f388f599c3af170b0b42a44592360d2581244881f306fe82985fd8f46e3b96928296610b2dc65a2f0216d31abd4e3bd73833410b291b083efa30390_0", "is_ignored": false, "lines": "\tcmd := exec.Command(\"which\", \"terraform\")", "message": "Found exec command", "metadata": {}, "metavars": {}, "severity": "ERROR"}, "path": "/Users/hoeg/Documents/security-research/hoeg.github/terraperm/module/terraform/cmd.go", "start": {"col": 9, "line": 14, "offset": 141}}, {"check_id": "var.folders.p6.h26fyb256tlfxlc2jlm_c5s40000gn.T.temp331558593.go-exec-command", "end": {"col": 45, "line": 40, "offset": 764}, "extra": {"engine_kind": "OSS", "fingerprint": "5c5b5bc77f388f599c3af170b0b42a44592360d2581244881f306fe82985fd8f46e3b96928296610b2dc65a2f0216d31abd4e3bd73833410b291b083efa30390_1", "is_ignored": false, "lines": "\tcmd := exec.Command(\"terraform\", action...)", "message": "Found exec command", "metadata": {}, "metavars": {}, "severity": "ERROR"}, "path": "/Users/hoeg/Documents/security-research/hoeg.github/terraperm/module/terraform/cmd.go", "start": {"col": 9, "line": 40, "offset": 728}}], "version": "1.22.0"}`

	r := result{}
	err := r.addRepoResult([]byte(input))
	if err != nil {
		t.Fatalf("expected to unmarshal result: %v", err)
	}

	if len(r.repoResults) == 0 {
		t.Fatal("expected results")
	}
}

func TestWithNoResults(t *testing.T) {
	input := `{"errors": [], "paths": {"_comment": "<add --verbose for a list of skipped paths>", "scanned": ["/Users/hoeg/Documents/security-research/hoeg.github/vault-plugin-secrets-dockerhub/backend.go", "/Users/hoeg/Documents/security-research/hoeg.github/vault-plugin-secrets-dockerhub/client.go", "/Users/hoeg/Documents/security-research/hoeg.github/vault-plugin-secrets-dockerhub/cmd/dockerhub/main.go", "/Users/hoeg/Documents/security-research/hoeg.github/vault-plugin-secrets-dockerhub/config.go", "/Users/hoeg/Documents/security-research/hoeg.github/vault-plugin-secrets-dockerhub/data.go", "/Users/hoeg/Documents/security-research/hoeg.github/vault-plugin-secrets-dockerhub/token.go"]}, "results": [], "version": "1.22.0"}`

	r := result{}
	err := r.addRepoResult([]byte(input))
	if err != nil {
		t.Fatalf("expected to unmarshal result: %v", err)
	}

	if len(r.repoResults) != 0 {
		t.Fatal("expected no result")
	}
}
