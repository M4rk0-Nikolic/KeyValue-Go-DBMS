package record

type Record struct {
	Key string
	Value []byte
	Tombstone bool
}