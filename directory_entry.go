package cfb

import (
	"encoding/binary"
	"io"
	"strings"
)

type rawDirectoryEntry struct {
	DirectoryEntryName       [64]byte
	DirectoryEntryNameLength uint16
	ObjectType               ObjectType
	ColorFlag                uint8
	LeftSiblingID            uint32
	RightSiblingID           uint32
	ChildID                  uint32
	CLSID                    [16]byte
	StateBits                uint32
	CreationTime             uint64
	ModifiedTime             uint64
	StartingSectorLocation   uint32
	StreamSize               uint64
}

type DirectoryEntry struct {
	f *File

	id       int
	children []*DirectoryEntry
	path     []string

	raw rawDirectoryEntry
}

func readDirectoryEntry(r io.Reader) (*DirectoryEntry, error) {
	var d DirectoryEntry
	if err := binary.Read(r, binary.LittleEndian, &d.raw); err != nil {
		return nil, err
	}
	if err := d.validate(); err != nil {
		return nil, err
	}
	return &d, nil
}

func (d *DirectoryEntry) validate() error {
	if strings.ContainsAny(d.Name(), "/\\:!") {
		return ErrValidation
	}
	return nil
}

func (d *DirectoryEntry) Name() string {
	if d.raw.DirectoryEntryNameLength < 2 {
		return ""
	}
	// fmt.Println("name len =", len(d.raw.DirectoryEntryName))
	// fmt.Println("name DirectoryEntryNameLength =", d.raw.DirectoryEntryNameLength)
	return decodeUTF16String(d.raw.DirectoryEntryName[:d.raw.DirectoryEntryNameLength-2])
}

func (d *DirectoryEntry) Path() string {
	return strings.Join(d.path, "/")
}

func (d *DirectoryEntry) Type() ObjectType {
	return d.raw.ObjectType
}

func (d *DirectoryEntry) Size() uint64 {
	return d.raw.StreamSize
}

func (d *DirectoryEntry) StartingSector() uint32 {
	return d.raw.StartingSectorLocation
}

func (d *DirectoryEntry) object() (Object, error) {
	switch d.Type() {
	default:
		return nil, ErrInvalidObject
	case StorageObject:
		return d.storage()
	case StreamObject:
		return d.stream()
	}
}
