package wal

import (
	"fmt"
	"log"
	"encoding/binary"
	"error"
	"os"
	"hash/crc32"
	"strconv"
)

/*
   +---------------+----------------+---------------+-----------+---------------+-----------------+-...-+--...--+---------+
   |    CRC (4B)   | Timestamp (8B) | Tombstone(1B) | Type (1B) | Key Size (8B) | Value Size (8B) | Key | Value | TC (8B) |
   +---------------+----------------+---------------+-----------+---------------+-----------------+-...-+--...--+---------+
   CRC = 32bit hash computed over the payload using CRC
   Key Size = Length of the Key data
   Tombstone = If this record was deleted and has a value
   Type = FULL, MIDDLE, START, END
   Value Size = Length of the Value data
   Key = Key data
   Value = Value data
   Timestamp = Timestamp of the operation in seconds
   TC = Trailing Canery
*/

const (
	CRC_SIZE = 4
	TIMESTAMP_SIZE = 8
	TOMBSTONE_SIZE = 1
	TYPE_SIZE = 1
	KEY_SIZE_SIZE = 8
	VALUE_SIZE_SIZE = 8
	TC_SIZE = 8
	
	CRC_START = 0
	TIMESTAMP_START = CRC_START + CRC_SIZE
	TOMBSTONE_START = TIMESTAMP_START + TIMESTAMP_SIZE
	TYPE_START = TOMBSTONE_START + TYPE_SIZE
	KEY_SIZE_START = TOMBSTONE_START + TOMBSTONE_SIZE
	VALUE_SIZE_START = KEY_SIZE_START + KEY_SIZE_SIZE
	KEY_START = VALUE_SIZE_START + VALUE_SIZE_SIZE
	TC_START = KEY_START + TC_SIZE
)

func CRC32(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

type WALFragmentType uint8

const (
	Full WALFragmentType = 0
	Start WALFragmentType = 1
	Middle WALFragmentType = 2
	End WALFragmentType = 3
)

// <START> <COMMIT>?

type WALEntry struct {
	CRC uint32
	Timestamp int64
	Tombstone bool
	Type WALFragmentType
	KeySize int64
	ValueSize int64
	Key []byte
	Value []byte
	TC int64
}

// 4 + 8 + 1 + 1 + 8 + 8 + 8 = 32 + 4 + 2 = 38 bytes of overhead required; if less - pad

// Maybe something like this though:
// type WALEntry struct {
//     Timestamp int64
//     Tombstone bool
//     Key       []byte
//     Value     []byte
// }

// type WALFragment struct {
//     CRC     uint32
//     Size    uint16
//     Type    WALFragmentType
//     Payload []byte
//     TC      uint32
// }

type WALSegment struct {
	SequenceNumber uint16
	Data []WALEntry
}

type SegWAL struct {
	Segments []WALSegment
	HotSegment WALSegment
	BlockManager *bm.BlockManager
}

func NewSegWAL(bm *bm.BlockManager) *SegWAL{
	err := os.MkdirAll("../data/wal", 0755)
	if err != nil {
		return err
	}

	entries, err2 := os.ReadDir("../data/wal")
	if err2 != nil{
		return err2
	}

	var wal_logs []uint16 = make([]uint16, 0)

	for _, entry := range entries{
		if !entry.IsDir(){
			if (strings.HasPrefix(entry.Name(), "wal_")) &&  (strings.HasSufix(entry.Name(), ".log")){

			}
		}
	}

	bm.WriteBlock("../data/wal")

	SegWAL{
		BlockManager: bm
	}



}

func (log *SegWAL) AppendRecord() error{
	if HotSegment == nil | len(SegWAL.Segments) == 0 {
		newseg := WALSegment {
			SequenceNumber: 1
		}
		SegWAL.Segments = append(SegWAL.Segments, newseg)
		SegWal.HotSegment = newseg
	}

	str := strconv.FormatUint(uint64(num), 10)

	if len(str) < 5 {
		str = strings.Repeat("0", 5-len(str)) + str
	}

	bm.WriteBlock(("../data/wal/wal." + str + ".log"))
	// bruh, how can i append if my only operations are Read and Write block according to the specs


}

func (log *SegWAL) Read() ([]byte, error){

}