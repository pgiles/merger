package main

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
)

const outFile = "output.csv"

func mainer() {

	// make sure there are only 2 args
	if len(os.Args) != 3 {
		log.Panic("\nUsage: command file1 file2")
	}

	// open the first file
	f1, e := os.Open(os.Args[1])
	if e != nil {
		log.Panic("\nUnable to open first file: ", e)
	}
	defer f1.Close()

	// open second file
	f2, e := os.Open(os.Args[2])
	if e != nil {
		log.Panic("\nUnable to open second file: ", e)
	}
	defer f2.Close()

	// clean up old output
	if _, err := os.Stat(outFile); !errors.Is(err, os.ErrNotExist) {
		// path does not exist
		os.Remove(outFile)
	}

	// create a file writer
	w, e := os.Create(outFile)
	if e != nil {
		log.Panic("\nUnable to create outFile: ", e)
	}
	defer func(w *os.File) {
		e = w.Close()
		if e != nil {
			log.Panic("\nUnable to close outFile: ", e)
		}
	}(w)

	// wrap the file readers with CSV readers
	cr1 := csv.NewReader(f1)
	cr2 := csv.NewReader(f2)

	// wrap the out file writer with a CSV writer
	cw := csv.NewWriter(w)

	// initialize the lines
	line1, b := readline(cr1)
	if !b {
		log.Panic("\nNo CSV lines in file 1.")
	}
	line2, b := readline(cr2)
	if !b {
		log.Panic("\nNo CSV lines in file 2.")
	}

	// copy the files according to similar rules of the merge step in Mergesort
	for {
		if compare(line1, line2) {
			writeline(cw, line1)
			if line1, b = readline(cr1); !b {
				copy(cr2, cw)
				break
			}
		} else {
			writeline(cw, line2)
			if line2, b = readline(cr2); !b {
				copy(cr1, cw)
				break
			}
		}
	}

	cw.Flush()

	// note the files will be closed here, since we defered it above
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

func writeline(w *csv.Writer, line []string) {
	e := w.Write(line)
	if e != nil {
		log.Panic("\nError writing file: ", e)
	}
}

func copy(r *csv.Reader, w *csv.Writer) {
	for line, b := readline(r); b; line, b = readline(r) {
		writeline(w, line)
	}
}

func compare(line1, line2 []string) bool {
	/* here, determine if line1 and line2 are in the correct order (line1 first)
	   if so, return true, otherwise false
	*/
	return true
}
