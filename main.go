package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
		"end": 7
	  },
	  {
		"start": 8,
		"end": 12
	  }
	]
}
`

	input := `John   1245
Peter  3545
Susan  6784
Sarah  4321
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

	fmt.Println(columns, conf.order)

	sr := strings.NewReader(input)
	scanner := bufio.NewScanner(sr)

	w := bufio.NewWriter(os.Stdout)

	for scanner.Scan() {
		line := scanner.Text()

		for _, v := range line {
			w.WriteRune(v)
		}
		w.WriteRune('\n')
	}

	w.Flush()
}
