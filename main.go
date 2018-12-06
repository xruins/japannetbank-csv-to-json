package main

import (
	"log"
	"os"
	"time"

	"github.com/xruins/japannetbank-csv-to-json/csv"
)

type Transaction struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Withdraw    int       `json:"withdraw"`
	Income      int       `json:"income"`
	Balance     int       `json:"balance"`
	Comment     string    `json:"comment"`
}

const (
	JSON = iota
	NdJSON
)

func main() {
	var format
	flag.Bool(&delimitWithNewLine, ")

	if len(os.Args) != 2 {
		log.Fatalf("usage: %s input_file output_file", os.Args[0])
	}

	input := os.Args[1]
	output := os.Args[2]

	ts, err := csv.ParseCSV(input)

}
