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

package atasmart

import (
	"fmt"

	"github.com/openebs/elves/smart/pkg/utilities"
)

// swapByteOrder swaps the order of every second byte in a byte slice and modifies the slice in place.
func (d *IdentDevData) swapByteOrder(b []byte) []byte {
	tmp := make([]byte, len(b))

	for i := 0; i < len(b); i += 2 {
		tmp[i], tmp[i+1] = b[i+1], b[i]
	}
	return tmp
}

// GetSerialNumber returns the serial number of a disk device using an ATA IDENTIFY command.
func (d *IdentDevData) GetSerialNumber() []byte {
	return d.swapByteOrder(d.SerialNumber[:])
}

// GetModelNumber returns the model number of a disk device using an ATA IDENTIFY command.
func (d *IdentDevData) GetModelNumber() []byte {
	return d.swapByteOrder(d.ModelNumber[:])
}

// GetFirmwareRevision returns the firmware revision of a disk device using an ATA IDENTIFY command.
func (d *IdentDevData) GetFirmwareRevision() []byte {
	return d.swapByteOrder(d.FirmwareRev[:])
}

// GetWWN returns the worldwide unique name for a disk
func (d *IdentDevData) GetWWN() string {
	NAA := d.WWN[0] >> 12
	IEEEOUI := (uint32(d.WWN[0]&0x0fff) << 12) | (uint32(d.WWN[1]) >> 4)
	UniqueID := ((uint64(d.WWN[1]) & 0xf) << 32) | (uint64(d.WWN[2]) << 16) | uint64(d.WWN[3])

	return fmt.Sprintf("%x %06x %09x", NAA, IEEEOUI, UniqueID)
}

// GetSectorSize returns logical and physical sector sizes of a disk
func (d *IdentDevData) GetSectorSize() (uint16, uint16) {
	// By default, we are assuming the physical and logical sector size to be 512
	// based on further check conditions, it would be altered.
	var LogSec, PhySec uint16 = 512, 512

	if (d.SectorSize & 0xc000) != 0x4000 {
		return LogSec, PhySec
	}
	// TODO : Add support for Long Logical/Physical Sectors (LLS/LPS)
	if (d.SectorSize & 0x2000) != 0x0000 {
		// Physical sector size is multiple of logical sector size
		PhySec <<= (d.SectorSize & 0x0f)
	}
	return LogSec, PhySec
}

// GetATAMajorVersion returns the ATA major version of a disk using an ATA IDENTIFY command.
func (d *IdentDevData) GetATAMajorVersion() (s string) {
	if (d.MajorVer == 0) || (d.MajorVer == 0xffff) {
		s = "This device does not report ATA major version"
		return
	}
	switch utilities.MSignificantBit(uint(d.MajorVer)) {
	case 1:
		s = "ATA-1"
	case 2:
		s = "ATA-2"
	case 3:
		s = "ATA-3"
	case 4:
		s = "ATA/ATAPI-4"
	case 5:
		s = "ATA/ATAPI-5"
	case 6:
		s = "ATA/ATAPI-6"
	case 7:
		s = "ATA/ATAPI-7"
	case 8:
		s = "ATA8-ACS"
	case 9:
		s = "ACS-2"
	case 10:
		s = "ACS-3"
	}

	return
}

// GetATAMinorVersion returns the ATA minor version using an ATA IDENTIFY command.
func (d *IdentDevData) GetATAMinorVersion() string {
	if (d.MinorVer == 0) || (d.MinorVer == 0xffff) {
		return "This device does not report ATA minor version"
	}
	// Since ATA minor version word is not a bitmask, we simply do a map lookup here
	if s, ok := ataMinorVersions[d.MinorVer]; ok {
		return s
	}
	return "unknown"
}

// AtaTransport returns the type of ata Transport being used such as serial ATA, parallel ATA, etc.
func (d *IdentDevData) AtaTransport() (s string) {
	if (d.AtaTransportMajor == 0) || (d.AtaTransportMajor == 0xffff) {
		s = "This device does not report Transport"
		return
	}

	switch d.AtaTransportMajor >> 12 {
	case 0x0:
		s = "Parallel ATA"
	case 0x1:
		s = "Serial ATA"

		switch utilities.MSignificantBit(uint(d.AtaTransportMajor & 0x0fff)) {
		case 0:
			s += " ATA8-AST"
		case 1:
			s += " SATA 1.0a"
		case 2:
			s += " SATA II Ext"
		case 3:
			s += " SATA 2.5"
		case 4:
			s += " SATA 2.6"
		case 5:
			s += " SATA 3.0"
		case 6:
			s += " SATA 3.1"
		case 7:
			s += " SATA 3.2"
		default:
			s += fmt.Sprintf(" SATA (%#03x)", d.AtaTransportMajor&0x0fff)
		}
	case 0xe:
		s = fmt.Sprintf("PCIe (%#03x)", d.AtaTransportMajor&0x0fff)
	default:
		s = fmt.Sprintf("Unknown (%#04x)", d.AtaTransportMajor)
	}

	return
}
