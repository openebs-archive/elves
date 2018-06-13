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

// Functions for SCSI-ATA Translation.

package scsismart

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/openebs/elves/smart/pkg/atasmart"
	"github.com/openebs/elves/smart/pkg/utilities"
)

// SATA is a simple wrapper around an embedded SCSIDevice type, which handles sending ATA
// commands via SCSI pass-through (SCSI-ATA Translation).
type SATA struct {
	SCSIDev
}

// AtaIdentify sends SCSI_ATA_PASSTHRU_16 command and read data from the response
// received based on the defined ATA IDENTIFY STRUCT in commands.go file
func (d *SATA) AtaIdentify() (atasmart.IdentDevData, error) {
	var identifyBuf atasmart.IdentDevData

	responseBuf := make([]byte, 512)

	cdb16 := CDB16{SCSIATAPassThru}
	cdb16[1] = 0x08 // ATA protocol (4 << 1, PIO data-in)
	cdb16[2] = 0x0e // BYT_BLOK = 1, T_LENGTH = 2, T_DIR = 1
	cdb16[14] = atasmart.AtaIdentifyDevice

	if err := d.sendSCSICDB(cdb16[:], &responseBuf); err != nil {
		return identifyBuf, fmt.Errorf("sendSCSICDB ATA IDENTIFY: %v", err)
	}

	binary.Read(bytes.NewBuffer(responseBuf), utilities.NativeEndian, &identifyBuf)

	return identifyBuf, nil
}

// GetBasicDiskInfo returns all the disk attributes and smart info for a particular SATA device
func (d *SATA) GetBasicDiskInfo(attrName string) (string, error) {
	// store data from the response received by calling AtaIdentify based on the defined ATA IDENTIFY STRUCT
	identifyBuf, err := d.AtaIdentify()
	if err != nil {
		return "", err
	}
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
	case SerialNumber:
		SerialNumber := string(identifyBuf.GetSerialNumber())
		return SerialNumber, nil
	case ModelNumber:
		ModelNumber := string(identifyBuf.GetModelNumber())
		return ModelNumber, nil
	case WWN:
		LuWWNDeviceID := identifyBuf.GetWWN()
		return LuWWNDeviceID, nil
	case FirmwareRev:
		FirmwareRev := string(identifyBuf.GetFirmwareRevision())
		return FirmwareRev, nil
	case AtATransport:
		AtATransport := string(identifyBuf.AtaTransport())
		return AtATransport, nil
	case ATAMajor:
		ATAMajor := identifyBuf.GetATAMajorVersion()
		return ATAMajor, nil
	case ATAMinor:
		ATAMinor := identifyBuf.GetATAMinorVersion()
		return ATAMinor, nil
	case Capacity:
		capacity, err := d.ReadCapacity()
		if err != nil {
			return "", err
		}
		newCapacity := strconv.FormatUint(capacity, 10)
		return newCapacity, nil
	case RPM:
		RPM := identifyBuf.RotationRate
		Rpm := strconv.FormatInt(int64(RPM), 10)
		return Rpm, nil
	case LogicalSize:
		LBSize, _ := identifyBuf.GetSectorSize()
		LogicalSize := strconv.FormatInt(int64(LBSize), 10)
		return LogicalSize, nil
	case PhysicalSize:
		_, PBSize := identifyBuf.GetSectorSize()
		PhysicalSize := strconv.FormatInt(int64(PBSize), 10)
		return PhysicalSize, nil
	default:
		return "Attribute Details not found", nil

	}
}
