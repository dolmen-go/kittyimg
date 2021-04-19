package kittyimg_test

import (
	"image"
	"os"

	// Plugin to decode GIF
	_ "image/gif"

	"github.com/dolmen-go/kittyimg"
)

func Example() {
	f, err := os.Open("dolmen.gif")
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
