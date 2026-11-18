package sstable

import (
	bm "github.com/M4rk0-Nikolic/KeyValue-Go-DBMS"
)

type SSTable struct{
	Path string
	BlockSize uint64
	BlockManager *bm.BlockManager
	CurrentOffset uint64
}

type FragmentType uint8

const (
	Full FragmentType = 0
	Start FragmentType = 1
	Middle FragmentType = 2
	End FragmentType = 3
)

type Fragment struct {
	CRC uint32
	Size uint16
	Type FragmentType
	Payload []byte
}

type ValueEntry struct {
	Timestamp uint64
	Tombstone uint64
	Value []byte
}

type IndexEntry struct {
	Offset uint64
	Key []byte
}

type SummaryEntry struct {
    Key         []byte
    IndexOffset uint64
}

type Summary struct {
    MinKey  []byte
    MaxKey  []byte
    Entries []SummaryEntry
}

type Footer struct {
    DataOffset       uint64
    FilterOffset     uint64
    SummaryOffset    uint64
    IndexOffset      uint64
    HashIndexOffset  uint64
    MetadataOffset   uint64
}






