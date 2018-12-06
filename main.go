package main

import (
	"io"
	"os"
	"time"

	"github.com/urfave/cli"
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

func main() {
	app := cli.NewApp()
	app.Name = "japannnetbank-csv-to-json"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "input", Value: "", Usage: "Input path of CSV file (default: stdin)"},
		cli.StringFlag{Name: "outout", Value: "", Usage: "Output path of converted data (default: stdout)"},
		cli.StringFlag{Name: "format", Value: "json", Usage: "Export format (json/ndjson)"},
	}
	app.Usage = "Convert CSV of Japan Net Bank transactions to JSON"
	app.Action = func(c *cli.Context) error {
		args := c.Args()
		input := args.Get(0)
		output := args.Get(1)

		ts, err := csv.ParseCSV(input)
		if err != nil {
			return err
		}

		var b []byte
		format := c.String("format")
		switch {
		case format == "ndjson":
			b, err = csv.MarshallNewLineDelimitedJSON(ts)
		case format == "json":
			b, err = csv.MarshallJSON(ts)
		default:
			b, err = csv.MarshallJSON(ts)
		}
		if err != nil {
			return err
		}

		var writer io.Writer
		if output == "" {
			writer = os.Stdout
		} else {
			file, err := os.Create(output)
			if err != nil {
				return err
			}
			defer file.Close()
			writer = file
		}

		writer.Write(b)
		return nil
	}
	app.Run(os.Args)
}
