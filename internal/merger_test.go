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
	headers := []string{"first_name", "ssn", "Transaction Date", "Category", "Amount", "Amount", "ssn"}
	w := bytes.NewBufferString("")
	m.combine(csv.NewWriter(w), files, headers)
	//fmt.Print(w)

	approvals.VerifyString(t, w.String())
}

func TestNegateValue(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  string
	}{
		{"negative value", "-50.00", "50.00"},
		{"positive value", "25.50", "25.50"},
		{"zero", "0", "0"},
		{"negative integer", "-100", "100"},
		{"empty string", "", ""},
		{"negative with whitespace", " -25.00", "25.00"},
		{"non-numeric", "text", "text"},
		{"negative non-numeric", "-text", "text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NegateValue(tt.input)
			if got != tt.want {
				t.Errorf("NegateValue(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCombineWithNegate(t *testing.T) {
	m := &Merger{NegateColumns: []string{"Amount"}}
	files := []string{"../cmd/fixtures/negative_test.csv"}
	headers := []string{"Date", "Amount", "Description"}
	w := bytes.NewBufferString("")
	m.combine(csv.NewWriter(w), files, headers)

	expected := `Date,Amount,Description
2024-01-01,50.00,Purchase 1
2024-01-02,25.50,Refund
2024-01-03,100.25,Purchase 2
`
	if w.String() != expected {
		t.Errorf("TestCombineWithNegate got:\n%s\nwant:\n%s", w.String(), expected)
	}
}
