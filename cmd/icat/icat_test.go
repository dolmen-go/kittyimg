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

package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func runMain(t *testing.T, name string, args ...string) (string, error) {
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

	var mainErr error
	t.Run(name, func(t *testing.T) {
		origStdout := os.Stdout
		os.Stdout = w
		t.Cleanup(func() {
			os.Stdout = origStdout
			w.Close()
		})

		mainErr = icatMain(w, args)
	})

	out := <-done
	if err != nil { // Report copy error
		// t.Logf("%T %T", err, errors.Unwrap(err))
		t.Error("copy error:", err)
	}
	return out, mainErr
}

func Test(t *testing.T) {
	out, err := runMain(t, "icat dolmen.gif", "../../dolmen.gif", "../../dolmen.gif")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Output:\n" + out)
	t.Logf("%q", out)
	if !strings.HasPrefix(out, "\x1b_Gq=1,a=T,f=32,s=420,v=66,t=d,o=z;eJzsndGt") ||
		!strings.HasSuffix(out, "9yLYll\x1b\\\n") {
		t.Error("unexpected output")
	}
}
