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

// SCSI generic IO functions.

package scsismart

import (
	"strconv"

	"golang.org/x/sys/unix"
)

// Dev is the top-level device interface. All supported device types must implement these interfaces.
type Dev interface {
	DevOpen
	DevClose
	Devinfo
}

// DevOpen is the interface which implements open method for opening a disk device
type DevOpen interface {
	Open() error
}

// DevClose is the interface which implements close method for closing a disk device
type DevClose interface {
	Close() error
}

// Devinfo is the interface which implements GetBasicDiskInfo method for getting a device details.
type Devinfo interface {
	GetBasicDiskInfo(attrName string) (string, error)
}

// Open returns error if a SCSI device returns error when opened
func (d *SCSIDev) Open() (err error) {
	d.fd, err = unix.Open(d.DevName, unix.O_RDWR, 0600)
	return err
}

// Close returns error if a SCSI device is not closed
func (d *SCSIDev) Close() error {
	return unix.Close(d.fd)
}

// DetectSCSIType returns the type of SCSI device
func DetectSCSIType(name string) (Dev, error) {
	device := SCSIDev{DevName: name}

	if err := device.Open(); err != nil {
		return nil, err
	}

	SCSIInquiry, err := device.SCSIInquiry()
	if err != nil {
		return nil, err
	}

	// Check if device is an ATA device (For an ATA device VendorIdentification value should be equal to ATA    )
	if SCSIInquiry.VendorID == [8]byte{0x41, 0x54, 0x41, 0x20, 0x20, 0x20, 0x20, 0x20} {
		return &SATA{device}, nil
	}

	return &device, nil
}

// GetBasicDiskInfo returns detail for a particular disk device attribute such as vendor,serial,etc
// Note : Now, it only returns basic disk info
func (d *SCSIDev) GetBasicDiskInfo(attrName string) (string, error) {
	// TODO : Fetch other disk attributes also such as serial no, wwn, etc
	// TODO : Return all the basic disk attributes available for a particular disk
	switch attrName {
	case Vendor:
		// Standard SCSI INQUIRY command
		InqRes, err := d.SCSIInquiry()
		if err != nil {
			return "", err
		}
		VendorID := InqRes.String()["VendorID"]
		return VendorID, nil
	case ProductDetail:
		// Standard SCSI INQUIRY command
		InqRes, err := d.SCSIInquiry()
		if err != nil {
			return "", err
		}
		ProductID := InqRes.String()["ProductID"]
		ProductRev := InqRes.String()["ProductRev"]
		ProductDetail := ProductID + " " + ProductRev
		return ProductDetail, nil
	case Capacity:
		capacity, err := d.ReadCapacity()
		if err != nil {
			return "", err
		}
		newCapacity := strconv.FormatUint(capacity, 10)
		return newCapacity, nil
	case RPM:
		// Get RPM using scsi mode pages
		RPM, err := d.GetScsiRPM()
		if err != nil {
			return "", err
		}
		return RPM, nil
	default:
		return "Attribute not found", nil

	}
}
