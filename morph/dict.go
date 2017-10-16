// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

package morph

import (
	"encoding/binary"
	"io"
	"unicode/utf8"
)

const (
	isLeafBit    = 1 << 31
	hasLeafBit   = 1 << 8
	extensionBit = 1 << 9
)

type dictionary []uint32

func newDictionary(r io.Reader) (dictionary, error) {
	var size uint32
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return nil, err
	}

	d := make(dictionary, size)
	if err := binary.Read(r, binary.LittleEndian, d); err != nil {
		return nil, err
	}

	return d, nil
}

func hasLeaf(n uint32) bool {
	return n&hasLeafBit != 0
}

func label(n uint32) uint32 {
	return n & (isLeafBit | 0xff)
}

func offset(n uint32) uint32 {
	return (n >> 10) << ((n & extensionBit) >> 6)
}

func (d dictionary) hasValue(index uint32) bool {
	return hasLeaf(d[index])
}

func (d dictionary) followByte(lbl byte, index uint32) uint32 {
	off := offset(d[index])
	next := index ^ off ^ uint32(lbl)
	if label(d[next]) != uint32(lbl) {
		return 0
	}
	return next
}

func (d dictionary) follow(s string, index uint32) uint32 {
	for i := 0; i < len(s); i++ {
		index = d.followByte(s[i], index)
		if index == 0 {
			return 0
		}
	}
	return index
}

func (d dictionary) followRune(r rune, index uint32) uint32 {
	buf := make([]byte, 4)
	n := utf8.EncodeRune(buf, r)
	for i := 0; i < n; i++ {
		index = d.followByte(buf[i], index)
		if index == 0 {
			return 0
		}
	}
	return index
}
