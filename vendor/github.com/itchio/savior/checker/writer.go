package checker

import (
	"fmt"
	"io"
)

type Writer struct {
	reference []byte
	offset    int64
}

var _ io.WriteSeeker = (*Writer)(nil)

func NewWriter(reference []byte) *Writer {
	return &Writer{
		reference: reference,
	}
}

func (cw *Writer) Write(buf []byte) (int, error) {
	n := 0
	for i := 0; i < len(buf); i++ {
		if cw.offset >= int64(len(cw.reference)) {
			return n, fmt.Errorf("out of bounds write: %d but max length is %d", cw.offset, len(cw.reference))
		}

		expected := cw.reference[cw.offset]
		actual := buf[i]
		if expected != actual {
			return n, fmt.Errorf("at byte %d, expected %x but got %x", cw.offset, expected, actual)
		}
		cw.offset++
		n++
	}
	return n, nil
}

func (cw *Writer) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekStart {
		return cw.offset, fmt.Errorf("unsupported whence value %d", whence)
	}

	if offset > int64(len(cw.reference)) {
		return cw.offset, fmt.Errorf("out of bounds seek: %d but max length is %d", cw.offset, len(cw.reference))
	}
	if offset < 0 {
		return cw.offset, fmt.Errorf("out of bounds seek: %d which is < 0", cw.offset)
	}

	cw.offset = offset
	return cw.offset, nil
}
