/*
Copyright © 2023 Paul Giles <pgilescapone@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"bufio"
	"fmt"
	"github.com/pgiles/merger/internal"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

// csvCmd represents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv",
	Args:  cobra.MinimumNArgs(1),
	Short: "Combine CSV files",
	Long: `Pass file paths or directories as arguments. 

Each file's contents (all rows, including headers) will be appended to 
the file passed before it resulting in a single CVS file named merged.csv
that contains the data from all files.

You can select the columns you'd like to use in the final (merged) result 
by using the interactive mode.
`,

	Example: "csv some/path/file.csv /a/file/to/append/append-me.csv\ncsv . -i",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := Files(args)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		if b, _ := cmd.Flags().GetBool("plan"); b == true {
			headers := internal.Headers(files)
			cmd.Println(prettyPrint(headers))
			return
		} else if b, _ := cmd.Flags().GetBool("interactive"); b == true {
			headers := internal.Headers(files)
			cmd.Println(prettyPrint(headers))
			selected := captureInteractiveInput()

			cols := matchSelected(headers, selected)
			// TODO have backend spit out a config.csv along with combined result
			new(internal.Merger).CombineCSVFiles(files, cols, nil)
			return
		}
		new(internal.Merger).Merge(files, nil)
	},
}

func Files(args []string) ([]string, error) {
	var fileList []string

	for _, a := range args {
		if fi, err := os.Stat(a); err != nil {
			return nil, err
		} else if fi.IsDir() {
			files, err := os.ReadDir(a)
			if err != nil {
				return nil, err
			}
			for _, f := range files {
				if strings.HasSuffix(strings.ToLower(f.Name()), ".csv") {
					fileList = append(fileList, strings.Join([]string{a, f.Name()}, "/"))
				}
			}
		} else {
			if strings.HasSuffix(strings.ToLower(a), ".csv") {
				fileList = append(fileList, a)
			}
		}
	}

	return fileList, nil
}
func matchSelected(headers [][]string, selected []string) []string {
	var tmpArr []string
	// Convert 2D array of each file's headers into a single array since that
	// is how the input is presented (numbered)
	for i := 0; i < len(headers); i++ {
		for _, header := range headers[i] {
			tmpArr = append(tmpArr, header)
		}
	}

	var arr []string
	for _, x := range selected {
		var idx, _ = strconv.Atoi(x)
		arr = append(arr, tmpArr[idx])
	}
	return arr
}
func prettyPrint(headers [][]string) string {
	var s string
	c := 0
	for i := 0; i < len(headers); i++ {
		for _, header := range headers[i] {
			s += fmt.Sprintf("[%d]:'%s' ", c, header)
			c++
		}
		s += "\n"
	}
	return s
}
func captureInteractiveInput() []string {
	// To create dynamic array
	arr := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Press RETURN when finished.")
	for {
		fmt.Print("Enter Text: ")
		// Scans a line from Stdin(Console)
		scanner.Scan()
		// Holds the string that scanned
		text := scanner.Text()
		if len(text) != 0 {
			fmt.Println(text)
			arr = append(arr, text)
		} else {
			break
		}

	}
	// Use collected inputs
	return arr
}

func init() {
	rootCmd.AddCommand(csvCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// csvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	csvCmd.Flags().BoolP("plan", "p", false, "Show the headers for each input file")
	csvCmd.Flags().BoolP("interactive", "i", false, "Pick your columns interactively and store as config for future runs")
	// TODO implement merge using a config file
	csvCmd.Flags().StringP("config", "c", "", "Use a set of headers configured from a previous run")
}
