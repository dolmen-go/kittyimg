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
	"bytes"
	"fmt"
	"image"
	"io"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/term"

	"github.com/dolmen-go/kittyimg"
	"github.com/dolmen-go/kittyimg/internal/writers"
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
		if err := transcode(os.Stdin, os.Stdout); err != nil {
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

			return transcode(f, os.Stdout)
		})(file)
		if err != nil {
			return err
		}
		os.Stdout.WriteString("\n")
	}

	return nil
}

func transcode(r io.Reader, w io.Writer) error {
	var buf bytes.Buffer
	in := io.TeeReader(r, &buf)
	cfg, format, err := image.DecodeConfig(in)
	if err != nil {
		return readError(r, err)
	}
	// Restart from byte 0
	in = io.MultiReader(&buf, r)

	// For PNG we send the raw file that probably has better compression
	// https://sw.kovidgoyal.net/kitty/graphics-protocol/#png-data
	if format == "png" {
		if _, err = fmt.Fprintf(w, "\033_Gq=1,a=T,f=100,s=%d,v=%d,", cfg.Width, cfg.Height); err != nil {
			return err
		}

		var pw writers.PayloadWriter
		pw.Reset(w)

		if _, err = io.Copy(&pw, in); err != nil {
			return err
		}
		return pw.Close()
	}

	img, _, err := image.Decode(in)
	if err != nil {
		return readError(r, err)
	}
	// return icat(w, img)
	return kittyimg.Fprint(w, img)
}

func readError(r io.Reader, err error) error {
	if r, ok := r.(interface{ Name() string }); ok {
		if name := r.Name(); name != "" {
			return fmt.Errorf("%s: %w", r.Name(), err)
		}
	}
	return err
}

func icat(w io.Writer, img image.Image) error {
	bounds := img.Bounds()

	// f=32 => RGBA
	_, err := fmt.Fprintf(w, "\033_Gq=1,a=T,f=32,s=%d,v=%d,t=d,", bounds.Dx(), bounds.Dy())
	if err != nil {
		return err
	}

	var zw writers.PayloadWriter
	zw.Reset(w)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 8 reduces this to the range [0, 255].
			if _, err = zw.Write([]byte{byte(r >> 8), byte(g >> 8), byte(b >> 8), byte(a >> 8)}); err != nil {
				return err
			}
		}
	}
	return zw.Close()
}
