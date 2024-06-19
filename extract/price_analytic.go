package extract

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/dsnet/try"
)

func AnalyzePriceCSV(filename string) (int, int) {
	f := try.E1(os.Open(filename))
	defer f.Close()

	r := csv.NewReader(f)
	count := 0
	successCount := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		// 50190,3B,Willem de Kooning | East Hampton VI,8600000.00,100
		count += 1
		if record[1] == "" || record[3] == "" {
			continue
		}
		successCount += 1
	}

	return count, successCount
}
