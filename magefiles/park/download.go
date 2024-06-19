package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dsnet/try"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/zhiminwen/quote"
)

type Download mg.Namespace

// supply file id as argument
func (Download) T01_download_from_google_drive(fileId string) {
	// https://drive.google.com/file/d/1_os4zWBinTip9G_q0Rtm5rLGC-JnzLEe/view?usp=sharing
	// id is then 1_os4zWBinTip9G_q0Rtm5rLGC-JnzLEe
	cmd := quote.CmdTemplate(`
		cd /data
		curl -c cookie -s -L "https://drive.google.com/uc?export=download&id={{ .fieldId }}"
	`, map[string]string{
		"fieldId": fileId,
	})
	content := try.E1(sh.Output("sh", "-c", cmd))
	doc := try.E1(goquery.NewDocumentFromReader(strings.NewReader(content)))

	action := doc.Find("form").First().AttrOr("action", "")
	cmd = quote.CmdTemplate(`
		cd /data
		curl -Lb ./cookie -o data.zip "{{ .action }}" --http1.1
		`, map[string]string{
		"dir":    workingDir,
		"action": action,
	})
	sh.RunV("sh", "-c", cmd)
}
