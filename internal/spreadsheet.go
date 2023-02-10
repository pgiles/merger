package internal

import "encoding/csv"

func Headers(files []string) [][]string {
	var r = make([][]string, len(files))
	for i, f := range files {
		line, _ := readline(csv.NewReader(openFile(f)))
		r[i] = line
	}
	return r
}

// ColumnIndexes returns the matching column index of a column position
func ColumnIndexes(headers []string, want []string) []int {
	var indexes []int
	for i := range headers {
		for j := range want {
			if headers[i] == want[j] {
				indexes = append(indexes, i)
				break
			}
		}
	}
	return indexes
}
