//go:build darwin
// +build darwin

package xid

import (
	"errors"
	"io/ioutil"
	"runtime"
	"syscall"
)

//readPlatformMachineID source: https://github.com/shirou/gopsutil/blob/master/host/host_syscall.go
func readPlatformMachineID() (string, error) {
	switch runtime.GOOS {
	case `windows`:
		return windows()
	case `linux`:
		{
			b, err := ioutil.ReadFile("/etc/machine-id")
			if err != nil || len(b) == 0 {
				b, err = ioutil.ReadFile("/sys/class/dmi/id/product_uuid")
			}
			return string(b), err
		}
	case `freebsd`:
		return syscall.Sysctl("kern.hostuuid"), nil
	case `darwin`:
		return syscall.Sysctl("kern.uuid"), nil
	default:
		return "", errors.New("not implemented")
	}

}
