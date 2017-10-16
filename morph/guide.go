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
)

type guide []byte

func newGuide(r io.Reader) (guide, error) {
	var size uint32
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return nil, err
	}

	g := make(guide, size*2)
	if err := binary.Read(r, binary.LittleEndian, g); err != nil {
		return nil, err
	}

	return g, nil
}

func (g guide) child(n uint32) byte {
	return g[n*2]
}

func (g guide) sibling(n uint32) byte {
	return g[n*2+1]
}
