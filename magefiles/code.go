package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/zhiminwen/quote"
)

type Code mg.Namespace

func (Code) T01_upload() {
	cmd := quote.CmdTemplate(`
		cd {{ .dir}}
		rm -rf {{ .list }}
	`, map[string]string{
		"dir":  workingDir,
		"list": "extract magefiles ocr performace video config.yaml go*",
	})
	bastion().Execute(cmd)

	cmd = quote.CmdTemplate(`
		tar --no-xattrs -cf - {{ .list }} | ssh ibm-brauc '(cd {{ .workingDir }}; tar xf -)'
	`, map[string]string{
		"list":       "extract magefiles ocr performance video config.yaml go*",
		"workingDir": workingDir,
	})
	sh.RunV("sh", "-c", cmd)
}
