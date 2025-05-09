package kittyimg_test

import (
	"embed"
	"image"
	"os"
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

func TestExample(*testing.T) {
	Example()
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

func TestExampleTranscode(*testing.T) {
	ExampleTranscode()
}
