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

// Miscellaneous utility functions

package utilities

import (
	"encoding/binary"
	"math/bits"
	"unsafe"
)

var (
	NativeEndian binary.ByteOrder
)

// init determines native endianness of a system
func init() {
	i := uint32(1)
	b := (*[4]byte)(unsafe.Pointer(&i))
	if b[0] == 1 {
		NativeEndian = binary.LittleEndian
	} else {
		NativeEndian = binary.BigEndian
	}
}

// MSignificantBit finds the most significant bit set in a uint
func MSignificantBit(i uint) int {
	if i == 0 {
		return 0
	}

	return bits.Len(i) - 1
}
