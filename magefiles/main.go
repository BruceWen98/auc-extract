package main

import (
	"os"

	"github.com/dsnet/try"
	"github.com/zhiminwen/magelib/tunnel"
	"github.com/zhiminwen/magetool/sshkit"
	"github.com/zhiminwen/quote"
)

var master *sshkit.SSHClient
var project, workingDir string

var bastionClient *sshkit.SSHClient
var tunnelPort, tunnelKey string
var tun *tunnel.Tunnel

func init() {
	os.Setenv("MAGEFILE_VERBOSE", "true")

	project = "auc-extract"
	workingDir = project
}

func bastion() *sshkit.SSHClient {
	if bastionClient != nil {
		return bastionClient
	}
	keyPath := try.E1(os.UserHomeDir()) + "/.ssh/id_rsa"
	master = try.E1(sshkit.NewSSHClient("bare02.ibmcloud.io.cpak", "22", "root", "", keyPath))
	bastionClient = try.E1(master.ProxySSHClient("192.168.10.220", "22", "ubuntu", "", keyPath))

	return bastionClient
}

func T00_test_ssh_connection() {
	cmd := quote.CmdTemplate(`
		hostname
	`, map[string]string{})
	bastion().Execute(cmd)
}

func T01_init_namespace() {
	cmd := quote.CmdTemplate(`
		mkdir -p {{ .dir }}
	`, map[string]string{
		"project": project,
		"dir":     workingDir,
	})
	bastion().Execute(cmd)
}
