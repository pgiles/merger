package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
)

const outputFile = "output.csv"

func main() {

	// clean up old output
	if _, err := os.Stat(outputFile); !errors.Is(err, os.ErrNotExist) {
		// path does not exist
		os.Remove(outputFile)
	}

	// create a file writer
	w, e := os.Create(outputFile)
	if e != nil {
		log.Panic("\nUnable to create outFile: ", e)
	}
	defer func(w *os.File) {
		e = w.Close()
		if e != nil {
			log.Panic("\nUnable to close outFile: ", e)
		}
	}(w)

	cw := csv.NewWriter(w)
	csvMerge(cw, os.Args)
}

func csvMerge(w *csv.Writer, files []string) {
	fmt.Println("input files:", files)
	src := openFile(files[1])
	defer src.Close()
	csrc := csv.NewReader(src)
	copy(csrc, w)
	w.Flush()

	for i := 1; i < len(files)-1; i++ {
		fmt.Print(files[1] + " <- " + files[i+1] + "\n")
		src := openFile(files[i+1])
		csvSrc := csv.NewReader(src)

		copy(csvSrc, w)
		src.Close()
		w.Flush()
	}
}

func openFile(file string) *os.File {
	// open the first file
	f, e := os.Open(file)
	if e != nil {
		log.Panic("\nUnable to open first file: ", e)
	}
	//defer f.Close()

	return f
}
