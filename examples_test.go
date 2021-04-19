package kittyimg_test

import (
	"embed"
	"image"
	"os"

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
