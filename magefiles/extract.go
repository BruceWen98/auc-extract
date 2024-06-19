package main

import (
	"auc-video-extract/extract"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/zhiminwen/quote"
)

type Extract mg.Namespace

// extract by config.yaml
func (Extract) T01_extract_into_csv() {
	conf := extract.ReadConfig("config.yaml")

	for _, auc := range conf.Aucs {
		log.Printf("processing auc:%s", auc.Name)
		if auc.Done {
			log.Printf("auc %s is done, skip", auc.Name)
			continue
		}
		now := time.Now()
		extract.ExtractAuction(conf.Threads, auc)
		log.Printf("time spent:%.2f min", time.Since(now).Minutes())
	}
}

// when rerun, just read from the saved text csv file and extract out the price
func (Extract) T10_extract_from_text_csv() {
	conf := extract.ReadConfig("config.yaml")

	for _, auc := range conf.Aucs {
		log.Printf("processing auc:%s", auc.Name)
		if auc.Done {
			log.Printf("auc %s is done, skip", auc.Name)
			continue
		}
		now := time.Now()
		extract.ExtractPriceFromTextFile(conf.Threads, auc)
		log.Printf("time spent:%.2f min", time.Since(now).Minutes())
	}
}

func (Extract) T02_download_csv(index int) {
	conf := extract.ReadConfig("config.yaml")
	auc := conf.Aucs[index]

	for _, f := range []string{auc.TextResultFile, auc.PriceResultFile} {
		file := fmt.Sprintf(`%s/%s`, auc.OutputPath, f)
		bastion().Download(file, filepath.Base(file))
	}
}

// func (Extract) T10_test_frame(frameID int) {
// 	conf := extract.ReadConfig("config.yaml")
// 	auc := conf.Aucs[0]
// 	o := ocr.NewOCR(frameID, fmt.Sprintf("frame%d.jpg", frameID), auc.InputPath, auc.OutputPath)
// 	text, err := o.Process(auc.CropSizeAndOffset, auc.PSM)
// 	if err != nil {
// 		log.Fatalf("Failed to process frame %d: %v", frameID, err)
// 	}
// 	log.Printf("text extracted:%s", text)
// }

// analysis the price csv file to find out success rate
func (Extract) T10_anlyze_price_csv(index int) {
	conf := extract.ReadConfig("config.yaml")
	auc := conf.Aucs[index]

	file := fmt.Sprintf(`%s/%s`, auc.OutputPath, auc.PriceResultFile)
	count, successCount := extract.AnalyzePriceCSV(file)

	log.Printf("Total records: %d, Success: %d, %.2f %%", count, successCount, float64(successCount)/float64(count)*100)
}

func (Extract) T11_download_csv(index int) {
	conf := extract.ReadConfig("config.yaml")
	auc := conf.Aucs[index]

	for _, f := range []string{auc.TextResultFile, auc.PriceResultFile} {
		file := fmt.Sprintf(`%s/%s`, auc.OutputPath, f)
		bastion().Download(file, filepath.Base(file))
	}
}

func (Extract) T12_download_all_csv() {
	conf := extract.ReadConfig("config.yaml")

	os.MkdirAll("results", 0755)
	for _, auc := range conf.Aucs {
		for _, f := range []string{auc.TextResultFile, auc.PriceResultFile} {
			file := fmt.Sprintf(`%s/%s`, auc.OutputPath, f)
			bastion().Download(file, "results/"+filepath.Base(file))
		}
	}
}

// remove the cropped images
func (Extract) T20_clean() {
	conf := extract.ReadConfig("config.yaml")

	for _, auc := range conf.Aucs {
		cmd := quote.CmdTemplate(`
			cd {{ .outPath }}
			find . -name "*.jpg" -delete
		`, map[string]string{
			"outPath": auc.OutputPath,
		})
		sh.RunV("bash", "-c", cmd)
	}
}
