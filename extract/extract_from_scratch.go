package extract

import (
	"auc-video-extract/ocr"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/dsnet/try"
)

type AuctionText struct {
	timeIndex int
	text      string
}

type AuctionPrice struct {
	timeIndex    int
	itemIndex    string
	title        string
	price        float64
	lowEstimate  float64
	highEstimate float64

	confidence int //for price
}

func ExtractAuction(threads int, auc AuctionConfig) {
	startTime := time.Now()

	ch := make(chan *ocr.Ocr)
	resCh := make(chan AuctionText)

	var wgExtract, wgResult sync.WaitGroup
	var muText, muPrice sync.Mutex
	textResult := []AuctionText{}
	priceResult := []AuctionPrice{}

	for i := 0; i < threads; i++ {
		wgExtract.Add(1)
		go func() {
			defer wgExtract.Done()
			for ocr := range ch {
				text, err := ocr.Process(auc.CropSizeAndOffset, auc.PSM)
				if err != nil {
					log.Printf("Failed to process frame %s: %v", ocr.File, err)
					continue
				}
				atext := AuctionText{timeIndex: ocr.Frame, text: text}
				resCh <- atext
				muText.Lock()
				textResult = append(textResult, atext)
				muText.Unlock()
			}
		}()
	}

	for i := 0; i < threads; i++ {
		wgResult.Add(1)
		go func() {
			defer wgResult.Done()
			for res := range resCh {
				itemIndex, title, price, confidence := SearchPrice(res.text, auc)
				lowEst, highEst := SearchEstimation(res.text, auc)

				percentDone := float64(res.timeIndex-auc.StartFrame) / float64(auc.EndFrame-auc.StartFrame) * 100.0
				if (res.timeIndex-auc.StartFrame)%auc.ReportInterval == 0 {
					log.Printf("Has processed %.2f%% items, elapsed time: %.2f min", percentDone, time.Since(startTime).Minutes())
				}
				muPrice.Lock()
				priceResult = append(priceResult, AuctionPrice{
					timeIndex:    res.timeIndex,
					itemIndex:    itemIndex,
					title:        title,
					price:        price,
					confidence:   confidence,
					lowEstimate:  lowEst,
					highEstimate: highEst,
				})
				muPrice.Unlock()
			}
		}()
	}

	try.E(os.MkdirAll(auc.OutputPath, 0755))
	for i := auc.StartFrame; i < auc.EndFrame; i += auc.Step {
		ch <- ocr.NewOCR(i, fmt.Sprintf("frame%d.jpg", i), auc.InputPath, auc.OutputPath)
	}
	close(ch)
	wgExtract.Wait()
	close(resCh)
	wgResult.Wait()

	saveTextResult(auc, textResult)
	savePriceResult(auc, priceResult)
}

func matchPrice(text, pat string, auc AuctionConfig) (bool, string, string, float64) {
	var itemIndex string
	var price float64
	var title string

	re := regexp.MustCompile(pat)
	reNonDigit := regexp.MustCompile(`[^0-9]+`)
	matches := re.FindStringSubmatch(text)
	if matches != nil { //matched
		var err error
		itemIndex = matches[auc.BidIndex]
		title = matches[auc.TitleIndex]
		priceText := matches[auc.PriceIndex]
		priceText = reNonDigit.ReplaceAllString(priceText, "")

		price, err = strconv.ParseFloat(priceText, 64)
		if err != nil {
			log.Printf("failed to convert price to int: %v", err)
		}
		return true, itemIndex, title, price
	}
	return false, itemIndex, title, price
}

func matchEstimation(text, pat string, auc AuctionConfig) (bool, float64, float64) {
	var lowEstimate float64
	var highEstimate float64

	re := regexp.MustCompile(pat)
	reNonDigit := regexp.MustCompile(`[^0-9]+`)
	matches := re.FindStringSubmatch(text)
	if matches != nil { //matched
		var err error
		lowEstimateText := matches[auc.EstimateLowIndex]
		highEstimateText := matches[auc.EstimateHighIndex]
		lowEstimateText = reNonDigit.ReplaceAllString(lowEstimateText, "")
		highEstimateText = reNonDigit.ReplaceAllString(highEstimateText, "")

		lowEstimate, err = strconv.ParseFloat(lowEstimateText, 64)
		if err != nil {
			log.Printf("failed to convert low estimate to int: %v", err)
		}
		highEstimate, err = strconv.ParseFloat(highEstimateText, 64)
		if err != nil {
			log.Printf("failed to convert high estimate to int: %v", err)
		}
		return true, lowEstimate, highEstimate
	}
	return false, lowEstimate, highEstimate
}

func SearchPrice(text string, auc AuctionConfig) (string, string, float64, int) {
	var itemIndex string
	var price float64
	var title string
	var matched bool

	matched, itemIndex, title, price = matchPrice(text, auc.ExtractPatternOneLine, auc)
	if matched {
		return itemIndex, title, price, 100
	}

	matched, itemIndex, title, price = matchPrice(text, auc.ExtractPatternMultipleLine, auc)
	if matched {
		return itemIndex, title, price, 85
	}
	return itemIndex, title, price, 0
}

func SearchEstimation(text string, auc AuctionConfig) (float64, float64) {
	var lowEstimate float64
	var highEstimate float64
	var matched bool

	matched, lowEstimate, highEstimate = matchEstimation(text, auc.ExtractEstimatePattern, auc)
	if matched {
		return lowEstimate, highEstimate
	}
	return lowEstimate, highEstimate
}

func saveTextResult(auc AuctionConfig, textResult []AuctionText) {
	textResultFile := fmt.Sprintf("%s/%s", auc.OutputPath, auc.TextResultFile)
	f := try.E1(os.Create(textResultFile))
	defer f.Close()

	w := csv.NewWriter(f)
	for _, res := range textResult {
		w.Write([]string{strconv.Itoa(res.timeIndex), res.text})
	}
	w.Flush()
}

func savePriceResult(auc AuctionConfig, priceResult []AuctionPrice) {
	priceResultFile := fmt.Sprintf("%s/%s", auc.OutputPath, auc.PriceResultFile)
	f := try.E1(os.Create(priceResultFile))
	defer f.Close()

	w := csv.NewWriter(f)
	for _, res := range priceResult {
		w.Write([]string{
			strconv.Itoa(res.timeIndex),
			res.itemIndex,
			res.title,
			strconv.FormatFloat(res.price, 'f', 2, 64),
			strconv.Itoa(res.confidence),
			strconv.FormatFloat(res.lowEstimate, 'f', 2, 64),
			strconv.FormatFloat(res.highEstimate, 'f', 2, 64),
		})
	}
	w.Flush()
}
