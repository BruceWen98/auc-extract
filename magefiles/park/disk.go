package main

import (
	"github.com/magefile/mage/mg"
	"github.com/zhiminwen/quote"
)

type Disk mg.Namespace

func (Disk) T01_add_2nd_disk() {
	cmd := quote.CmdTemplate(`
				qemu-img create -f raw /data1/kvm-images/{{ .vmName }}-data-disk.img {{ .disk2Size}}G
				virsh attach-disk {{ .vmName }} --source /data1/kvm-images/{{ .vmName }}-data-disk.img --target vdb --persistent
		`, map[string]string{
		"vmName":    "test-bed",
		"disk2Size": "200",
	})

	master.Execute(cmd)
}

func (Disk) T02_create_filesystem() {
	cmd := quote.CmdTemplate(`
		sudo pvcreate /dev/vdb
		sudo vgcreate vg-data /dev/vdb
		sudo lvcreate -l 100%FREE -n lv-data vg-data
		sudo mkfs.ext4 /dev/vg-data/lv-data
		sudo mkdir -p /data
		sudo mount /dev/vg-data/lv-data /data
		echo "/dev/vg-data/lv-data /data ext4 defaults 0 0" | sudo -a /etc/fstab
	`, map[string]string{})

	bastion().Execute(cmd)
}
