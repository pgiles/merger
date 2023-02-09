# Merger

Tool for appending CSV files

TODO:
Tool for combining CSV files with different headers (first line) into a single CSV file.

## Commands
To learn how to use this tool, you can rely upon the feedback and examples offered by the CLI.
```
$ go build
$ merger                          
Merger is a tool for combining CSV files with different
headers (first line) into a single CSV file.

Usage:
  merger [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  csv         Combine CSV files
  help        Help about any command

Flags:
  -h, --help      help for merger
  -v, --version   version for merger

Use "merger [command] --help" for more information about a command.
```

## Logging
Logging output has the following configuration options.

| Environment Variable | Options                      |
|----------------------|------------------------------|
| LOG_LEVEL            | debug, info, error (default) |
| LOG_FORMAT           | json, text (default)         |
| LOG_FILE             | file name                    |
| LOG_SOURCE           | any, not present (default)   |

Setting LOG_SOURCE to any value adds a ("source",
"file:line") attribute to the output indicating the source code position of
the log statement.

## Usage
This is a [go app](https://go.dev/doc/effective_go).  All the `go build`, `go run`, `go test`, etc. work.

## License
[MIT License](LICENSE)