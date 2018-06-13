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

// SCSI command definitions.

package scsismart

// SCSI generic (sg)
// See dxfer_direction http://sg.danny.cz/sg/p/sg_v3_ho.html
const (
	SGDxferNone      = -1 //SCSI Test Unit Ready command
	SGDxferToDev     = -2 //SCSI WRITE command
	SGDxferFromDev   = -3 //SCSI READ command
	SGDxferToFromDev = -4

	SGInfoOkMask = 0x1
	SGInfoOk     = 0x0

	SGIO = 0x2285

	// DefaultTimeout in millisecs
	DefaultTimeout = 20000
)

// Constants being used by switch case for returning disk details
const (
	Vendor        = "Vendor"
	ProductDetail = "ProductDetail"
	Capacity      = "Capacity"
	LogicalSize   = "LogicalSize"
	PhysicalSize  = "PhysicalSize"
	SerialNumber  = "SerialNumber"
	WWN           = "LuWWNDeviceID"
	FirmwareRev   = "FirmwareRevision"
	ModelNumber   = "ModelNumber"
	RPM           = "RPM"
	ATAMajor      = "ATAMajorVersion"
	ATAMinor      = "ATAMinorVersion"
	AtATransport  = "AtaTransport"
)

// SCSIDev structure
type SCSIDev struct {
	DevName string // SCSI device name
	fd      int    // File descriptor for the scsi device
}

// sg_io_hdr_t structure See http://sg.danny.cz/sg/p/sg_v3_ho.html
type sgIOHeader struct {
	interfaceID    int32   // 'S' for SCSI generic (required)
	dxferDirection int32   // data transfer direction
	cmdLen         uint8   // SCSI command length (<= 16 bytes)
	mxSBLen        uint8   // max length to write to sbp
	iovecCount     uint16  // 0 implies no scatter gather
	dxferLen       uint32  // byte count of data transfer
	dxferp         uintptr // points to data transfer memory or scatter gather list
	cmdp           uintptr // points to command to perform
	sbp            uintptr // points to sense_buffer memory
	timeout        uint32  // MAX_UINT -> no timeout (unit: millisec)
	flags          uint32  // 0 -> default, see SG_FLAG...
	packID         int32   // unused internally (normally)
	usrPtr         uintptr // unused internally
	status         uint8   // SCSI status
	maskedStatus   uint8   // shifted, masked scsi status
	msgStatus      uint8   // messaging level data (optional)
	SBLenwr        uint8   // byte count actually written to sbp
	hostStatus     uint16  // errors from host adapter
	driverStatus   uint16  // errors from software driver
	resid          int32   // dxfer_len - actual_transferred
	duration       uint32  // time taken by cmd (unit: millisec)
	info           uint32  // auxiliary information
}

type sgIOErr struct {
	scsiStatus   uint8
	hostStatus   uint16
	driverStatus uint16
	senseBuf     [32]byte
}

/* These structs are not in use ------------------------------------------------------
// DiskAttr is struct being used for returning all the available disk details (both basic and smart)
// For now, only basic disk attr are being fetched so it is returning only basic attrs
type DiskAttr struct {
	BasicDiskAttr
}

// BasicDiskAttr is the structure being used for returning basic disk details
type BasicDiskAttr struct {
	SCSIInquiry      InquiryResponse
	UserCapacity     uint64
	LBSize           uint16
	PBSize           uint16
	SerialNumber     string
	LuWWNDeviceID    string
	FirmwareRevision string
	ModelNumber      string
	RotationRate     uint16
	ATAMajorVersion  string
	ATAMinorVersion  string
	AtaTransport     string
}

// SmartDiskAttr is the structure defined for smart disk attrs (Note : Not being used yet)
type SmartDiskAttr struct {
}

*/

// InquiryResponse is the struct for SCSI INQUIRY response
type InquiryResponse struct {
	_          [2]byte // type of device
	Version    byte
	_          [5]byte
	VendorID   [8]byte  // Vendor Identification
	ProductID  [16]byte // Product Identification
	ProductRev [4]byte  // Product Revision Level
}
