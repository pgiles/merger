package cmd

import (
	"bytes"
	"fmt"
	"github.com/pgiles/merger/internal"
	"os"
	"reflect"
	"testing"
)

func TestMatchSelected(t *testing.T) {
	var tests = []struct {
		headers  [][]string
		selected []string
		want     []string
	}{
		{
			[][]string{
				{"Transaction Date", "Post Date", "Category", "Amount"},
				{"Field", "Type", "Null", "Key", "Default", "Extra"},
			},
			[]string{"0", "3", "4"},
			[]string{"Transaction Date", "Amount", "Default"},
		},
	}

	//t.Run enables running “subtests”, one for each table entry. These are shown separately when executing go test -v.
	for _, tt := range tests {
		testName := fmt.Sprintf("%s", tt.headers)
		t.Run(testName, func(t *testing.T) {
			ans := matchSelected(tt.headers, tt.selected)
			if reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got: '%s', want: '%s'", ans, tt.want)
			}
		})
	}
}

func TestFiles(t *testing.T) {
	var tests = []struct {
		args []string
		want []string
	}{
		{
			[]string{"./fixtures"},
			[]string{"test.csv", "test_info.csv", "transactions.csv"},
		},
	}

	//t.Run enables running “subtests”, one for each table entry. These are shown separately when executing go test -v.
	for _, tt := range tests {
		testName := fmt.Sprintf("%s", tt.args)
		t.Run(testName, func(t *testing.T) {
			ans, err := Files(tt.args)
			if err != nil {
				t.Error(err)
			}
			if reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got: '%s', want: '%s'", ans, tt.want)
			}
		})
	}
}
func TestCSV(t *testing.T) {
	// delete result file
	defer func() {
		_ = os.Remove(internal.DefaultOutputFile)
	}()

	cmd := rootCmd.Root()
	cmd.SetArgs([]string{"csv", "./fixtures"})
	b := bytes.NewBufferString("")
	cmd.SetOut(b)

	err := cmd.Execute()
	if err != nil {
		t.Fail()
	}

}
