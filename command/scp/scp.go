package scp

import (
	"errors"
	"github.com/mhewedy/vermin/command"
	"github.com/mhewedy/vermin/db"
	"github.com/mhewedy/vermin/ip"
	"strings"
)

const (
	CopySeparator = ":"
)

type vmPath struct {
	name string
	file string
}

func Copy(src string, dest string) error {

	srcVmPath, srcIsVm := toVmPath(src)
	destVmPath, destIsVm := toVmPath(dest)

	if srcIsVm && destIsVm {
		return copyBetweenVMs(srcVmPath, destVmPath)
	} else if srcIsVm {
		return copyFromVM(srcVmPath, dest)
	} else if destIsVm {
		return copyToVM(src, destVmPath)
	} else {
		return errors.New("src/dest one of them should be vm")
	}
}

// convert <name>:<file> string to vmPath{name, file}
func toVmPath(srcDest string) (vmPath, bool) {
	if !strings.Contains(srcDest, CopySeparator) {
		return vmPath{}, false
	}

	vmAndPath := strings.Split(srcDest, CopySeparator)

	return vmPath{vmAndPath[0], vmAndPath[1]}, true
}

func copyFromVM(vmPath vmPath, localFile string) error {
	ipAddr, err := ip.Find(vmPath.name, false)
	if err != nil {
		return err
	}

	_, err = command.Scp(vmPath.name, db.GetUsername(vmPath.name)+"@"+ipAddr+":"+vmPath.file, localFile).Call()
	return err
}

func copyToVM(localFile string, vmPath vmPath) error {
	ipAddr, err := ip.Find(vmPath.name, false)
	if err != nil {
		return err
	}

	_, err = command.Scp(vmPath.name, localFile, db.GetUsername(vmPath.name)+"@"+ipAddr+":"+vmPath.file).Call()
	return err
}

func copyBetweenVMs(srcVmPath vmPath, destVmPath vmPath) error {

	srcIPAddr, err := ip.Find(srcVmPath.name, false)
	if err != nil {
		return err
	}

	destIPAddr, err := ip.Find(destVmPath.name, false)
	if err != nil {
		return err
	}

	_, err = command.Scp(srcVmPath.name, "-3",
		db.GetUsername(srcVmPath.name)+"@"+srcIPAddr+":"+srcVmPath.file,
		db.GetUsername(destVmPath.name)+"@"+destIPAddr+":"+destVmPath.file).Call()
	return err
}

// ----

func CopyToVM(vmName string, localFile string, vmFile string) error {
	return copyToVM(localFile, vmPath{vmName, vmFile})
}
