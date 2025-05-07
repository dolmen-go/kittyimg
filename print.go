// Package kittyimg provides utilities to show image in a graphic terminal emulator supporting kitty's "terminal graphics protocol".
//
// See https://sw.kovidgoyal.net/kitty/graphics-protocol.html.
package kittyimg

import (
	"fmt"
	"image"
	"io"

	"github.com/dolmen-go/kittyimg/internal/writers"
)

func Fprint(w io.Writer, img image.Image) error {
	bounds := img.Bounds()

	// f=32 => RGBA
	_, err := fmt.Fprintf(w, "\033_Gq=1,a=T,f=32,s=%d,v=%d,t=d,", bounds.Dx(), bounds.Dy())
	if err != nil {
		return err
	}

	buf := make([]byte, 0, min(bounds.Dx()*bounds.Dy()*4, 16384)) // Multiple of 4 (RGBA)

	// var p payloadWriter
	var p writers.ZlibPayloadWriter
	p.Reset(w)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if len(buf) == cap(buf) {
				if _, err = p.Write(buf); err != nil {
					return err
				}
				buf = buf[:0]
			}
			r, g, b, a := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 8 reduces this to the range [0, 255].
			buf = append(buf, byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
		}
	}

	if _, err = p.Write(buf); err != nil {
		return err
	}
	return p.Close()
}

func Fprintln(w io.Writer, img image.Image) error {
	err := Fprint(w, img)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{'\n'})
	return err
}
