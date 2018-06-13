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

// Definition for ATA commands

package atasmart

// Table 10 of X3T13/2008D (ATA-3) Revision 7b, See http://www.scs.stanford.edu/11wi-cs140/pintos/specs/ata-3-std.pdf
// Table 31 of T13/1699-D Revision 6a, See http://www.t13.org/documents/uploadeddocuments/docs2008/d1699r6a-ata8-acs.pdf
// Table 47 of T13/2161-D Revision 5, See http://www.t13.org/Documents/UploadedDocuments/docs2013/d2161r5-ATAATAPI_Command_Set_-_3.pdf
// Table 57 of T13/BSR INCITS 529 Revision 18 , See http://t13.org/Documents/UploadedDocuments/docs2017/di529r18-ATAATAPI_Command_Set_-_4.pdf
var ataMinorVersions = map[uint16]string{
	0x0001: "ATA-1 X3T9.2/781D prior to revision 4",       // obsolete
	0x0002: "ATA-1 published, ANSI X3.221-1994",           // obsolete
	0x0003: "ATA-1 X3T9.2/781D revision 4",                // obsolete
	0x0004: "ATA-2 published, ANSI X3.279-1996",           // obsolete
	0x0005: "ATA-2 X3T10/948D prior to revision 2k",       // obsolete
	0x0006: "ATA-3 X3T10/2008D revision 1",                // obsolete
	0x0007: "ATA-2 X3T10/948D revision 2k",                // obsolete
	0x0008: "ATA-3 X3T10/2008D revision 0",                // obsolete
	0x0009: "ATA-2 X3T10/948D revision 3",                 // obsolete
	0x000a: "ATA-3 published, ANSI X3.298-1997",           // obsolete
	0x000b: "ATA-3 X3T10/2008D revision 6",                // obsolete
	0x000c: "ATA-3 X3T13/2008D revision 7 and 7a",         // obsolete
	0x000d: "ATA/ATAPI-4 X3T13/1153D version 6",           // obsolete
	0x000e: "ATA/ATAPI-4 T13/1153D version 13",            // obsolete
	0x000f: "ATA/ATAPI-4 X3T13/1153D version 7",           // obsolete
	0x0010: "ATA/ATAPI-4 T13/1153D version 18",            // obsolete
	0x0011: "ATA/ATAPI-4 T13/1153D version 15",            // obsolete
	0x0012: "ATA/ATAPI-4 published, ANSI NCITS 317-1998",  // obsolete
	0x0013: "ATA/ATAPI-5 T13/1321D version 3",             // obsolete
	0x0014: "ATA/ATAPI-4 T13/1153D version 14",            // obsolete
	0x0015: "ATA/ATAPI-5 T13/1321D revision 1",            // obsolete
	0x0016: "ATA/ATAPI-5 published, ANSI NCITS 340-2000",  // obsolete
	0x0017: "ATA/ATAPI-4 T13/1153D revision 17",           // obsolete
	0x0018: "ATA/ATAPI-6 T13/1410D version 0",             // obsolete
	0x0019: "ATA/ATAPI-6 T13/1410D version 3a",            // obsolete
	0x001a: "ATA/ATAPI-7 T13/1532D version 1",             // obsolete
	0x001b: "ATA/ATAPI-6 T13/1410D version 2",             // obsolete
	0x001c: "ATA/ATAPI-6 T13/1410D version 1",             // obsolete
	0x001d: "ATA/ATAPI-7 published, ANSI INCITS 397-2005", // obsolete
	0x001e: "ATA/ATAPI-7 T13/1532D version 0",             // obsolete
	0x001f: "ACS-3 revision 3b",
	0x0021: "ATA/ATAPI-7 T13/1532D version 4a",            // obsolete
	0x0022: "ATA/ATAPI-6 published, ANSI INCITS 361-2002", // obsolete
	0x0027: "ATA8-ACS version 3c",
	0x0028: "ATA8-ACS version 6",
	0x0029: "ATA8-ACS version 4",
	0x0031: "ACS-2 revision 2",
	0x0033: "ATA8-ACS version 3e",
	0x0039: "ATA8-ACS version 4c",
	0x0042: "ATA8-ACS version 3f",
	0x0052: "ATA8-ACS version 3b",
	0x005e: "ACS-4 revision 5",
	0x006d: "ACS-3 revision 5",
	0x0082: "ACS-2 published, ANSI INCITS 482-2012",
	0x0107: "ATA8-ACS version 2d",
	0x010a: "ACS-3 published, ANSI INCITS 522-2014",
	0x0110: "ACS-2 revision 3",
	0x011b: "ACS-3 revision 4",
}

// IdentDevData struct is an ATA IDENTIFY DEVICE struct. ATA8-ACS defines this as a page of 16-bit words.
// _ (underscore) is used here to skip the words which we don't want to parse or get the data while parsing
// the ata identify device data struct page.
type IdentDevData struct {
	_                 [10]uint16  // ...
	SerialNumber      [20]byte    // Word 10..19, device serial number.
	_                 [3]uint16   // ...
	FirmwareRev       [8]byte     // Word 23..26, device firmware revision.
	ModelNumber       [40]byte    // Word 27..46, device model number.
	_                 [33]uint16  // ...
	MajorVer          uint16      // Word 80, major version number.
	MinorVer          uint16      // Word 81, minor version number.
	_                 [24]uint16  // ...
	SectorSize        uint16      // Word 106, Logical/physical sector size.
	_                 [1]uint16   // ...
	WWN               [4]uint16   // Word 108..111, WWN (World Wide Name).
	_                 [105]uint16 // ...
	RotationRate      uint16      // Word 217, nominal media rotation rate.
	_                 [4]uint16   // ...
	AtaTransportMajor uint16      // Word 222, Transport major version number.
	_                 [33]uint16  // ...
} // 512 bytes
