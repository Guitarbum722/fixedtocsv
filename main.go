package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
