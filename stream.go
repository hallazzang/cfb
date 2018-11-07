package cfb

import "fmt"

type Stream struct {
	d *DirectoryEntry
	r *SectorReader
}

func (d *DirectoryEntry) stream() (*Stream, error) {
	if d.Type() != StreamObject {
		return nil, ErrWrongObjectType
	}
	var sr *SectorReader
	if d.Size() < uint64(d.f.header.raw.MiniStreamCutoffSize) {
		tsr, err := newSectorReader(d.f.r, d.f.header.SectorSize(), d.f.directoryEntries[0].StartingSector(), d.f.fat, func(sector uint32) int64 {
			return int64((sector + 1) * d.f.header.SectorSize())
		})
		if err != nil {
			return nil, err
		}
		sr, err = newSectorReader(tsr, d.f.header.MiniSectorSize(), d.StartingSector(), d.f.minifat, func(sector uint32) int64 {
			return int64(sector * d.f.header.MiniSectorSize())
		})
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		sr, err = newSectorReader(d.f.r, d.f.header.SectorSize(), d.StartingSector(), d.f.fat, func(sector uint32) int64 {
			return int64((sector + 1) * d.f.header.SectorSize())
		})
		if err != nil {
			return nil, err
		}
	}
	sr.maxOffset = int64(d.Size())
	return &Stream{d, sr}, nil
}

func (s *Stream) String() string {
	return fmt.Sprintf("Stream{Path:%q Size:%d}", s.Path(), s.Size())
}

func (s *Stream) Name() string {
	return s.d.Name()
}

func (s *Stream) Path() string {
	return s.d.Path()
}

func (s *Stream) Type() ObjectType {
	return StreamObject
}

func (s *Stream) Size() uint64 {
	return s.d.Size()
}

func (s *Stream) ReadAt(b []byte, offset int64) (int, error) {
	return s.r.ReadAt(b, offset)
}

func (s *Stream) Seek(offset int64, whence int) (int64, error) {
	return s.r.Seek(offset, whence)
}

func (s *Stream) Read(b []byte) (int, error) {
	return s.r.Read(b)
}
