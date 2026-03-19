//go:build go1.24

/*
   Copyright 2021-2026 Olivier Mengué.

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

package kittyimg_test

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"iter"
	"maps"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/dolmen-go/kittyimg"
)

const (
	// Enforce canonical Base64
	base64CharRE  = `[A-Za-z0-9+/]`
	base64CharRE2 = `[AQgw]`             // See https://go.dev/play/p/ui8tmhV-YLH
	base64CharRE3 = `[AEIMQUYcgkosw048]` // See https://go.dev/play/p/HVF_A6wJcOo
	base64REInf   = `(?:` + base64CharRE + `{4})*(?:` + base64CharRE + "(?:" + base64CharRE2 + `=|` + base64CharRE + base64CharRE3 + `)=)?`

	// Payload is limited to 4096 bytes
	// 4096 / 4 = 1024 blocks
	// But package regexp has a limit of 250 repetitions (see issue #78222), so we have to split:
	// 1024 = 250*4 + 23 + 1
	base64Chars4   = `(?:` + base64CharRE + `{4})`
	base64Chars250 = base64Chars4 + `{0,250}`
	payloadRE      = base64Chars250 + base64Chars250 + base64Chars250 + base64Chars250 + base64Chars4 + `{0,23}(?:` + base64CharRE + "(?:" + base64CharRE2 + `==|` + base64CharRE + `(?:` + base64CharRE3 + `=|` + base64CharRE + `{2})))?`
)

var (
	blockRE = regexp.MustCompile(`` +
		"\033_G" +
		"(?<params>(?:[a-zA-Z]=[^,;\033]{1,11}(?:,[a-zA-Z]=[^,;\033]{1,11})*)?)" +
		"(?:;(?<payload>" + payloadRE + ")?)" +
		"\033\\\\",
	)
	kvRE = regexp.MustCompile(`^(?:([a-zA-Z])=([^,;\033]+)(?:,|$))+$`)
)

func countBlocks(s string) int {
	m := blockRE.FindAllStringIndex(s, -1)
	if m == nil {
		return 0
	}
	return len(m)
}

type Params map[byte]string

func (p Params) String() string {
	if p == nil || len(p) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, k := range slices.Sorted(maps.Keys(p)) {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte(k)
		sb.WriteByte('=')
		sb.WriteString(p[k])
	}
	return sb.String()
}

type Block struct {
	Params  Params
	Payload []byte
}

func addParam(params *Params, p string) {
	if *params == nil {
		*params = Params{p[0]: p[2:]}
		return
	}
	(*params)[p[0]] = p[2:]
}

func parseParams(p string) (params Params) {
	if p == "" {
		return nil
	}
	for {
		i := strings.IndexByte(p, ',')
		if i == -1 {
			addParam(&params, p)
			break
		}
		addParam(&params, p[:i])
		p = p[i+1:]
	}
	return
}

func (bl *Block) init(params string, payload string) {
	bl.Params = parseParams(params)
	if payload != "" {
		bl.Payload, _ = base64.StdEncoding.DecodeString(payload)
	}
}

func (bl *Block) UnmarshalText(b []byte) error {
	m := blockRE.FindSubmatchIndex(b)
	if m == nil || m[0] != 0 || m[1] != len(b) {
		return errors.New("invalid block")
	}
	bl.init(string(b[m[2]:m[3]]), string(b[m[4]:m[5]]))
	return nil
}

func extractBlocks(s []byte) iter.Seq[*Block] {
	matches := blockRE.FindAllSubmatchIndex(s, -1)
	if len(matches) == 0 {
		// panic("no match")
		return func(yield func(*Block) bool) {}
	}
	anchor := 0
	for _, m := range matches {
		if m[0] != anchor {
			panic("should match contiguously")
		}
		anchor = m[1]
	}
	if anchor != len(s) {
		panic(fmt.Errorf("should match full string but found %q", s[anchor:]))
	}
	return func(yield func(*Block) bool) {
		for i, m := range matches {
			matches[i] = nil // free memory early
			//fmt.Printf("Params: %s\n", s[m[2]:m[3]])
			//if len(m) > 2 {
			//	fmt.Printf("Payload: %q\n", s[m[4]:m[5]])
			//}

			var bl Block
			/*
				err := bl.UnmarshalText([]byte(s[m[0]:m[1]]))
				if err != nil {
					panic(err)
				}
			*/
			bl.init(string(s[m[2]:m[3]]), string(s[m[4]:m[5]]))

			if !yield(&bl) {
				break
			}
		}
		return
	}
}

func testDecode(t *testing.T, filepath string, expectedLen int) {
	f, err := os.Open(filepath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	if err := kittyimg.Transcode(&buf, f); err != nil {
		t.Fatal(err)
	}

	gotLen := buf.Len()
	t.Logf("Output: %d bytes", gotLen)

	i := 0
	for bl := range extractBlocks(buf.Bytes()) {
		t.Log("-- Block", i, "--")
		t.Log("Params:", bl.Params)
		t.Logf("Payload: %d bytes", len(bl.Payload))
		i++
	}

	if i == 0 {
		t.Fatal("Decode failure!")
	}
	if gotLen != expectedLen {
		t.Fatalf("Length: got %d, expected %d", gotLen, expectedLen)
	}
}

func TestImagePNG3069(t *testing.T) {
	testDecode(t, "testdata/go-favicon-3069.png", 4121)
}

func TestImagePNG3070(t *testing.T) {
	testDecode(t, "testdata/go-favicon-3070.png", 4125)
}

func TestImagePNG3071(t *testing.T) {
	testDecode(t, "testdata/go-favicon-3071.png", 4125)
}

// Test encoding of a PNG image file of 3072 bytes, which is a base64 payload of 4096.
func TestImagePNG3072(t *testing.T) {
	testDecode(t, "testdata/go-favicon-3072.png", 4125)
}

// Test encoding of a PNG image file of 3073 bytes, which is two blocks (3072+1 => 4096+4).
func TestImagePNG3073(t *testing.T) {
	testDecode(t, "testdata/go-favicon-3073.png", 4139)
}
