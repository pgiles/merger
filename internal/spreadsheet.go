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
	wantMap := make(map[string]bool)
	for i, header := range headers {
		indexMap[header] = i
	}
	var indexes []int
	added := make(map[int]bool)
	for _, w := range want {
		if wantMap[w] {
			continue
		}
		wantMap[w] = true
		if index, ok := indexMap[w]; ok && !added[index] {
			indexes = append(indexes, index)
			added[index] = true
		}
	}
	return indexes
}
