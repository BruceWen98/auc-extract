package performance

import (
	"auc-video-extract/extract"
	"auc-video-extract/ocr"
	"fmt"
	"log"
	"testing"
)

func TestProcess(t *testing.T) {
	conf := extract.ReadConfig("../config.yaml")
	frameID := 458805
	auc := conf.Aucs[0]
	o := ocr.NewOCR(frameID, fmt.Sprintf("frame%d.jpg", frameID), auc.InputPath, auc.OutputPath)
	text, err := o.Process(auc.CropSizeAndOffset, auc.PSM)
	if err != nil {
		log.Fatalf("Failed to process frame %d: %v", frameID, err)
	}
	log.Printf("text extracted:%s", text)
}
