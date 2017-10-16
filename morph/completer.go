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

type completer struct {
	dict       dictionary
	guide      guide
	lastIndex  uint32
	indexStack []uint32
	key        []byte
}

func newCompleter(d dictionary, g guide) *completer {
	return &completer{
		dict:  d,
		guide: g,
	}
}

func (c *completer) start(index uint32, prefix string) {
	c.key = []byte(prefix)
	if len(c.guide) > 0 {
		c.indexStack = []uint32{index}
	}
}

func (c *completer) next() bool {
	if len(c.indexStack) == 0 {
		return false
	}

	index := c.indexStack[len(c.indexStack)-1]

	if c.lastIndex != 0 {
		for {
			siblingLabel := c.guide.sibling(index)
			if len(c.key) > 0 {
				c.key = c.key[:len(c.key)-1]
			}

			c.indexStack = c.indexStack[:len(c.indexStack)-1]
			if len(c.indexStack) == 0 {
				return false
			}

			index = c.indexStack[len(c.indexStack)-1]
			if siblingLabel != 0 {
				if index = c.follow(siblingLabel, index); index == 0 {
					return false
				}
				break
			}
		}
	}

	return c.findTerminal(index)
}

func (c *completer) follow(label byte, index uint32) uint32 {
	index = c.dict.followByte(label, index)
	if index == 0 {
		return 0
	}

	c.key = append(c.key, label)
	c.indexStack = append(c.indexStack, index)
	return index
}

func (c *completer) findTerminal(index uint32) bool {
	for !c.dict.hasValue(index) {
		label := c.guide.child(index)
		if index = c.dict.followByte(label, index); index == 0 {
			return false
		}
		c.key = append(c.key, label)
		c.indexStack = append(c.indexStack, index)
	}
	c.lastIndex = index
	return true
}
