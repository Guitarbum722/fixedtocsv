package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

type fixedWidthConfig struct {
	ColumnLens []struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"columnLens"`
	order []int // to order the keys in ColumnLens
}

func main() {
	configInput := `{
		"columnLens": [
	  {
		"start": 0,
		"end": 6
	  },
	  {
		"start": 7,
		"end": 21
	  },
	  {
		"start": 22,
		"end": 38
	  },
	  {
		"start": 39,
		"end": 63
	  },
	  {
		"start": 64,
		"end": 101
	  },
	  {
		"start": 102,
		"end": 136
	  },
	  {
		"start": 137,
		"end": 149
	  },
	  {
		"start": 150,
		"end": 163
	  }
	]
}
`

	conf := &fixedWidthConfig{}

	err := json.Unmarshal([]byte(configInput), conf)
	if err != nil {
		log.Fatalln("err parsing config file :\n", err)
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
