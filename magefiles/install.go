package main

import (
	"github.com/magefile/mage/mg"
	"github.com/zhiminwen/quote"
)

type Install mg.Namespace

func (Install) T01_install_tool() {
	cmd := quote.CmdTemplate(`
		sudo apt install -y tesseract-ocr-eng tesseract-ocr-chi-sim tesseract-ocr-chi-tra libtesseract-dev libleptonica-dev
		sudo apt install -y imagemagick
	`, map[string]string{})

	bastion().Execute(cmd)
}

func (Install) T02_install_go() {
	cmd := quote.CmdTemplate(`
		cd {{ .dir }}
		curl -LO https://go.dev/dl/go1.22.4.linux-amd64.tar.gz
		sudo rm -rf /usr/local/go
		sudo tar -C /usr/local -xzf go1.22.4.linux-amd64.tar.gz
		sudo ln -s /usr/local/go/bin/go /usr/local/bin/go
	`, map[string]string{
		"dir": workingDir,
	})

	bastion().Execute(cmd)
}

func (Install) T03_install_mage() {
	cmd := quote.CmdTemplate(`
		cd {{ .dir }}
		// git clone https://github.com/magefile/mage
		cd mage
		mkdir -p $HOME/go/bin
		go run bootstrap.go
	`, map[string]string{
		"dir": workingDir,
	})
	bastion().Execute(cmd)
	// after that do: sudo cp mage /usr/local/bin
}
