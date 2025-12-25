package internal

import (
	"encoding/csv"
	"errors"
	"fmt"
	log "golang.org/x/exp/slog"
	"io"
	"os"
	"strings"
)

const DefaultOutputFile = "merged.csv"
const OutputConfigFileName = "cfg.csv"

type Merger struct {
	OutputFileName string
	GenerateConfig bool
	NegateColumns  []string
}

func (m *Merger) Merge(filenames []string, outputFilename *string) {
	f := m.outputFile(outputFilename)
	defer closeFile(f)

	cw := csv.NewWriter(f)
	m.AppendCSVFiles(cw, filenames)
}

func (m *Merger) CombineCSVFiles(filenames []string, cols []string, outputFilename *string) {
	f := m.outputFile(outputFilename)
	defer closeFile(f)

	cw := csv.NewWriter(f)
	m.combine(cw, filenames, cols)
}

func (m *Merger) combine(w *csv.Writer, files []string, columns []string) {
	//first try at this will be a naive impl:
	// 1. read in records of each input file, write columns with matching headers; load everything into memory
	log.Debug("columns to keep", "columns", columns)
	log.Debug("columns to negate", "negate", m.NegateColumns)

	// Build a set of column names to negate for quick lookup
	negateSet := make(map[string]bool)
	for _, col := range m.NegateColumns {
		negateSet[col] = true
	}

	for _, f := range files {
		reader := csv.NewReader(openFile(f))
		records, _ := reader.ReadAll()
		rows := make([][]string, len(records))
		indexes := ColumnIndexes(records[0], columns)

		for i := 0; i < len(records); i++ {
			var cols []string
			for _, col := range indexes {
				value := records[i][col]
				// Apply negation if this column is in the negate list and it's not the header row
				// Use records[0][col] to get the actual column name from the file's header
				if i > 0 && negateSet[records[0][col]] {
					value = NegateValue(value)
				}
				cols = append(cols, value) //the columns we are using in the output file
			}
			rows[i] = append(rows[i], cols...) //the rows for this file
		}

		err := w.WriteAll(rows)
		if err != nil {
			LogPanic("", err)
		}
		w.Flush()
		fmt.Printf("%v <- %s\n", m.OutputFileName, f)
	}
	m.GenerateConfigFile(columns)

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

func (m *Merger) GenerateConfigFile(header []string) {
	if !m.GenerateConfig {
		return
	}
	enc := csv.NewWriter(DeleteAndCreateFile(OutputConfigFileName)) //Lazy here; client can't choose config file name
	e := enc.Write(header)
	if e != nil {
		LogPanic("Unable to create config file.", e, "file", OutputConfigFileName)
	}
	enc.Flush()
	fmt.Printf("generated %s\n", OutputConfigFileName)
}

func LoadConfigFile(f string) []string {
	reader := csv.NewReader(openFile(f))
	records, _ := reader.ReadAll()
	return records[0]
}

func (m *Merger) outputFile(outputFilename *string) *os.File {
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

// NegateValue converts a negative numeric string to its positive equivalent.
// If the value starts with a minus sign, it removes the minus sign.
// If the value is already positive or non-numeric, it returns the value unchanged.
func NegateValue(value string) string {
	trimmed := strings.TrimSpace(value)
	if strings.HasPrefix(trimmed, "-") {
		return strings.TrimPrefix(trimmed, "-")
	}
	return value
}
