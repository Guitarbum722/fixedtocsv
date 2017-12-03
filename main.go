package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const defaultConfigFile = "config.json"

type fixedWidthConfig struct {
	ColumnLens []struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"columnLens"`
	order []int // to order the keys in ColumnLens
}

func main() {
	var confInput []byte
	var err error

	confInput, err = ioutil.ReadFile(defaultConfigFile)

	if err != nil {
		log.Fatalln("unable to read config file : ", err)
	}

	conf := &fixedWidthConfig{}

	err = json.Unmarshal(confInput, conf)
	if err != nil {
		log.Fatalln("err parsing config file :", err)
	}

	columns := make(map[int]int)

	for i, v := range conf.ColumnLens {
		conf.order = append(conf.order, i)
		columns[i] = v.End - v.Start
	}

	sr := strings.NewReader(input)
	scanner := bufio.NewScanner(sr)

	fp, err := os.Create("output.csv")
	if err != nil {
		log.Fatalln(err)
	}
	w := bufio.NewWriter(fp)

	for scanner.Scan() {
		line := scanner.Text()
		var fields = make([]string, 0, len(conf.ColumnLens))

		// split line into a slice of strings based on length configuration
		for i := range conf.order {
			fields = append(fields, line[conf.ColumnLens[i].Start:conf.ColumnLens[i].End])
		}

		for i := range fields {
			fields[i] = strings.Trim(fields[i], " ")
		}

		w.WriteString(strings.Join(fields, ",") + "\n")
	}

	w.Flush()
}
