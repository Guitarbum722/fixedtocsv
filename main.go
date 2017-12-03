package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	usage = `Usage: fixedtocsv [-c] [-d] [-f] [-o]
	Options:
	  -h | --help  help
	  -c           input configuration file (default: "config.json" in current directory) .
	  -d           output delimiter (default: comma ",").
	  -f           input file name (Required).
	  -o           output file name (default: "output.csv" in current directory).
	  -csv         if provided, configuration file is a CSV. (-c flag required with this option)
`
)

var (
	dFlag   = flag.String("d", ",", "")
	cFlag   = flag.String("c", "config.json", "")
	fFlag   = flag.String("f", "", "must specify an input file")
	oFlag   = flag.String("o", "output.csv", "")
	csvFlag = flag.Bool("csv", false, "")
)

type configType byte

// types of configuration files
const (
	jsonConfig configType = iota
	csvConfig
)

type columnLen struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type fixedWidthConfig struct {
	ColumnLens []columnLen `json:"columnLens"`
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

	var configFileType = jsonConfig

	if *csvFlag {
		configFileType = csvConfig
	}

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

	// read config file
	confInput, err := ioutil.ReadFile(*cFlag)
	if err != nil {
		log.Fatalln("unable to read config file : ", err)
	}

	sw := newScanWriter(ifp, ofp, confInput, configFileType) // initialize scanWriter

	// convert fixed width file to csv
	sw.convert()
	if err := sw.flush(); err != nil {
		log.Fatalln("err converting file : ", err)
	}
}

func newScanWriter(i io.Reader, o io.Writer, c []byte, t configType) *scanWriter {
	sw := &scanWriter{
		s: bufio.NewScanner(i),
		w: bufio.NewWriter(o),
	}

	sw.loadConfig(c, t)

	return sw
}

func (sw *scanWriter) loadConfig(confInput []byte, t configType) {
	c := &fixedWidthConfig{}

	if t == csvConfig {
		reader := csv.NewReader(bytes.NewReader(confInput))
		rows, err := reader.ReadAll()
		if err != nil {
			log.Fatalln("err reading csv config file : ", err)
		}
		c.ColumnLens = make([]columnLen, len(rows[1:]))

		// TODO: skip header row or option to do so?
		for i, v := range rows[1:] {
			start, err := strconv.Atoi(v[1])
			if err != nil {
				log.Fatalln("a non integer was provided in config file : ", err)
			}
			end, err := strconv.Atoi(v[2])
			if err != nil {
				log.Fatalln("a non integer was provided in config file : ", err)
			}

			c.ColumnLens[i].Start = start
			c.ColumnLens[i].End = end
		}
	} else {
		if err := json.Unmarshal(confInput, c); err != nil {
			log.Fatalln("err parsing config file :", err)
		}
	}

	sw.conf = c
}

func (sw *scanWriter) convert() {
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
}

func (sw *scanWriter) flush() error {
	return sw.w.Flush()
}
