package csv

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Transaction struct {
	Timestamp   string `json:"timestamp"`
	Description string `json:"description"`
	Withdraw    int    `json:"withdraw"`
	Income      int    `json:"income"`
	Balance     int    `json:"balance"`
	Comment     string `json:"comment"`
}

const dateLayout = "%04s-%02s-%02T%02s:%02s:%02s|JST"

func ParseCSV(file_path string) ([]*Transaction, error) {
	inputFile, err := os.Open(file_path)
	if err != nil {
		return nil, err
	}
	defer inputFile.Close()

	reader := csv.NewReader(transform.NewReader(inputFile, japanese.ShiftJIS.NewDecoder()))
	reader.LazyQuotes = true

	transactions := []*Transaction{}
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		// columns: 年, 月, 日, 時, 分, 秒, 取引順番号, 摘要, お支払金額, お預り金額, 残高, メモ
		timestamp := fmt.Sprintf(dateLayout, line[0], line[1], line[2], line[3], line[4], line[5])
		withdraw, _ := strconv.Atoi(line[8])
		income, _ := strconv.Atoi(line[8])
		balance, _ := strconv.Atoi(line[8])
		tran := &Transaction{
			Timestamp:   timestamp,
			Description: line[7],
			Withdraw:    withdraw,
			Income:      income,
			Balance:     balance,
			Comment:     line[11],
		}
		transactions = append(transactions, tran)
	}

	return transactions, nil
}

// MarshallJSON marshallizes Transactions as regular JSON
func MarshallJSON(ts []*Transaction) ([]byte, error) {
	b, err := json.Marshal(ts)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// MarshallNewLineDelimitedJSON marshallizes Transactions as new-line delimited JSON
func MarshallNewLineDelimitedJSON(ts []*Transaction) ([]byte, error) {
	var ret []byte
	var wg sync.WaitGroup
	mutex := &sync.Mutex{}
	for _, t := range ts {
		wg.Add(1)
		tChan := make(chan *Transaction, 1)
		errChan := make(chan error, 1)
		go func(tr *Transaction) {
			b, err := json.Marshal(t)
			if err != nil {
				tChan <- tr
				errChan <- err
			} else {
				mutex.Lock()
				w := append(b, '\n')
				ret = append(ret, w...)
				mutex.Unlock()
				wg.Done()
			}
		}(t)

		if err := <-errChan; err != nil {
			t := <-tChan
			return nil, fmt.Errorf("failed to marshallize following (Transaction: %#v, err: %s)", t, err)
		}
	}
	return ret, nil
}
