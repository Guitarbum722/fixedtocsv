package main

import (
	"bytes"
	"strings"
	"testing"
)

const sampleConfig = `{
    "columnLens": [
        {
            "start": 0,
            "end": 7
        },
        {
            "start": 7,
            "end": 22
		}
	]
}
`

var convertCases = []struct {
	input     string
	resultBuf int // fields trimmed, plus delimiter, plus newline
}{
	{
		"1      Jeffry         ",
		9,
	},
	{
		"20     CaspŁr         ",
		11,
	},
	{
		"4333333Doll           ",
		13,
	},
}

func TestNewScanWriter(t *testing.T) {
	sw := newScanWriter(
		strings.NewReader("One  Two  Three"),
		bytes.NewBufferString(""),
		[]byte(sampleConfig),
		jsonConfig,
	)

	if sw == nil {
		t.Fatal("failed initializing scanWriter\n")
	}
}

func TestConvert(t *testing.T) {
	for _, tt := range convertCases {
		sw := newScanWriter(
			strings.NewReader(tt.input),
			bytes.NewBufferString(""),
			[]byte(sampleConfig),
			jsonConfig,
		)
		sw.convert()

		bufGot := sw.w.Buffered()
		if bufGot != tt.resultBuf {
			t.Fatalf("convert ( %v ) = %v; want %v", tt.input, bufGot, tt.resultBuf)
		}
	}
}
