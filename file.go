package cfb

import (
	"io"
	"io/ioutil"
	"strings"
)

type File struct {
	header           *FileHeader
	r                io.ReaderAt
	fat              []uint32
	minifat          []uint32
	directoryEntries []*DirectoryEntry
}

func New(r io.ReaderAt) (*File, error) {
	f := &File{r: r}
	h, err := readFileHeader(&offsetReader{f.r, 0})
	if err != nil {
		return nil, err
	}
	f.header = h
	if err := f.buildFAT(); err != nil {
		return nil, err
	}
	if err := f.buildMiniFAT(); err != nil {
		return nil, err
	}
	if err := f.buildDirectoryTree(); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *File) Header() *FileHeader {
	return f.header
}

func (f *File) Objects() ([]Object, error) {
	var os []Object
	for _, d := range f.directoryEntries[1:] {
		o, err := d.object()
		if err != nil {
			return nil, err
		}
		os = append(os, o)
	}
	return os, nil
}

func (f *File) Get(path string) (Object, error) {
	ps := strings.Split(path, "/")

	d := f.directoryEntries[0]
	for _, p := range ps {
		var t *DirectoryEntry
		found := false
		for _, c := range d.children {
			if c.Name() == p {
				t = c
				found = true
				break
			}
		}
		if !found {
			return nil, ErrObjectNotFound
		}
		d = t
	}

	return d.object()
}

func (f *File) buildFAT() error {
	sectors := make([]uint32, f.header.raw.NumberOfFATSectors)
	copy(sectors, f.header.raw.DIFAT[:f.header.raw.NumberOfFATSectors])

	if f.header.raw.NumberOfDIFATSectors > 0 {
		s := f.header.raw.FirstDIFATSectorLocation
		for i := 0; i < int(f.header.raw.NumberOfDIFATSectors); i++ {
			b, err := f.readSector(s)
			if err != nil {
				return err
			}
			difat, err := bytesToUint32s(b)
			if err != nil {
				return err
			}
			sectors = append(sectors, difat[:len(difat)-1]...)
			s = difat[len(difat)-1]
		}
	}

	var b []byte
	for _, s := range sectors {
		chunk, err := f.readSector(s)
		if err != nil {
			return err
		}
		b = append(b, chunk...)
	}

	fat, err := bytesToUint32s(b)
	if err != nil {
		return err
	}
	f.fat = fat
	return nil
}

func (f *File) buildMiniFAT() error {
	sr, err := newSectorReader(f.r, f.header.SectorSize(), f.header.raw.FirstMiniFATSectorLocation, f.fat, func(sector uint32) int64 {
		return int64((sector + 1) * f.header.SectorSize())
	})
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(sr)
	if err != nil {
		return err
	}
	minifat, err := bytesToUint32s(b)
	if err != nil {
		return err
	}
	f.minifat = minifat
	return nil
}

func (f *File) readDirectoryEntries() ([]*DirectoryEntry, error) {
	var ds []*DirectoryEntry
	sr, err := newSectorReader(f.r, f.header.SectorSize(), f.header.raw.FirstDirectorySectorLocation, f.fat, func(sector uint32) int64 {
		return int64((sector + 1) * f.header.SectorSize())
	})
	if err != nil {
		return nil, err
	}
	for i := 0; ; i++ {
		d, err := readDirectoryEntry(sr)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if d.Type() == UnknownObject {
			continue
		}
		d.f = f
		d.id = i
		ds = append(ds, d)
	}
	return ds, nil
}

func (f *File) buildDirectoryTree() error {
	ds, err := f.readDirectoryEntries()
	if err != nil {
		return err
	}

	var walk func(uint32, *DirectoryEntry, []string)
	walk = func(id uint32, parent *DirectoryEntry, prefixes []string) {
		if id == noStream {
			return
		}
		d := ds[id]
		parent.children = append(parent.children, d)

		walk(d.raw.LeftSiblingID, parent, prefixes)
		walk(d.raw.RightSiblingID, parent, prefixes)

		ps := make([]string, len(prefixes))
		copy(ps, prefixes)
		ps = append(ps, d.Name())
		d.path = ps
		walk(d.raw.ChildID, d, ps)
	}
	walk(ds[0].raw.ChildID, ds[0], []string{})

	f.directoryEntries = ds
	return nil
}

func (f *File) readSector(sector uint32) ([]byte, error) {
	b := make([]byte, f.header.SectorSize())
	if n, err := f.r.ReadAt(b, int64((sector+1)*f.header.SectorSize())); err != nil {
		return nil, err
	} else if n != int(f.header.SectorSize()) {
		return nil, ErrInsufficientData
	}
	return b, nil
}
