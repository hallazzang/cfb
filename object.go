package cfb

type Object interface {
	Name() string
	Path() string
	Type() ObjectType
	Size() uint64
	ReadAt([]byte, int64) (int, error)
	Seek(int64, int) (int64, error)
	Read([]byte) (int, error)
}
