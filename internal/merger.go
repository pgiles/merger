package internal

import (
	"encoding/csv"
	"errors"
	"fmt"
	log "golang.org/x/exp/slog"
	"io"
	"os"
)

const DefaultOutputFile = "merged.csv"

type Merger struct {
	OutputFileName string
}

func (m *Merger) Merge(filenames []string, outputFilename *string) {
	f := outputFile(m, outputFilename)
	defer closeFile(f)

	cw := csv.NewWriter(f)
	m.AppendCSVFiles(cw, filenames)
}

func (m *Merger) CombineCSVFiles(filenames []string, cols []string, outputFilename *string) {
	f := outputFile(m, outputFilename)
	defer closeFile(f)

	cw := csv.NewWriter(f)
	m.combine(cw, filenames, cols)
}

func (m *Merger) combine(w *csv.Writer, files []string, columns []string) {
	//first try at this will be a naive impl:
 func (m *Merger) combine(w *csv.Writer, files []string, columns []string) {
 	log.Debug("columns to keep", "columns", columns)
 	headerMap := make(map[string]int)
 	for _, f := range files {
 		reader := csv.NewReader(openFile(f))
 		records, _ := reader.ReadAll()
 		rows := make([][]string, len(records))
 		indexes := ColumnIndexes(records[0], columns)
 		for i, col := range indexes {
 			headerMap[records[0][col]] = i
 		}
 		for i := 0; i < len(records); i++ {
 			var cols []string
 			for _, col := range indexes {
 				cols = append(cols, records[i][col]) //the columns we are using in the output file
 			}
 			for j := 0; j < len(cols); j++ {
 				if j < len(rows[i]) {
 					rows[i][j] += cols[j]
 				} else {
 					rows[i] = append(rows[i], cols[j])
 				}
 			}
 		}
 
 		err := w.WriteAll(rows)
 		if err != nil {
 			LogPanic("", err)
 		}
 		w.Flush()
 		fmt.Printf("%v <- %s\n", m.OutputFileName, f)
 	}
 }

// AppendCSVFiles appends the files in the array to the output file (writer)
func (m *Merger) AppendCSVFiles(w *csv.Writer, files []string) {
	log.Debug("input files", "files", files)
	for i := 0; i < len(files); i++ {
		src := openFile(files[i])
		csvSrc := csv.NewReader(src)

		copyTo(csvSrc, w)
		closeFile(src)
		w.Flush()
		fmt.Printf("%v <- %s\n", m.OutputFileName, files[i])
	}
}

// DeleteAndCreateFile returned file is a resource that must be closed after the program is finished with it
func DeleteAndCreateFile(f string) *os.File {
	// delete old file
	if _, err := os.Stat(f); !errors.Is(err, os.ErrNotExist) {
		// path does not exist
		err := os.Remove(f)
		if err != nil {
			LogPanic("Unable to delete output file.", err, "file", f)
		}
	}

	// create a file writer
	w, e := os.Create(f)
	if e != nil {
		LogPanic("Unable to create output file.", e, "file", f)
	}

	return w
}

func outputFile(m *Merger, outputFilename *string) *os.File {
	m.OutputFileName = DefaultOutputFile
	if outputFilename != nil {
		m.OutputFileName = *outputFilename
	}

	return DeleteAndCreateFile(m.OutputFileName)
}

func closeFile(src *os.File) {
	err := src.Close()
	if err != nil {
		LogPanic("Unable to open file.", err)
	}
}

func openFile(file string) *os.File {
	// open the file
	f, e := os.Open(file)
	if e != nil {
		LogPanic("Unable to open file.", e, "file", file)
	}

	return f
}
func readline(r *csv.Reader) ([]string, bool) {
	line, e := r.Read()
	if e != nil {
		if e == io.EOF {
			return nil, false
		}
		LogPanic("Error reading file.", e, "file", r)
	}
	return line, true
}

func writeLine(w *csv.Writer, line []string) {
	e := w.Write(line)
	if e != nil {
		LogPanic("Error writing file.", e, "file", w)
	}
}

func copyTo(r *csv.Reader, w *csv.Writer) {
	for line, b := readline(r); b; line, b = readline(r) {
		writeLine(w, line)
	}
}