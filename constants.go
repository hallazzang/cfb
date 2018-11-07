package cfb

//go:generate stringer -type ObjectType
type ObjectType uint8

const (
	maxRegSect uint32 = 0xfffffffa
	difatSect  uint32 = 0xfffffffc
	fatSect    uint32 = 0xfffffffd
	endOfChain uint32 = 0xfffffffe
	freeSect   uint32 = 0xffffffff

	maxRegSID uint32 = 0xfffffffa
	noStream  uint32 = 0xffffffff

	UnknownObject     ObjectType = 0x00
	StorageObject     ObjectType = 0x01
	StreamObject      ObjectType = 0x02
	RootStorageObject ObjectType = 0x05
)
