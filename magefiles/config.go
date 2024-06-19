package main

import (
	"auc-video-extract/extract"
	"log"

	"github.com/magefile/mage/mg"
)

type Conf mg.Namespace

func (Conf) T01_print_config() {
	conf := extract.ReadConfig("config.yaml")
	log.Printf("config:%#v", conf)
}
