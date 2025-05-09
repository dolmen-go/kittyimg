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
	"embed"
	"image"
	"io"
	"os"
	"strings"
	"testing"

	// Plugin to decode GIF
	_ "image/gif"

	"github.com/dolmen-go/kittyimg"
)

//go:embed dolmen.gif
var files embed.FS

func Example() {
	f, err := files.Open("dolmen.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	kittyimg.Fprintln(os.Stdout, img)
}

func ExampleTranscode() {
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
	if !strings.HasPrefix(out, "\x1b_Gq=1,a=T,f=32,s=420,v=66,t=d,o=z,m=0;eJzsndGt") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestExampleTranscode(t *testing.T) {
	out := captureExampleOutput(t, "ExampleTranscode", ExampleTranscode)
	if !strings.HasPrefix(out, "\x1b_Gq=1,a=T,f=32,s=420,v=66,t=d,o=z,m=0;eJzsndGt") {
		t.Fatalf("unexpected output: %q", out)
	}
}
