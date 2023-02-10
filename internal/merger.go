package internal

import (
	"encoding/csv"
	"errors"
	"fmt"
	log "golang.org/x/exp/slog"
	"io"
	"os"
)

const defaultOutputFile = "merged.csv"

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

func (m *Merger) combine(w *csv.Writer, files []string, cols []string) {
	//first try at this will be a naive impl:
	// 1. read in records of each input file, write columns with matching headers; load everything into memory
	log.Debug("cols to keep", cols)
	for _, f := range files {
		reader := csv.NewReader(openFile(f))
		records, _ := reader.ReadAll()
		filteredRecords := make([][]string, len(records))
		for i := 0; i < len(records); i++ {
			cols := []string{records[i][0], records[i][2]}
			filteredRecords[i] = append(filteredRecords[i], cols...)
		}

		err := w.WriteAll(filteredRecords)
		if err != nil {
			LogPanic("", err)
		}
		w.Flush()
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

func ShowHeaders(files []string) [][]string {
	var r = make([][]string, len(files))
	for i, f := range files {
		line, _ := readline(csv.NewReader(openFile(f)))
		r[i] = line
	}
	return r
}

func outputFile(m *Merger, outputFilename *string) *os.File {
	m.OutputFileName = defaultOutputFile
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

// LogPanic logs to ERROR (would prefer to log as FATAL, but I'm not going to
// create custom levels to do it) and panics. It is a necessary convenience
// method that is here in absence of a log.Panic in golang.org/x/exp/slog. It's
// not ideal.
func LogPanic(msg string, err error, args ...any) {
	//os.Setenv("LOG_SOURCE", "1")
	args = append(args, log.Any("level", "FATAL"))
	log.Error(msg, err, args...)
	panic(fmt.Sprintf("\n%v", err))
}
