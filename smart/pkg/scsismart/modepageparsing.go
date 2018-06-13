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

// Various SCSI mode pages

package scsismart

import (
	"encoding/binary"
)

const (
	// SCSI mode pages code
	RigidDiskGeometryPage = 0x04

	// Mode page control field default value
	ModePageControlDefault = 2
)

// GetScsiRPM is used to fetch the Rotation per minute (RPM) of a scsi device using scsi mode page
// i.e., Rigid disk geometry page
// This is just a sample usage of scsi mode sense command to parse scsi mode pages
// WIP ...
func (d *SCSIDev) GetScsiRPM() (string, error) {
	//  Getting the response from Rigid Disk Geometry scsi mode page
	DiskGeomPageResp, err := d.modeSense(RigidDiskGeometryPage, 0, ModePageControlDefault)
	if err != nil {
		return "", err
	}
	bdLen := DiskGeomPageResp[3]
	offset := bdLen + 4
	RPM := string(binary.BigEndian.Uint16(DiskGeomPageResp[offset+20:]))

	return RPM, err
}
