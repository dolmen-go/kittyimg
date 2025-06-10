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
	"os"
	"strings"
	"testing"

	_ "image/gif"

	"github.com/dolmen-go/kittyimg"
)

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

func TestExampleTranscode_gif(t *testing.T) {
	out := captureExampleOutput(t, "ExampleTranscode_gif", ExampleTranscode_gif)
	t.Log(out)
	if !strings.HasPrefix(out, "\x1b_Gq=1,a=T,f=32,s=420,v=66,t=d,o=z,m=0;eJzsndGt") {
		t.Fatalf("unexpected output: %q", out)
	}
}
