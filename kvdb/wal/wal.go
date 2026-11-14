package wal

import (
	"fmt"
	"log"
	"encoding/binary"
	"error"
	"os"
	"hash/crc32"
	"strconv"
	bm "github.com/M4rk0-Nikolic/KeyValue-Go-DBMS"
)

type WAL interface {
    Append(entry Entry) error
    AppendBatch(entries []Entry) error
    Recover() ([]Entry, error)
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

type Operation struct {
	Tombstone bool
	KeySize uint64
	ValueSize uint64
	Key []byte
	Value []byte
}

type Entry struct {
	Timestamp uint64
	NumberOfOperations uint16
	Operations []Operation
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
	maxBlock := wal.BlocksInSegment

	fragmentStarted := false

	buffer := []byte{}

	for index, entry := range toBeReplayed {
		walEntry := []byte{}
		for currentBlock := 0 ; currentBlock < maxBlock ; currentBlock++ {
			data := segwal.BlockManager.ReadBlock(entry.Name(), currentBlock)
			for ;;{
				fragmentHeader := data[:7]
				crc32 := binary.LittleEndian.Uint32(fragmentHeader)
				size := binary.LittleEndian.Uint16(fragmentHeader[4:])
				typ := FragmentType(fragmentHeader[6])

				switch typ{
				case Full:
					if crc32 != crc32.ChecksumIEEE(data[7:size]){
						log.Println("Failed to recover one record.")
						buffer = []byte{}
					}
					buffer = append(buffer, data[7:size])
					break
				case End:
					if fragmentStarted == false {
						data = data[7+size:]
						continue
					}
				case Start:
					fragmentStarted = true
					
				case End:
					fragmentStarted = false

				case Middle:
					
				}

				data = data[7+size:]
			}
			DeserializeWALEntry()
		}
	}

}

func (segwal *SegWAL) SwapSegment() {
	oldNum := segwal.HotSegment.SequenceNumber
	newSeg := WALSegment{ SequenceNumber: uint16(oldNum + 1) Path: GetSegPath(oldNum + 1) }
	segwal.Segments = append(segwal.Segments, segwal.HotSegment)
	segwal.HotSegment = newSeg
	segwal.SegmentBytes = 0
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

	segmentSize := segwal.BlockSize * segwal.BlocksInSegment
	blockIndex := segwal.SegmentBytes / segwal.BlockSize

	first := true
	walentry := Entry {
		Timestamp: time.Now().Unix(),
		NumberOfOperations: 1,
		Operations: []Operation{Operation{
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

		remaining := segmentSize - (segwal.SegmentBytes % segwal.BlockSize)

		if remaining < TOTAL {
			copy(block[segwal.BlockSize - remaining: segwal.BlockSize], make([]byte, remaining))
			segwal.BlockManager.Write(segwal.HotSegment.Path, blockIndex, block)
		}

		n := min(len(data), remaining-7)
		last := n == len(data)

		var typ FragmentType
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

		copy(block[segwal.BlockSize-remaining:], SerializeFragment(Fragment{ 
			Type: typ,
			Payload: payload,
			Hash := crc32.ChecksumIEEE(payload),
			Size := len(payload)
		}))

		BlockManager.WriteBlock(segwal.HotSegment.Path, blockIndex, block)

		if type != Full && type != Last{
			blockIndex += 1
		}

		segwal.SegmentBytes += uint64(7 + n)

		data = data[n:]
		first = false
	}

}