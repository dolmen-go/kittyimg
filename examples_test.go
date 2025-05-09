/*
   Copyright 2021-2025 Olivier Mengué.

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
	"embed"
	"image"
	"io"
	"os"
	"strings"
	"testing"

	_ "image/gif"
	_ "image/png"

	"github.com/dolmen-go/kittyimg"
)

// Source: https://go.dev/play/p/XN6x3L23Vok
var favicon = []byte{
	0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n',
	0, 0, 0, 13, 'I', 'H', 'D', 'R',
	0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x0f, 0x04, 0x03, 0x00, 0x00, 0x00,
	0x1f, 0x5d, 0x52, 0x1c, // CRC
	0, 0, 0, 15, 'P', 'L', 'T', 'E',
	0x7a, 0xdf, 0xfd, 0xfd, 0xff, 0xfc, 0x39, 0x4d, 0x52, 0x19, 0x16, 0x15, 0xc3, 0x8d, 0x76,
	0xc7, 0x36, 0x2c, 0xf5, // CRC
	0, 0, 0, 64, 'I', 'D', 'A', 'T',
	0x08, 0xd7, 0x95, 0xc9, 0xd1, 0x0d, 0xc0, 0x20, 0x0c, 0x03, 0xd1, 0x23, 0x5d, 0xa0, 0x49, 0x17,
	0x20, 0x4c, 0xc0, 0x10, 0xec, 0x3f, 0x53, 0x8d, 0xc2, 0x02, 0x9c, 0xfc, 0xf1, 0x24, 0xe3, 0x31,
	0x54, 0x3a, 0xd1, 0x51, 0x96, 0x74, 0x1c, 0xcd, 0x18, 0xed, 0x9b, 0x9a, 0x11, 0x85, 0x24, 0xea,
	0xda, 0xe0, 0x99, 0x14, 0xd6, 0x3a, 0x68, 0x6f, 0x41, 0xdd, 0xe2, 0x07, 0xdb, 0xb5, 0x05, 0xca,
	0xdb, 0xb2, 0x9a, 0xdd, // CRC
	0, 0, 0, 0, 'I', 'E', 'N', 'D', 0xae, 0x42, 0x60, 0x82,
}

func Example() {
	f := bytes.NewReader(favicon)

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	kittyimg.Fprintln(os.Stdout, img)
}

func ExampleTranscode_png() {
	f := bytes.NewReader(favicon)

	kittyimg.Transcode(os.Stdout, f)
	os.Stdout.WriteString("\n")
}

//go:embed dolmen.gif
var files embed.FS

func ExampleTranscode_gif() {
	f, err := files.Open("dolmen.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	kittyimg.Transcode(os.Stdout, f)
	os.Stdout.WriteString("\n")
}

func captureExampleOutput(t *testing.T, name string, example func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("pipe:", err)
	}
	done := make(chan string)
	go func() {
		defer r.Close()
		var buf strings.Builder
		_, err = io.Copy(&buf, r)
		done <- buf.String()
	}()
	t.Run(name, func(t *testing.T) {
		origStdout := os.Stdout
		os.Stdout = w
		t.Cleanup(func() {
			os.Stdout = origStdout
			w.Close()
		})
		example()
	})
	out := <-done
	if err != nil { // Report copy error
		// t.Logf("%T %T", err, errors.Unwrap(err))
		t.Error("copy error:", err)
	}
	return out
}

func TestExample(t *testing.T) {
	out := captureExampleOutput(t, "Example", Example)
	t.Log(out)
	if !strings.HasPrefix(out, "\x1b_Gq=1,a=T,f=32,s=16,v=15,t=d,o=z,m=0;eJz6+//Pfxi29A") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestExampleTranscode_png(t *testing.T) {
	out := captureExampleOutput(t, "ExampleTranscode_png", ExampleTranscode_png)
	t.Log(out)
	// PNG file is directly transmitted
	if !strings.HasPrefix(out, "\x1b_Gq=1,a=T,f=100,s=16,v=15,m=0;iVBORw0KGgoA") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestExampleTranscode_gif(t *testing.T) {
	out := captureExampleOutput(t, "ExampleTranscode_gif", ExampleTranscode_gif)
	t.Log(out)
	if !strings.HasPrefix(out, "\x1b_Gq=1,a=T,f=32,s=420,v=66,t=d,o=z,m=0;eJzsndGt") {
		t.Fatalf("unexpected output: %q", out)
	}
}
