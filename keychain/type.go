package keychain

type keychainHeader struct {
	Signature        [4]byte
	MajorVersion     int16
	MinorVersion     int16
	Unknown1         int32
	TableArrayOffset int32
	Unknown2         int32
}

type keychainTableList struct {
	DataSize       int32
	NumberOfTables int32
	Table          []keychainTable
}

type keychainTableHeader struct {
	DataSize             int32
	RecordType           RecordTypes
	NumberOfRecord       int32
	RecordListOffset     int32
	Unknown1             int32
	Unknown2             int32
	NumberOfRecordOffset int32
	RecordOffset         []int32
}

type keychainTable struct {
	offset     int32
	Header     keychainTableHeader
	RecordList []Record
}

type RecordHeader struct {
	DataSize    int32
	RecordIndex int32
	Unknown1    int32
	Unknown2    int32
	KeyDataSize int32
	Unknown3    int32
}

type Record struct {
	Header          RecordHeader
	AttributeOffset []int32
}

const (
	signature = "kych"
)

type Keychain struct {
	Header    keychainHeader
	TableList keychainTableList
}
