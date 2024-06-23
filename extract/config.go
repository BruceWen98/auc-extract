package extract

import (
	"os"

	"github.com/dsnet/try"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Threads int             `yaml:"threads"`
	Aucs    []AuctionConfig `yaml:"aucs"`
}

type AuctionConfig struct {
	Name           string `yaml:"name"`
	StartFrame     int    `yaml:"start_frame"`
	EndFrame       int    `yaml:"end_frame"`
	Step           int    `yaml:"step"`
	ReportInterval int    `yaml:"report_interval"`

	InputPath                  string `yaml:"input_path"`
	OutputPath                 string `yaml:"output_path"`
	ExtractPatternOneLine      string `yaml:"extract_pattern_one_line"`
	ExtractPatternMultipleLine string `yaml:"extract_pattern_multiple_line"`
	BidIndex                   int    `yaml:"bid_index"`
	TitleIndex                 int    `yaml:"title_index"`
	PriceIndex                 int    `yaml:"price_index"`

	ExtractEstimatePattern string `yaml:"extract_estimation_pattern"`
	EstimateLowIndex       int    `yaml:"estimation_low_index"`
	EstimateHighIndex      int    `yaml:"estimation_high_index"`

	CropSizeAndOffset string `yaml:"crop_size_and_offset"`
	PSM               string `yaml:"psm"`
	//if -1 then not using PSM

	TextResultFile  string `yaml:"text_result_file"`
	PriceResultFile string `yaml:"price_result_file"`

	Done bool `yaml:"done"`
}

func ReadConfig(file string) *Config {
	var config Config
	content := try.E1(os.ReadFile(file))
	try.E(yaml.Unmarshal(content, &config))

	for i := range config.Aucs {
		if config.Aucs[i].BidIndex == 0 {
			config.Aucs[i].BidIndex = 1
		}
		if config.Aucs[i].TitleIndex == 0 {
			config.Aucs[i].TitleIndex = 2
		}
		if config.Aucs[i].PriceIndex == 0 {
			config.Aucs[i].PriceIndex = 3
		}
	}
	return &config
}
