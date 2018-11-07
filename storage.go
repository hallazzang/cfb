package cfb

import (
	"fmt"
	"io"
)

type Storage struct {
	d *DirectoryEntry
}

func (d *DirectoryEntry) storage() (*Storage, error) {
	if d.Type() != StorageObject {
		return nil, ErrWrongObjectType
	}
	return &Storage{d}, nil
}

func (s *Storage) String() string {
	return fmt.Sprintf("Storage{Path:%q}", s.Path())
}

func (s *Storage) Name() string {
	return s.d.Name()
}

func (s *Storage) Path() string {
	return s.d.Path()
}

func (s *Storage) Type() ObjectType {
	return StorageObject
}

func (s *Storage) Size() uint64 {
	return 0
}

func (s *Storage) ReadAt(b []byte, offset int64) (int, error) {
	return 0, io.EOF
}

func (s *Storage) Seek(offset int64, whence int) (int64, error) {
	// TODO: error check
	return 0, nil
}

func (s *Storage) Read(b []byte) (int, error) {
	return 0, io.EOF
}
