package kittyimg

import (
	"compress/zlib"
	"encoding/base64"
	"io"
)

const (
	// https://sw.kovidgoyal.net/kitty/graphics-protocol.html#remote-client
	chunkEncSize = 4096
	chunkRawSize = (chunkEncSize / 4) * 3
)

// streamPayload is an io.WriteCloser that encodes the payload data in the stream.
// https://sw.kovidgoyal.net/kitty/graphics-protocol.html#remote-client
type streamPayload struct {
	bufEnc [chunkEncSize]byte
	bufRaw [chunkRawSize]byte
	n      int
	w      io.Writer
}

func (spw *streamPayload) Reset(w io.Writer) {
	spw.w = w
	spw.n = 0
}

func (spw *streamPayload) encode() error {
	// fmt.Fprintln(os.Stderr, len(bufRaw), "=>", (len(bufRaw)+2)/3*4)

	base64.StdEncoding.Encode(spw.bufEnc[:], spw.bufRaw[:spw.n])
	_, err := spw.w.Write(spw.bufEnc[:(spw.n+2)/3*4])
	spw.n = 0
	return err
}

func (spw *streamPayload) Write(b []byte) (n int, err error) {
	for len(b) > 0 {
		if spw.n == cap(spw.bufRaw) {
			_, err = spw.w.Write([]byte("m=1;"))
			if err != nil {
				return
			}
			err = spw.encode()
			if err != nil {
				return
			}
			_, err = spw.w.Write([]byte("\033\\\033_G"))
			if err != nil {
				return
			}
		}

		l := copy(spw.bufRaw[spw.n:], b)
		spw.n += l
		n += l
		b = b[l:]
	}
	return
}

// Close closes the Writer, flushing any unwritten data to the underlying io.Writer, but does not close the underlying io.Writer.
func (spw *streamPayload) Close() (err error) {
	if spw.n == 0 {
		_, err = spw.w.Write([]byte("m=0;\033\\"))
		return
	}
	_, err = spw.w.Write([]byte("m=0;"))
	if err != nil {
		return
	}
	err = spw.encode()
	if err != nil {
		return
	}
	_, err = spw.w.Write([]byte("\033\\"))
	return
}

// zlibPayload is an io.WriteCloser.
// https://sw.kovidgoyal.net/kitty/graphics-protocol.html#compression
type zlibPayload struct {
	buffer [16384]byte
	n      int
	spw    streamPayload
	zw     *zlib.Writer
}

func (zp *zlibPayload) Reset(w io.Writer) {
	_, _ = w.Write([]byte("o=z,"))
	zp.spw.Reset(w)
	zp.zw = zlib.NewWriter(&zp.spw)
	zp.n = 0
}

func (zp *zlibPayload) Write(b []byte) (n int, err error) {
	for len(b) > 0 {
		if zp.n == cap(zp.buffer) {
			_, err = zp.zw.Write(zp.buffer[:])
			if err != nil {
				return
			}
			zp.n = 0
		}
		m := copy(zp.buffer[zp.n:], b)
		zp.n += m
		n += m
		b = b[m:]
	}
	return
}

// Close closes the Writer, flushing any unwritten data to the underlying io.Writer, but does not close the underlying io.Writer.
func (zp *zlibPayload) Close() error {
	if zp.n > 0 {
		if _, err := zp.zw.Write(zp.buffer[:zp.n]); err != nil {
			return err
		}
	}
	if err := zp.zw.Close(); err != nil {
		return err
	}
	return zp.spw.Close()
}
