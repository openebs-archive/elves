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

// Package smartinfo is a pure Go SMART library.
//
package smartinfo

import (
	"fmt"
	"path/filepath"

	"github.com/openebs/elves/smart/pkg/ioctl"
	"github.com/openebs/elves/smart/pkg/scsismart"
)

// ScanDevices discover and return the list of scsi devices.
func ScanDevices() []scsismart.SCSIDev {
	devices := make([]scsismart.SCSIDev, 0)
	// Find all SCSI disk devices
	files, err := filepath.Glob("/dev/sd*[^0-9]")
	if err != nil {
		return devices
	}

	for _, file := range files {
		devices = append(devices, scsismart.SCSIDev{DevName: file})
	}

	return devices
}

// SCSIBasicDiskInfo returns details(disk attributes and their values such as vendor,serialno,etc) of a disk
func SCSIBasicDiskInfo(device string, attrName string) (string, error) {

	// Check if required permissions are present or not for accessing a device
	err := ioctl.CheckBinaryPerm()
	if err != nil {
		return "", fmt.Errorf("error while checking device access permissions, Error: %+v", err)
	}

	if device == "" {
		return "", fmt.Errorf("no disk device path given to get the disk details")
	}

	d, err := scsismart.DetectSCSIType(device)
	if err != nil {
		return "", fmt.Errorf("error in detecting type of SCSI device, Error: %+v", err)
	}
	defer d.Close()

	AttrDetail, err := d.GetBasicDiskInfo(attrName)
	if err != nil {
		return AttrDetail, fmt.Errorf("error getting %q detail of disk device having devpath %q", attrName, device)
	}
	return AttrDetail, nil
}
