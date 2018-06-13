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

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/openebs/elves/smart/pkg/ioctl"
	"github.com/openebs/elves/smart/pkg/utilities"
)

// SCSI commands being used
const (
	SCSIInquiry      = 0x12
	SCSIModeSense    = 0x1a
	SCSIReadCapacity = 0x25
	SCSIATAPassThru  = 0x85

	// Minimum length of standard INQUIRY response
	INQRespLen = 36
)

// SCSI Command Descriptor Block types
type CDB6 [6]byte
type CDB10 [10]byte
type CDB16 [16]byte

// runSCSIGen executes SCSI generic commands i.e sgIO commands
func (d *SCSIDev) runSCSIGen(header *sgIOHeader) error {
	if err := ioctl.Ioctl(uintptr(d.fd), SGIO, uintptr(unsafe.Pointer(header))); err != nil {
		return err
	}

	// See http://www.t10.org/lists/2status.htm for SCSI status codes
	if header.info&SGInfoOkMask != SGInfoOk {
		err := sgIOErr{
			scsiStatus:   header.status,
			hostStatus:   header.hostStatus,
			driverStatus: header.driverStatus,
		}
		return err
	}

	return nil
}

func (e sgIOErr) Error() string {
	return fmt.Sprintf("SCSI status: %#02x, host status: %#02x, driver status: %#02x",
		e.scsiStatus, e.hostStatus, e.driverStatus)
}

// sendSCSICDB sends a SCSI Command Descriptor Block to the device and writes the response into the
// supplied []byte pointer.
func (d *SCSIDev) sendSCSICDB(cdb []byte, respBuf *[]byte) error {
	senseBuf := make([]byte, 32)

	// Populate all the required fields of "sg_io_hdr_t" struct
	header := sgIOHeader{
		interfaceID:    'S',
		dxferDirection: SGDxferFromDev,
		timeout:        DefaultTimeout,
		cmdLen:         uint8(len(cdb)),
		mxSBLen:        uint8(len(senseBuf)),
		dxferLen:       uint32(len(*respBuf)),
		dxferp:         uintptr(unsafe.Pointer(&(*respBuf)[0])),
		cmdp:           uintptr(unsafe.Pointer(&cdb[0])),
		sbp:            uintptr(unsafe.Pointer(&senseBuf[0])),
	}

	return d.runSCSIGen(&header)
}

// SCSIInquiry sends an INQUIRY command to a SCSI device and returns an InquiryResponse struct.
func (d *SCSIDev) SCSIInquiry() (InquiryResponse, error) {
	var response InquiryResponse

	respBuf := make([]byte, INQRespLen)

	cdb := CDB6{SCSIInquiry}
	binary.BigEndian.PutUint16(cdb[3:], uint16(len(respBuf)))

	if err := d.sendSCSICDB(cdb[:], &respBuf); err != nil {
		return response, err
	}

	binary.Read(bytes.NewBuffer(respBuf), utilities.NativeEndian, &response)

	return response, nil
}

func (inquiry InquiryResponse) String() map[string]string {
	InqRespMap := make(map[string]string)
	InqRespMap["VendorID"] = fmt.Sprintf("%.8s", inquiry.VendorID)
	InqRespMap["ProductID"] = fmt.Sprintf("%.16s", inquiry.ProductID)
	InqRespMap["ProductRev"] = fmt.Sprintf("%.4s", inquiry.ProductRev)
	return InqRespMap
}

// modeSense function is used to send a SCSI MODE SENSE(6) command to a device.
// TODO : Implement SCSI MODE SENSE(10) command also
func (d *SCSIDev) modeSense(pageNo, subPageNo, pageCtrl uint8) ([]byte, error) {
	respBuf := make([]byte, 64)

	cdb := CDB6{SCSIModeSense}
	cdb[2] = (pageCtrl << 6) | (pageNo & 0x3f)
	cdb[3] = subPageNo
	cdb[4] = uint8(len(respBuf))

	if err := d.sendSCSICDB(cdb[:], &respBuf); err != nil {
		return respBuf, err
	}

	return respBuf, nil
}

// ReadCapacity sends a SCSI READ CAPACITY(10) command to a device and returns the capacity in bytes.
func (d *SCSIDev) ReadCapacity() (uint64, error) {
	respBuf := make([]byte, 8)
	cdb := CDB10{SCSIReadCapacity}

	if err := d.sendSCSICDB(cdb[:], &respBuf); err != nil {
		return 0, err
	}

	lastLBA := binary.BigEndian.Uint32(respBuf[0:]) // max. addressable LBA
	LBsize := binary.BigEndian.Uint32(respBuf[4:])  // logical block size
	capacity := (uint64(lastLBA) + 1) * uint64(LBsize)

	return capacity, nil

}
