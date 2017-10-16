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
	"reflect"
	"testing"
)

var testCases = []struct {
	word string
	want [3][]string
}{
	{"олакрез", [3][]string{nil, nil, nil}},
	{"еж", [3][]string{
		{"ёж", "ёж", "ёж"},
		{"ёж", "ёж", "ёж"},
		{"NOUN,anim,masc sing,nomn", "NOUN,inan,masc sing,nomn", "NOUN,inan,masc sing,accs"},
	}},
	{"черта", [3][]string{
		{"черта", "чёрта", "чёрта"},
		{"черта", "чёрт", "чёрт"},
		{"NOUN,inan,femn sing,nomn", "NOUN,anim,masc sing,gent", "NOUN,anim,masc sing,accs"},
	}},
	{"чёрт", [3][]string{
		{"чёрт"},
		{"чёрт"},
		{"NOUN,anim,masc sing,nomn"},
	}},
	{"черт", [3][]string{
		{"черт", "чёрт"},
		{"черта", "чёрт"},
		{"NOUN,inan,femn plur,gent", "NOUN,anim,masc sing,nomn"},
	}},
	{"изба", [3][]string{
		{"изба"},
		{"изба"},
		{"NOUN,inan,femn sing,nomn"},
	}},
	{"квакании", [3][]string{
		{"квакании"},
		{"кваканье"},
		{"NOUN,inan,neut sing,loct,V-ie"},
	}},
	{"залом", [3][]string{
		{"залом", "залом", "залом", "залом", "залом"},
		{"залом", "залом", "зал", "залом", "зало"},
		{"NOUN,inan,masc sing,nomn", "NOUN,inan,masc sing,accs", "NOUN,inan,masc sing,ablt", "NOUN,anim,masc sing,nomn", "NOUN,inan,neut sing,ablt"},
	}},
}

func TestParse(t *testing.T) {
	for _, tc := range testCases {
		words, norms, tags := Parse(tc.word)
		if !reflect.DeepEqual(words, tc.want[0]) {
			t.Errorf("Parse(%q): want words %v, got %v", tc.word, tc.want[0], words)
		}
		if !reflect.DeepEqual(norms, tc.want[1]) {
			t.Errorf("Parse(%q): want norms %v, got %v", tc.word, tc.want[1], norms)
		}
		if !reflect.DeepEqual(tags, tc.want[2]) {
			t.Errorf("Parse(%q): want tags %v, got %v", tc.word, tc.want[2], tags)
		}
	}
}
