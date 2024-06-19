package extract

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/dsnet/try"
)

// assuming the text is already extracted saved into a csv file
// then just read the csv file and extract out the bidding price
func ExtractPriceFromTextFile(threads int, auc AuctionConfig) {
	textCh := make(chan AuctionText)
	priceResult := []AuctionPrice{}

	var muPrice sync.Mutex
	var wg sync.WaitGroup
	// startTime := time.Now()
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for res := range textCh {
				itemIndex, title, price, confidence := SearchPrice(res.text, auc)
				// percentDone := float64(res.timeIndex-auc.StartFrame) / float64(auc.EndFrame-auc.StartFrame) * 100.0
				// if (res.timeIndex-auc.StartFrame)%auc.ReportInterval == 0 {
				// 	log.Printf("Has processed %.2f%% items, elapsed time: %.2f min", percentDone, time.Since(startTime).Minutes())
				// }
				muPrice.Lock()
				priceResult = append(priceResult, AuctionPrice{
					timeIndex:  res.timeIndex,
					itemIndex:  itemIndex,
					title:      title,
					price:      price,
					confidence: confidence,
				})
				muPrice.Unlock()
			}
		}()
	}

	f := try.E1(os.Open(fmt.Sprintf("%s/%s", auc.OutputPath, auc.TextResultFile)))
	defer f.Close()

	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}
		frameId := try.E1(strconv.Atoi(record[0]))
		text := record[1]

		textCh <- AuctionText{timeIndex: frameId, text: text}
	}

	close(textCh)
	wg.Wait()
	savePriceResult(auc, priceResult)
}
