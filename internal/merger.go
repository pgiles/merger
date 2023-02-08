package internal

import (
	"encoding/csv"
	"errors"
	"fmt"
	log "golang.org/x/exp/slog"
	"io"
	"os"
)

const outputFile = "merged.csv"

type Merger struct {
	OutputFileName string
}

func (m *Merger) Merge(filenames []string, outputFilename *string) {
	m.OutputFileName = outputFile
	if outputFilename != nil {
		m.OutputFileName = *outputFilename
	}
	w := DeleteAndCreateFile(m.OutputFileName)
	defer closeFile(w)

	cw := csv.NewWriter(w)
	m.AppendCSVFiles(cw, filenames)
}

// AppendCSVFiles appends the files in the array to the outputFile (writer)
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
