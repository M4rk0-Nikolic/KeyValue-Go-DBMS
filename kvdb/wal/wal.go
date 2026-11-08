package wal

import (
	"fmt"
	"log"
	"encoding/binary"
	"error"
	"hash/crc32"
)

type WALFragmentType string

const (
	Start WALFragmentType = "START"
	Middle WALFragmentType = "MIDDLE"
	End WALFragmentType = "END"
)

// <START> <COMMIT>?

type WALEntry struct {
	CRC uint32
	Key []byte
	Value []byte
	Tombstone bool
	Type string
	TC int64
}

type WALSegment struct {
	Path string
	Data []WALEntry
}

type WAL strcut {
	Segments []WALSegment
}


func (log *WALEntry) Append(entry []byte) error{

}

func (log *WALEntry) Read() ([]byte, error){

}