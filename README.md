### fixedtocsv

Convert a data-positioned (fixed width) file to CSV format.

```go
go get -u github.com/Guitarbum722/fixedtocsv

cd fixedtocsv

make build
```

#### Usage:
```
Usage: fixedtocsv [-c] [-d] [-f] [-o]
        Options:
          -h | --help  help
          -c           input configuration file (default: "config.json" in current directory)
          -d           output delimiter (default: comma ",")
          -f           input file name (Required)
          -o           output file name (default: "output.csv" in current directory)
```

The only required flag is `-f` which specifies the input file.

A JSON file is expected to provide the configuration of the input data (being the fixed width flat file).

Here is an example of the configuration file.  Most importantly, each object in `columnLens` represents a field in the file.  The order of the fields will be evaluated in the order in which they appear in the JSON.

```json
{
    "columnLens": [
        {
            "start": 0,
            "end": 6
        },
        {
            "start": 7,
            "end": 21
        }
    ]
}
```