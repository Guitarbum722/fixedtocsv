package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type fixedWidthConfig struct {
	ColumnLens []struct {
		Start int `json:"start"`
		End   int `json:"end"`
	} `json:"columnLens"`
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

	input := `This is a header line to be skipped
John   1245
Peter  3545
Susan  6784
Sarah  4321
`

	conf := &fixedWidthConfig{}

	err := json.Unmarshal([]byte(configInput), conf)
	if err != nil {
		log.Fatalln("err parsing json :\n", err)
	}

	columns := make(map[int]int)

	for i, v := range conf.ColumnLens {
		columns[i] = v.End - v.Start
	}

	fmt.Println(columns)

	sr := strings.NewReader(input)
	scanner := bufio.NewScanner(sr)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
