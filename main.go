package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	usage = `Usage: fixedtocsv [-d] [-f] [-o]
	Options:
	  -h | --help  help
	  -d           output delimiter (default: comma ",")
	  -f           input configuration file (default: "config.json" in current directory) 
	  -o           output file name (default: "output.csv" in current directory)
`
)

var (
	dFlag = flag.String("d", ",", "")
	fFlag = flag.String("f", "config.json", "")
	oFlag = flag.String("o", "output.csv", "")
)

type fixedWidthConfig struct {
	ColumnLens []struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"columnLens"`
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}
	flag.Parse()

	var confInput []byte
	var err error

	confInput, err = ioutil.ReadFile(*fFlag)

	if err != nil {
		log.Fatalln("unable to read config file : ", err)
	}

	conf := &fixedWidthConfig{}

	err = json.Unmarshal(confInput, conf)
	if err != nil {
		log.Fatalln("err parsing config file :", err)
	}

	sr := strings.NewReader(input)
	scanner := bufio.NewScanner(sr)

	fp, err := os.Create(*oFlag)
	if err != nil {
		log.Fatalln(err)
	}
	w := bufio.NewWriter(fp)

	for scanner.Scan() {
		line := scanner.Text()
		var fields = make([]string, 0, len(conf.ColumnLens))

		// split line into a slice of strings based on length configuration
		for i := range conf.ColumnLens {
			fields = append(fields, line[conf.ColumnLens[i].Start:conf.ColumnLens[i].End])
		}

		for i := range fields {
			fields[i] = strings.Trim(fields[i], " ")
		}

		w.WriteString(strings.Join(fields, *dFlag) + "\n")
	}

	w.Flush()
}
