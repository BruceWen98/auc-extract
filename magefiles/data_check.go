package main

import (
	"auc-video-extract/extract"
	"auc-video-extract/ocr"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/dsnet/try"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type DataCheck mg.Namespace

func (DataCheck) T01_find_min_and_max_frame(aucIndex int) {
	conf := extract.ReadConfig("config.yaml")
	auc := conf.Aucs[aucIndex]

	entries := try.E1(os.ReadDir(auc.InputPath))
	frames := []int{}
	re := regexp.MustCompile(`^frame\d+.jpg$`)
	for _, e := range entries {
		name := e.Name()
		//somefile is started with ._frame535185.jpg. need to ignore, so use regex to match
		if re.MatchString(name) {
			num := strings.ReplaceAll(strings.ReplaceAll(name, ".jpg", ""), "frame", "")
			index, err := strconv.Atoi(num)
			if err != nil {
				log.Printf("Failed to parse frame number from %s: %v", name, err)
			}
			frames = append(frames, index)
		}
	}

	sort.Ints(frames)
	log.Printf("frames total count:%d", len(frames))
	log.Printf("first 5 frames:%v", frames[:5])
	log.Printf("last 5 frames:%v", frames[len(frames)-5:])
	log.Printf("suggest report inerval:%.2f", float64(float64(frames[len(frames)-1]-frames[0])/float64(auc.Step)*0.1))

}

// sample the middle frame for sizing and text ocr
func (DataCheck) T02_check_middle_frame(aucIndex int) {
	conf := extract.ReadConfig("config.yaml")
	auc := conf.Aucs[aucIndex]

	frameId := auc.StartFrame + auc.Step*int((auc.EndFrame-auc.StartFrame)/auc.Step/2)
	log.Printf("middle frame:%d", frameId)
	file := fmt.Sprintf("%s/frame%d.jpg", auc.InputPath, frameId)
	log.Printf("image info with identify command:")
	sh.RunV("identify", file)

	os.MkdirAll(auc.OutputPath, 0755)
	o := ocr.NewOCR(frameId, fmt.Sprintf("frame%d.jpg", frameId), auc.InputPath, auc.OutputPath)
	text, err := o.Process(auc.CropSizeAndOffset, auc.PSM)
	if err != nil {
		log.Fatalf("Failed to process frame %d: %v", frameId, err)
	}
	log.Printf("text extracted:%s", text)
}

// test the frames from 10% to 100% every 10% to validate the OCR extraction
func (DataCheck) T03_sample_frame(aucIndex int) {
	conf := extract.ReadConfig("config.yaml")
	auc := conf.Aucs[aucIndex]

	frames := []int{}
	for _, percent := range []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0} {
		id := auc.StartFrame + auc.Step*int(float64(auc.EndFrame-auc.StartFrame)/float64(auc.Step)*percent)
		frames = append(frames, id)
	}
	for _, frameId := range frames {
		log.Printf("------------frame:%d---------------", frameId)
		o := ocr.NewOCR(frameId, fmt.Sprintf("frame%d.jpg", frameId), auc.InputPath, auc.OutputPath)
		text, err := o.Process(auc.CropSizeAndOffset, auc.PSM)
		if err != nil {
			log.Fatalf("Failed to process frame %d: %v", frameId, err)
		}
		log.Printf("text extracted:%s", text)

		itemIndex, title, price, confidence := extract.SearchPrice(text, auc)
		log.Printf("itemIndex:%s, title:%s, price:%.2f, confidence:%d", itemIndex, title, price, confidence)
		log.Printf("-----------------------")
	}
}

// check a specific frame for extraction and matching
func (DataCheck) T04_extract_a_frame(aucIndex int, frameId int) {
	conf := extract.ReadConfig("config.yaml")
	auc := conf.Aucs[aucIndex]

	o := ocr.NewOCR(frameId, fmt.Sprintf("frame%d.jpg", frameId), auc.InputPath, auc.OutputPath)
	text, err := o.Process(auc.CropSizeAndOffset, auc.PSM)
	if err != nil {
		log.Fatalf("Failed to process frame %d: %v", frameId, err)
	}
	log.Printf("text extracted:%s", text)

	itemIndex, title, price, confidence := extract.SearchPrice(text, auc)
	log.Printf("itemIndex:%s, title:%s, price:%.2f, confidence:%d", itemIndex, title, price, confidence)
}

// download the frame jpg files
func (DataCheck) T04_download_frame_picture(aucIndex int, frameId int) {
	conf := extract.ReadConfig("config.yaml")
	auc := conf.Aucs[aucIndex]

	file := fmt.Sprintf(`"%s/frame%d.jpg"`, auc.InputPath, frameId)
	log.Printf("download file:%s", file)
	bastion().Download(file, fmt.Sprintf("frame%d.jpg", frameId))
}
