package internal

import (
	"bytes"
	"encoding/csv"
	"testing"

	approvals "github.com/approvals/go-approval-tests"
)

// Tests with E2E (End-to-End) prefix are not executed by "go test" for hopefully obvious reasons.

func init() {
	approvals.UseFolder("fixtures")
}

func TestCombine(t *testing.T) {
	m := new(Merger)
	files := []string{"../cmd/fixtures/test.csv", "../cmd/fixtures/transactions.CSV"}
	headers := []string{"first_name", "ssn", "Transaction Date", "Category", "Amount", "ssn"}
	w := bytes.NewBufferString("")
	m.combine(csv.NewWriter(w), files, headers)

	expected := `expected output`
	actual := w.String()
	if actual != expected {
		t.Errorf("TestCombine failed. Expected: %s, Actual: %s", expected, actual)
	}
}