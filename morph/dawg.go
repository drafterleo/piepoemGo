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
	"encoding/base64"
	"os"
	"unicode/utf8"
)

const payloadSeparator = '\x01'

type dawg struct {
	dct   dictionary
	guide guide
}

func newDAWG(fn string) (*dawg, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d, err := newDictionary(f)
	if err != nil {
		return nil, err
	}

	g, err := newGuide(f)
	if err != nil {
		return nil, err
	}

	return &dawg{
		dct:   d,
		guide: g,
	}, nil
}

func b64d(p []byte) []byte {
	enc := base64.StdEncoding
	dst := make([]byte, enc.DecodedLen(len(p)))
	n, err := enc.Decode(dst, p)
	if err != nil {
		panic(err)
	}
	return dst[:n]
}

func (d *dawg) valuesForIndex(index uint32) [][]byte {
	var values [][]byte
	completer := newCompleter(d.dct, d.guide)
	completer.start(index, "")
	for completer.next() {
		values = append(values, b64d(completer.key))
	}
	return values
}

type item struct {
	key    string
	values [][]byte
}

func (d *dawg) similarItemsRecursive(prefix string, key []rune, index uint32) []item {
	var items []item

	startPos := utf8.RuneCountInString(prefix)
	endPos := len(key)
	wordPos := startPos

	for wordPos < endPos {
		r := key[wordPos]
		if r == 'е' {
			if next := d.dct.follow("ё", index); next != 0 {
				newPrefix := prefix + string(key[startPos:wordPos]) + "ё"
				items = append(items,
					d.similarItemsRecursive(newPrefix, key, next)...,
				)
			}
		}
		if index = d.dct.followRune(r, index); index == 0 {
			break
		}
		wordPos++
	}
	if wordPos == endPos {
		if index = d.dct.followByte(payloadSeparator, index); index != 0 {
			foundKey := prefix + string(key[startPos:])
			value := d.valuesForIndex(index)
			items = append([]item{{foundKey, value}}, items...)
		}
	}

	return items
}

func (d *dawg) similarItems(key string) []item {
	return d.similarItemsRecursive("", []rune(key), 0)
}
