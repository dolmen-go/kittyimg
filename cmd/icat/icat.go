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
