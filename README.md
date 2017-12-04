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

A JSON or CSV file is expected to provide the configuration of the input data (being the fixed width flat file).

Currently, a file named `config.json` is expected by default.  Please see the flags above to specify a different file name.

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

If you want to use a CSV for the configuration, then it should have the below format (headers are expected).  You will need to use the `-csv` option AND the `-c` flag with a file name is REQUIRED in this case.

```
column,start,end
0,0,6
1,7,21
2,22,38
3,39,63
4,64,101
5,102,136
6,137,149
7,150,163
```
