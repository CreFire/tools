//go:build windows

package xid

import (
	"fmt"
	"syscall"
	"unsafe"
)

//readPlatformMachineID source: https://github.com/shirou/gopsutil/blob/master/host/host_syscall.go
func readPlatformMachineID() (string, error) {
	var h syscall.Handle
	utf16PtrFromString, err := syscall.UTF16PtrFromString(`SOFTWARE\Microsoft\Cryptography`)
	if err != nil {
		return "", err
	}
	err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE,
		utf16PtrFromString,
		0, syscall.KEY_READ|syscall.KEY_WOW64_64KEY, &h)
	if err != nil {
		return "", err
	}
	defer func() {
		err := syscall.RegCloseKey(h)
		if err != nil {
		}
	}()
	// len(`{`) + len(`abcdefgh-1234-456789012-123345456671` * 2) + len(`}`)
	//2 == bytes/UTF16
	const syscallRegBufLen = 74
	const uuidLen = 36
	var (
		regBuf  [syscallRegBufLen]uint16
		bufLen  = uint32(syscallRegBufLen)
		valType uint32
	)
	utf16Ptr, err := syscall.UTF16PtrFromString(`MachineGuid`)
	if err != nil {
		return "", err
	}
	err = syscall.RegQueryValueEx(h, utf16Ptr, nil, &valType, (*byte)(unsafe.Pointer(&regBuf[0])), &bufLen)
	if err != nil {
		return "", err
	}
	hostId := syscall.UTF16ToString(regBuf[:])
	if len(hostId) != uuidLen {
		return "", fmt.Errorf("HostID incorrect: %q\n", hostId)
	}
	return hostId, nil
}
