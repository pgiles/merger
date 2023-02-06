package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const outputFile = "merged.csv"

func main() {
	w := DeleteAndCreateFile(outputFile)
	defer closeFile(w)

	cw := csv.NewWriter(w)
	AppendCSVFiles(cw, os.Args[1:])
}

// AppendCSVFiles appends the files in the array to the outputFile (writer)
func AppendCSVFiles(w *csv.Writer, files []string) {
	fmt.Println("input files:", files)

	for i := 0; i < len(files); i++ {
		fmt.Printf("%v <- %s\n", outputFile, files[i])
		src := openFile(files[i])
		csvSrc := csv.NewReader(src)

		copyTo(csvSrc, w)
		closeFile(src)
		w.Flush()
	}
}

// DeleteAndCreateFile returned file is a resource that must be closed after the program is finished with it
func DeleteAndCreateFile(f string) *os.File {
	// delete old file
	if _, err := os.Stat(f); !errors.Is(err, os.ErrNotExist) {
		// path does not exist
		err := os.Remove(f)
		if err != nil {
			log.Panic(err)
		}
	}

	// create a file writer
	w, e := os.Create(f)
	if e != nil {
		log.Panic("\nUnable to create output File: ", e)
	}

	return w
}

func closeFile(src *os.File) {
	err := src.Close()
	if err != nil {
		log.Panic(err)
	}
}

func openFile(file string) *os.File {
	// open the file
	f, e := os.Open(file)
	if e != nil {
		log.Panic("\nUnable to open file: ", e)
	}

	return f
}
func readline(r *csv.Reader) ([]string, bool) {
	line, e := r.Read()
	if e != nil {
		if e == io.EOF {
			return nil, false
		}
		log.Panic("\nError reading file: ", e)
	}
	return line, true
}

func writeLine(w *csv.Writer, line []string) {
	e := w.Write(line)
	if e != nil {
		log.Panic("\nError writing file: ", e)
	}
}

func copyTo(r *csv.Reader, w *csv.Writer) {
	for line, b := readline(r); b; line, b = readline(r) {
		writeLine(w, line)
	}
}
