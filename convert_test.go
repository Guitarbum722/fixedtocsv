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
            "end": 6
        },
        {
            "start": 7,
            "end": 21
		}
	]
}
`

var convertCases = []struct {
	input     string
	resultBuf int
}{
	{
		"1      Jeffry         ",
		9,
	},
	{
		"20     CaspŁr         ",
		11,
	},
}

func TestNewScanWriter(t *testing.T) {
	sw := newScanWriter(strings.NewReader("One  Two  Three"), bytes.NewBufferString(""), []byte(sampleConfig))

	if sw == nil {
		t.Fatal("failed initializing scanWriter\n")
	}
}

func TestConvert(t *testing.T) {
	for _, tt := range convertCases {
		sw := newScanWriter(strings.NewReader(tt.input), bytes.NewBufferString(""), []byte(sampleConfig))
		sw.convert()

		bufGot := sw.w.Buffered()
		if bufGot != tt.resultBuf {
			t.Fatalf("convert ( %v ) = %v; want %v", tt.input, bufGot, tt.resultBuf)
		}
	}
}