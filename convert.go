package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"strings"
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
		fields := make([]string, len(sw.conf.ColumnLens))

		// split line into a slice of strings based on length configuration
		// then trim surrounding space.
		for i := range sw.conf.ColumnLens {
			fields[i] = line[sw.conf.ColumnLens[i].Start:sw.conf.ColumnLens[i].End]
			fields[i] = strings.Trim(fields[i], " ")
		}

		sw.w.WriteString(strings.Join(fields, *dFlag) + "\n")
	}
}

func (sw *scanWriter) flush() error {
	return sw.w.Flush()
}
