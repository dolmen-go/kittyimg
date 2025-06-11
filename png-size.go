//go:build ignore

// Command png-size adjusts the size of PNG file to reach a target size.
//
// png-size injects a "tEXt" chunk (non-compressed text) of the appropriate length
// to grow the file.
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"hash/crc32"
	"os"
)

const (
	pngSignature = "\x89PNG\r\n\x1a\n"
	pngTEXTlabel = "tEXt"
)

var pngIEND = []byte{0, 0, 0, 0, 'I', 'E', 'N', 'D', 0xae, 0x42, 0x60, 0x82}

func main() {
	var targetPNGSize int
	flag.IntVar(&targetPNGSize, "size", 0, "target size")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Println("Usage: go run png-size.go -size=<target-size> <input.png> <output.png>")
		os.Exit(1)
	}

	inputPath := flag.Arg(0)
	outputPath := flag.Arg(1)
	maxInputSize := targetPNGSize - 13 - len("comment") // Maximum allowed size for the input PNG

	pngData, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Printf("Error reading input PNG file: %v\n", err)
		os.Exit(1)
	}

	currentSize := len(pngData)
	fmt.Printf("Input PNG size: %d bytes\n", currentSize)

	if currentSize > maxInputSize {
		fmt.Printf("Error: Input PNG size (%d bytes) must be <= %d bytes to allow adding tEXt chunk.\n", currentSize, maxInputSize)
		os.Exit(1)
	}

	// 2. Vérifier la signature PNG
	if !bytes.HasPrefix(pngData, []byte(pngSignature)) {
		fmt.Println("Error: Not a valid PNG file (invalid signature).")
		os.Exit(1)
	}

	// Verify IEND chunk
	iendOffset := len(pngData) - 12
	if !bytes.Equal(pngData[iendOffset:], pngIEND) {
		fmt.Println("Error: Not a valid PNG file (invalid IEND chunk).")
		os.Exit(1)
	}

	// Key: "comment" (7 byte)
	// Value: X bytes
	// Total tEXt chunk size = 4 (len) + 4 (type) + 7 (key) + 1 (null) + X (value) + 4 (CRC) = 20 + X

	// Bytes already in the file + (14 + X) = targetPNGSize
	// X = targetPNGSize - currentSize - 14

	bytesNeededForTextValue := targetPNGSize - currentSize - 20 // 14 est la surcharge pour le chunk tEXt "comment"

	if bytesNeededForTextValue < 0 {
		fmt.Printf("Error: Input PNG size (%d bytes) is already too large. Need %d more bytes for tEXt, but result is negative.\n", currentSize, bytesNeededForTextValue)
		os.Exit(1)
	}

	fmt.Printf("Bytes needed for tEXt value: %d\n", bytesNeededForTextValue)

	// gives 0123 once encoded as base64
	const patternBase = "\xd3\x5d\xb7" + "\xd3\x5d\xb7"
	pattern := patternBase[iendOffset%4+2:][:3] // FIXME test formula with different input files

	// 5. Construire le chunk tEXt "comment"
	tEXt := append(
		[]byte("tEXt"+"comment"+"\x00"), // Null separator
		bytes.Repeat([]byte(pattern), (bytesNeededForTextValue+len(pattern)-1)/len(pattern))[:bytesNeededForTextValue]...,
	)

	tEXtLength := uint32(len(tEXt) - 4)

	// Calculer le CRC pour le chunk tEXt (type + données)
	crc := crc32.New(crc32.IEEETable)
	crc.Write(tEXt)
	tEXt = binary.BigEndian.AppendUint32(tEXt, crc.Sum32())

	pngData = binary.BigEndian.AppendUint32(
		pngData[:iendOffset],
		tEXtLength,
	)
	pngData = append(pngData, tEXt...)
	pngData = append(pngData, pngIEND...)

	finalSize := len(pngData)

	fmt.Printf("Final PNG size: %d bytes\n", finalSize)

	if finalSize != targetPNGSize {
		fmt.Printf("Error: Final PNG size is %d, expected %d. This might be due to unexpected PNG structure or rounding issues.\n", finalSize, targetPNGSize)
		os.Exit(1)
	}

	err = os.WriteFile(outputPath, pngData, 0644)
	if err != nil {
		fmt.Printf("Error writing output PNG file: %v\n", err)
		os.Exit(1)
	}
}
