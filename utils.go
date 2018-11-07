package cfb

import (
	"bytes"
	"encoding/binary"
	"io"
	"unicode/utf16"
)

type offsetReader struct {
	r      io.ReaderAt
	offset int64
}

func (o *offsetReader) Read(b []byte) (int, error) {
	return o.r.ReadAt(b, o.offset)
}

func bytesToUint32s(b []byte) ([]uint32, error) {
	var ret []uint32
	r := bytes.NewReader(b)
	for {
		var u uint32
		if err := binary.Read(r, binary.LittleEndian, &u); err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		ret = append(ret, u)
	}
	return ret, nil
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func decodeUTF16String(b []byte) string {
	u := make([]uint16, len(b)/2)
	for i := range u {
		u[i] = (uint16(b[i*2+1]) << 8) | uint16(b[i*2])
	}
	return string(utf16.Decode(u))
}

func isEqualBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
