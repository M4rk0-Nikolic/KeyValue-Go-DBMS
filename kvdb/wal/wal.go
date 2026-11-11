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

type WAL interface {
    Append(entry WALEntry) error
    AppendBatch(entries []WALEntry) error
    Recover() ([]WALEntry, error)
}

const (
	CRC_SIZE = 4
	SIZE_SIZE = 2
	TYPE_SIZE = 1
	TOTAL = CRC_SIZE + SIZE_SIZE + TYPE_SIZE

	TIMESTAMP_SIZE = 8
	TOMBSTONE_SIZE = 1
	KEY_SIZE_SIZE = 8
	VALUE_SIZE_SIZE = 8
	TC_SIZE = 8
)

type WALFragmentType uint8

const (
	Full WALFragmentType = 0
	Start WALFragmentType = 1
	Middle WALFragmentType = 2
	End WALFragmentType = 3
)

type WALFragment struct {
	CRC uint32
	Size uint16
	Type WALFragmentType
	Payload []byte
}

type WALOperation struct {
	Tombstone bool
	KeySize uint64
	ValueSize uint64
	Key []byte
	Value []byte
}

type WALEntry struct {
	Timestamp uint64
	NumberOfOperations uint16
	Operations []WALOperation
	TC uint64
}

type WALSegment struct {
	SequenceNumber uint16
	Path string
}

type SegWAL struct {
	BlockSize uint64
	BlocksInSegment uint64
	SegmentBytes uint64
	Segments []WALSegment
	HotSegment WALSegment
	BlockManager *bm.BlockManager
}

func (op WALOperation)Serialize() []byte{
	length := len(op.Key) + len(op.Value) + TOMBSTONE_SIZE + KEY_SIZE_SIZE + VALUE_SIZE_SIZE
	serialized := make([]byte, length)
	offset := 0
	if op.Tombstone {
        serialized[offset] = 1
    }
    offset += TOMBSTONE_SIZE
	binary.LittleEndian.PutUint64(serialized[offset:offset+8], op.KeySize)
	offset += KEY_SIZE_SIZE
	binary.LittleEndian.PutUint64(serialized[offset:offset+8], op.ValueSize)
	offset += VALUE_SIZE_SIZE
	offset += copy(serialized[offset:], op.Key)
	offset += copy(serialized[offset:], op.Value)
	return serialized
}


func (wale WALEntry)Serialize() []byte{
	serialized := make([]byte, TIMESTAMP_SIZE + NUMBER_OF_OPERATIONS_SIZE)
	offset := 0
	binary.LittleEndian.PutUint16(wale.Timestamp, serialized[offset:offset+8])
	offset += TIMESTAMP_SIZE
	binary.LittleEndian.PutUint16(wale.NumberOfOperations, serialized[offset:offset+2])
	for _, op := range wale.Operations{
		serialized = append(serialized, op.Serialize()...)
	}
	tc := make([]byte, TC_SIZE)
    binary.LittleEndian.PutUint64(tc, entry.TC)
    serialized = append(serialized, tc...)

	return serialized
}

func GetSegPath(num uint16) {
	str := strconv.FormatUint(uint64(num), 10)

	if len(str) < 5 {
		str = "../data/wal/wal" + strings.Repeat("0", 5-len(str)) + str + ".log"
	}

	return str
}

func (segwal SegWAL)Recover(memtable *memtable.Memtable){
	err := os.MkdirAll("../data/wal", 0755)
	if err != nil {
		return err
	}

	entries, err2 := os.ReadDir("../data/wal")
	if err2 != nil{
		return err2
	}

	var wal_logs []WALSegment = []WALSegment{}

	for _, entry := range entries{
		if !entry.IsDir(){
			if (strings.HasPrefix(entry.Name(), "wal_")) &&  (strings.HasSufix(entry.Name(), ".log")){
				wal_logs = append(wal_logs, WALSegment{ SequenceNumber: uint16(_) Path: entry.Name() })
			}
		}
	}

	// Determine which WAL Segments have been flushed to disk via MANIFEST file
	// Remaining WAL Segments need to be replayed

	toBeReplayed := entries

	for

}

func (segwal *SegWAL) SwapSegment() {
	old_num := segwal.HotSegment.SequenceNumber
	new_seg := WALSegment{ SequenceNumber: uint16(old_num + 1) Path: GetSegPath(old_num + 1) }
	segwal.Segments = append(segwal.Segments, segwal.HotSegment)
	segwal.HotSegment = new_se
	segwal.SegmentBytes = 0
}

func SerializeFragment(fragment WALFragment) {
	result := make([]byte, fragment.Size+2+4)
	
	binary.LittleEndian.PutUint32(result, fragment.Hash)
	binary.LittleEndian.PutUint16(result[4:], fragment.Size)

	copy(result[6:], fragment.Payload)

	return result
}

func (segwal *SegWAL) Append(record record.Record) error{
	if segwal.HotSegment == nil | len(segwal.Segments) == 0 {
		newseg := WALSegment {
			SequenceNumber: 1,
			Path: GetSegPath(1)
		}
		segwal.Segments = append(segwal.Segments, newseg)
		segwal.HotSegment = newseg
	}

	segmentSize := segwal.BlockSize * wal.BlocksInSegment
	blockIndex := tipOffset / wal.BlockSize

	first := true
	walentry := WALEntry {
		Timestamp: time.Now().Unix(),
		NumberOfOperations: 1,
		Operations: []WALOperation{WALOperation{
			Tombstone: record.Tombstone
			Key: record.Key
			Value: record.Value
			KeySize: len(record.Key)
			ValueSize: len(record.Value)
		}},
		TC: 0xDEADBEEFFEEDFACE
		}
	data := walentry.Serialize()

	for len(data) > 0 {
		block := BlockManager.ReadBlock(segwal.HotSegment.Path, blockIndex)
		
		if segwal.SegmentBytes == segmentSize {
			if err := segwal.SwapSegment(); err != nil {
				return err
			}
			blockIndex = 0
		}

		remaining := segmentSize - (wal.SegmentBytes % segwal.BlockSize)

		if remaining < TOTAL {
			copy(block[segwal.BlockSize - remaining: segwal.BlockSize], make([]byte, remaining))
			segwal.BlockManager.Write(segwal.HotSegment.Path, blockIndex, block)
		}

		n := min(len(data), remaining-7)
		last := n == len(data)

		var typ WALFragmentType
		switch {
		case first && last:
			typ = Full
		case first:
			typ = Start
		case last:
			typ = End
		case !first && !last:
			typ = Middle
		}

		payload := data[:n]

		copy(block[segwal.BlockSize-remaining:], SerializeFragment(WALFragment{ 
			Type: typ,
			Payload: payload,
			Hash := crc32.ChecksumIEEE(payload),
			Size := len(payload)
		}))

		BlockManager.WriteBlock(segwal.HotSegment.Path, blockIndex, block)

		if type != Full{
			blockIndex += 1
		}

		segwal.SegmentBytes += uint64(7 + n)

		data = data[n:]
		first = false
	}

}