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
	indexMap := make(map[string]int)
	wantSet := make(map[string]bool)
	for i, header := range headers {
		indexMap[header] = i
	}
	for _, w := range want {
		wantSet[w] = true
	}
	var indexes []int
	for w := range wantSet {
		if index, ok := indexMap[w]; ok {
			indexes = append(indexes, index)
		}
	}
	return indexes
}
