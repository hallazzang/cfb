package cfb

import (
	"encoding/binary"
	"io"
)

type rawFileHeader struct {
	HeaderSignature              [8]byte
	HeaderCLSID                  [16]byte
	MinorVersion                 uint16
	MajorVersion                 uint16
	ByteOrder                    uint16
	SectorShift                  uint16
	MiniSectorShift              uint16
	Reserved                     [6]byte
	NumberOfDirectorySectors     uint32
	NumberOfFATSectors           uint32
	FirstDirectorySectorLocation uint32
	TransactionSignatureNumber   uint32
	MiniStreamCutoffSize         uint32
	FirstMiniFATSectorLocation   uint32
	NumberOfMiniFATSectors       uint32
	FirstDIFATSectorLocation     uint32
	NumberOfDIFATSectors         uint32
	DIFAT                        [109]uint32
}

type FileHeader struct {
	raw rawFileHeader
}

func readFileHeader(r io.Reader) (*FileHeader, error) {
	var h FileHeader
	if err := binary.Read(r, binary.LittleEndian, &h.raw); err != nil {
		return nil, err
	}
	if err := h.validate(); err != nil {
		return nil, err
	}
	return &h, nil
}

func (h *FileHeader) validate() error {
	// TODO: return reason of validation error
	if !isEqualBytes(h.raw.HeaderSignature[:], []byte("\xd0\xcf\x11\xe0\xa1\xb1\x1a\xe1")) {
		return ErrValidation
	} else if !isEqualBytes(h.raw.HeaderCLSID[:], []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")) {
		return ErrValidation
	} else if h.raw.MinorVersion != 0x3e {
		return ErrValidation
	} else if h.raw.MajorVersion != 0x3 && h.raw.MajorVersion != 0x4 {
		return ErrValidation
	} else if h.raw.ByteOrder != 0xfffe {
		return ErrValidation
	} else if h.raw.SectorShift != 0x9 && h.raw.SectorShift != 0xc {
		return ErrValidation
	} else if h.raw.MiniSectorShift != 0x6 {
		return ErrValidation
	} else if !isEqualBytes(h.raw.Reserved[:], []byte("\x00\x00\x00\x00\x00\x00")) {
		return ErrValidation
	} else if h.raw.MiniStreamCutoffSize != 0x1000 {
		return ErrValidation
	}
	return nil
}

func (h *FileHeader) SectorSize() uint32 {
	return 1 << h.raw.SectorShift
}

func (h *FileHeader) MiniSectorSize() uint32 {
	return 1 << h.raw.MiniSectorShift
}
