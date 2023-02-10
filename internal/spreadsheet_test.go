package internal

import (
	"fmt"
	"testing"
)

// Tests with E2E (End-to-End) prefix are not executed by "go test" for hopefully obvious reasons.

func TestShowHeaders(t *testing.T) {
	//[[Date Category Amount] [first_name last_name ssn]]
	var tests = []struct {
		files []string
		want  [][]string
	}{
		{
			[]string{"../cmd/fixtures/transactions.csv", "../cmd/fixtures/test_info.csv"},
			[][]string{
				{"Transaction Date", "Post Date", "Category", "Amount"},
				{"Field", "Type", "Null", "Key", "Default", "Extra"},
			},
		},
	}
	//t.Run enables running “subtests”, one for each table entry. These are shown separately when executing go test -v.
	for _, tt := range tests {
		testName := fmt.Sprintf("%s", tt.files)
		t.Run(testName, func(t *testing.T) {
			ans := Headers(tt.files)
			for i, w := range tt.want {
				for j, jw := range w {
					if ans[i][j] != jw {
						t.Errorf("got: '%s', want: '%s'", ans[i], tt.want[i])
					}
				}
			}
		})
	}
}
