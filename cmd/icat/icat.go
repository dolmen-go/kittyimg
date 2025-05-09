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

// icat - Print images in kitty/ghostty terminal emulators
//
// Usage
//
//	icat < file.png
//	icat file.png [file.png [...]]
//
// Install
//
//	go install github.com/dolmen-go/kittyimg/cmd/icat@latest
//
// Description
//
//	icat kitty.png
//
// is equivalent to:
//
//	kitten icat --transfer-mode=stream --align=left kitty.png
package main

import (
	"fmt"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/term"

	"github.com/dolmen-go/kittyimg"
)

func main() {
	var status int
	if err := _main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		status = 1
	}
	os.Exit(status)
}

func _main() error {
	if (len(os.Args) == 1 || os.Args[1] == "-") && !term.IsTerminal(int(os.Stdin.Fd())) {
		if err := kittyimg.Transcode(os.Stdout, os.Stdin); err != nil {
			return err
		}
		os.Stdout.WriteString("\n")
		return nil
	}

	for _, file := range os.Args[1:] {
		err := (func(file string) error {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()

			return kittyimg.Transcode(os.Stdout, f)
		})(file)
		if err != nil {
			return err
		}
		os.Stdout.WriteString("\n")
	}

	return nil
}
