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

type WALEntry struct {
	Timestamp uint64
	Tombstone bool
	KeySize uint64
	ValueSize uint64
	Key []byte
	Value []byte
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

func (segwal SegWAL) SyncHotSegment() error{
	if file, err := os.Openfile(segwal.HotSegment, O_RDONLY, 0644); err == nil {
		file.Sync()
		file.Close()
		return nil
	}
	return err
}

func (wale WALEntry)Serialize() []byte{
	length := TOTAL + len(wale.Key) + len(wale.Value)
	serialized := make([]byte, length)
	offset := 0
	binary.LittleEndian.PutUint64(wale.Timestamp, serialized[offset:offset+8])
	offset += TIMESTAMP_SIZE
	serialized[offset] = byte(wale.Tombstone)
	offset += TOMBSTONE_SIZE
	binary.LittleEndian.PutUint64(wale.KeySize, serialized[offset:offset+8])
	offset += KEY_SIZE_SIZE
	binary.LittleEndian.PutUint64(wale.ValueSize, serialized[offset:offset+8])
	offset += VALUE_SIZE_SIZE
	offset += copy(serialized[offset:], wale.Key)
	offset += copy(serialized[offset:], wale.Value)
	binary.LittleEndian.PutUint64(wale.TC, serialized[offset:offset+8])

	return serialized
}

func GetSegPath(num uint16) {
	str := strconv.FormatUint(uint64(num), 10)

	if len(str) < 5 {
		str = "../data/wal/wal" + strings.Repeat("0", 5-len(str)) + str + ".log"
	}

	return str
}

func CRC32(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
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

	var wal_logs []WALSegment = make([]WALSegment, 0)

	for _, entry := range entries{
		if !entry.IsDir(){
			if (strings.HasPrefix(entry.Name(), "wal_")) &&  (strings.HasSufix(entry.Name(), ".log")){
				wal_logs = append(wal_logs, WALSegment{ SequenceNumber: uint16(_) Path: entry })
			}
		}
	}

	// Need to implement other things before finishing this

	SegWAL{
		BlockManager: bm
		Segments : wal_logs
	}

}

func (segwal *SegWAL) SwapSegment() {
	segwal.SyncHotSegment()
	old_num := segwal.HotSegment.SequenceNumber
	new_seg := WALSegment{ SequenceNumber: uint16(old_num + 1) Path: GetSegPath(old_num + 1) }
	segwal.Segments = append(segwal.Segments, segwal.HotSegment)
	segwal.HotSegment =  new_seg
	os.CreateFile(new_seg.Path)
}

func SerializeFragment(fragment WALFragment) {
	result := make([]byte, fragment.Size+2+4)
	
	binary.LittleEndian.PutUint32(result, fragment.Hash)
	binary.LittleEndian.PutUint16(result[4:], fragment.Size)

	copy(result[6:], fragment.Payload)

	return result
}

func (segwal *SegWAL) AppendRecordUnBuffered(record record.Record) error{
	if segwal.HotSegment == nil | len(segwal.Segments) == 0 {
		newseg := WALSegment {
			SequenceNumber: 1,
			Path: GetSegPath(1)
		}
		segwal.Segments = append(segwal.Segments, newseg)
		segwal.HotSegment = newseg
	}

	file, err := os.OpenFile(segwal.HotSegment.Path, O_CREATE | O_APPEND, 0644)
	if err != nil{
		log.Fatal("Can not open log file.")
	}
	defer file.Close()

	first := true
	walentry := WALEntry {
		Timestamp: time.Now().Unix()
		Tombstone: record.Tombstone
		Key: record.Key
		Value: record.Value
		KeySize: len(record.Key)
		ValueSize: len(record.Value)
		TC: 0xDEADBEEFFEEDFACE
		}
	data := walentry.Serialize()

	total_size := segwal.BlockSize * wal.BlocksInSegment

	for len(data) > 0 {
		if segwal.SegmentBytes == total_size {
			if err := segwal.SwapSegment(); err != nil {
				return err
			}
		}

		remaining := total_size - (wal.SegmentBytes % segwal.BlockSize)

		if remaining < TOTAL {
			if file, err := os.OpenFile(segwal.HotSegment.Path, O_APPEND, 0644); err != nil{
				return err
			}
			file.Write(make([]byte, remaining))
			file.Sync()
			file.Close()
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

		SerializeFragment(WALFragment{ 
			Type: typ,
			Payload: payload,
			Hash := CRC32(payload),
			Size := len(payload)
		})

		segwal.SegmentBytes += uint64(7 + n)

		data = data[n:]
		first = false
	}

}


func (segwal *SegWAL) ReadSegment(record record.Record) {

}


func (segwal *SegWAL) AppendRecordBuffered(record record.Record) {

}