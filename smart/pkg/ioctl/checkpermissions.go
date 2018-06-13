/*
Copyright 2018 The OpenEBS Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// See Linux man-pages http://man7.org/linux/man-pages/man2/capset.2.html

package ioctl

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	linuxCapabilityVersion3 = 0x20080522
	capSysRawIO             = 1 << 17
	capSysAdmin             = 1 << 21
)

type userCapHeader struct {
	version uint32
	pid     int
}

type userCapData struct {
	effective   uint32
	permitted   uint32
	inheritable uint32
}

type userCapsV3 struct {
	hdr  userCapHeader
	data [2]userCapData
}

// CheckBinaryPerm invokes the linux CAPGET syscall which checks for necessary capabilities required for a binary to access a device.
// Note : If the binary is executed as root, it automatically has all capabilities set.
func CheckBinaryPerm() error {
	userCaps := new(userCapsV3)
	userCaps.hdr.version = linuxCapabilityVersion3

	_, _, err := unix.RawSyscall(unix.SYS_CAPGET, uintptr(unsafe.Pointer(&userCaps.hdr)), uintptr(unsafe.Pointer(&userCaps.data)), 0)
	if err != 0 {
		return fmt.Errorf("linux capget syscall has failed: %+v", err)
	}

	if (userCaps.data[0].effective&capSysRawIO == 0) && (userCaps.data[0].effective&capSysAdmin == 0) {
		return fmt.Errorf("capSysRawIO and capSysAdmin are not in effect, device access will fail")
	}
	return nil
}
