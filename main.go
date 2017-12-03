package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	usage = `Usage: fixedtocsv [-c] [-d] [-f] [-o]
	Options:
	  -h | --help  help
	  -c           input configuration file (default: "config.json" in current directory) 
	  -d           output delimiter (default: comma ",")
	  -f           input file name (Required)
	  -o           output file name (default: "output.csv" in current directory)
`
)

var (
	dFlag = flag.String("d", ",", "")
	cFlag = flag.String("c", "config.json", "")
	fFlag = flag.String("f", "", "must specify an input file")
	oFlag = flag.String("o", "output.csv", "")
)

type fixedWidthConfig struct {
	ColumnLens []struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"columnLens"`
}

// scanWriter scans data positioned input and writes to the configured writer.
type scanWriter struct {
	conf *fixedWidthConfig
	s    *bufio.Scanner
	w    *bufio.Writer
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}
	flag.Parse()

	var ifp *os.File
	ifp, err := os.Open(*fFlag)
	if err != nil {
		log.Fatalln("unable to open input file : ", err)
	}
	defer ifp.Close()

	ofp, err := os.Create(*oFlag)
	if err != nil {
		log.Fatalln(err)
	}
	defer ofp.Close()

	sw := newScanWriter(ifp, ofp, *cFlag) // initialize scanWriter

	// convert fixed width file to csv
	if err := sw.convert(); err != nil {
		log.Fatalln("err converting file : ", err)
	}
}

func newScanWriter(i io.Reader, o io.Writer, c string) *scanWriter {
	sw := &scanWriter{
		s: bufio.NewScanner(i),
		w: bufio.NewWriter(o),
	}
	sw.loadConfig(c)

	return sw
}

func (sw *scanWriter) loadConfig(fileName string) {
	confInput, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("unable to read config file : ", err)
	}
	c := &fixedWidthConfig{}

	if err = json.Unmarshal(confInput, c); err != nil {
		log.Fatalln("err parsing config file :", err)
	}

	sw.conf = c
}

func (sw *scanWriter) convert() error {
	for sw.s.Scan() {
		line := sw.s.Text()
		var fields = make([]string, 0, len(sw.conf.ColumnLens))

		// split line into a slice of strings based on length configuration
		for i := range sw.conf.ColumnLens {
			fields = append(fields, line[sw.conf.ColumnLens[i].Start:sw.conf.ColumnLens[i].End])
		}

		for i := range fields {
			fields[i] = strings.Trim(fields[i], " ")
		}

		sw.w.WriteString(strings.Join(fields, *dFlag) + "\n")
	}

	return sw.w.Flush()
}
