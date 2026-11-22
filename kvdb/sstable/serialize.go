package sstable

import (
    "encoding/binary"
	bm "github.com/M4rk0-Nikolic/KeyValue-Go-DBMS"
	"github.com/M4rk0-Nikolic/KeyValue-Go-DBMS/config"
)

var BlockSize uint64 = config.Global.BlockManager.BlockSizeInKB * 1024

type SSTable struct{
	Path string
	BlockManager *bm.BlockManager
	CurrentOffset uint64
}

type ValueEntry struct {
	Timestamp uint64
	Tombstone bool
    ValueLength uint64
	Value []byte
}

func (value ValueEntry)Serialize struct {
    buffer := []byte{}

    buffer = binary.AppendUvarint(buffer, value.Timestamp)

    if value.Tombstone {
        buffer = append(buffer, byte(1))
    } else{
        buffer = append(buffer, byte(0))
    }

    buffer = append(buffer, value.Value...)
    return buffer
}


func DeserializeValueEntry(data []byte) (ValueEntry, error) {
    var value ValueEntry

    timestamp, n := binary.Uvarint(data)
    if n <= 0 {
        return value, errors.New("Invalid timestamp")
    }
    value.Timestamp = timestamp
    data = data[n:]

    if len(data) < 1 {
        return value, errors.New("Missing tombstone")
    }

    value.Tombstone = data[0] == 1
    data = data[1:]

    valueLength, n = binary.Uvarint(data)
    if n <= 0 {
        return value, errors.New("Invalid value length")
    }
    data = data[n:]
    value.ValueLength = valueLength

    if valueLength > uint64(len(data)) {
        return value, errors.New("Incomplete value")
    }

    value.Value = make([]byte, valueLength)
    copy(value.Value, data[:])

    return value, nil
}

type IndexEntry struct {
	Offset uint64
    KeyLength uint64
	Key []byte
}

func (index IndexEntry)Serialize struct {
    buffer := []byte{}

    buffer = binary.AppendUvarint(buffer, index.Offset)

    buffer = binary.AppendUvarint(buffer, index.KeyLength)

    buffer = append(buffer, value.Value...)
    return buffer
}

func DeserializeIndexEntry(data []byte) (IndexEntry, error) {
    var value IndexEntry

    offset, n := binary.Uvarint(data)
    if n <= 0 {
        return value, errors.New("Invalid offset")
    }
    value.Offset = offset
    data = data[n:]

    keyLength, n = binary.Uvarint(data)
    if n <= 0 {
        return value, errors.New("Invalid key length")
    }
    value.KeyLength = keyLength
    data = data[n:]

    if keyLength > uint64(len(data)) {
        return value, errors.New("Incomplete key")
    }

    value.Key = make([]byte, keyLength)
    copy(value.Key, data[:])

    return value, nil
}

type Section struct {
    Start uint64
    End uint64
}

type Footer struct {
    DataOffset       Section
    FilterOffset     Section
    SummaryOffset    Section
    IndexOffset      Section
    HashIndexOffset  Section
    MetadataOffset   Section
}

func (footer Footer)Serialize struct {
    buffer := []byte{}

    buffer = binary.AppendUvarint(buffer, footer.DataOffset.Start)
    buffer = binary.AppendUvarint(buffer, footer.DataOffset.End)

    buffer = binary.AppendUvarint(buffer, footer.FilterOffset.Start)
    buffer = binary.AppendUvarint(buffer, footer.FilterOffset.End)

    buffer = binary.AppendUvarint(buffer, footer.SummaryOffset.Start)
    buffer = binary.AppendUvarint(buffer, footer.SummaryOffset.End)

    buffer = binary.AppendUvarint(buffer, footer.IndexOffset.Start)
    buffer = binary.AppendUvarint(buffer, footer.IndexOffset.End)

    buffer = binary.AppendUvarint(buffer, footer.HashIndexOffset.Start)
    buffer = binary.AppendUvarint(buffer, footer.HashIndexOffset.End)

    buffer = binary.AppendUvarint(buffer, footer.MetadataOffset.Start)
    buffer = binary.AppendUvarint(buffer, footer.MetadataOffset.End)

    return buffer
}

func ReadUvarint(data *[]byte) (uint64, error) {
    value, n := binary.Uvarint(*data)

    if n <= 0 {
        return 0, errors.New("Invalid uint64")
    }

    *data = (*data)[n:]

    return value, nil
}

func DeserializeFooter(data []byte) (Footer, error) {
    var value Footer
    var err error

    value.DataOffset.Start, err = ReadUvarint(&data)
    if err != nil { return value, err }
    value.DataOffset.End, err = ReadUvarint(&data)
    if err != nil { return value, err }

    value.FilterOffset.Start, err = ReadUvarint(&data)
    if err != nil { return value, err }
    value.FilterOffset.End, err = ReadUvarint(&data)
    if err != nil { return value, err }

    value.SummaryOffset.Start, err = ReadUvarint(&data)
    if err != nil { return value, err }
    value.SummaryOffset.End, err = ReadUvarint(&data)
    if err != nil { return value, err }

    value.IndexOffset.Start, err = ReadUvarint(&data)
    if err != nil { return value, err }
    value.IndexOffset.End, err = ReadUvarint(&data)
    if err != nil { return value, err }

    value.HashIndexOffset.Start, err = ReadUvarint(&data)
    if err != nil { return value, err }
    value.HashIndexOffset.End, err = ReadUvarint(&data)
    if err != nil { return value, err }

    value.MetadataOffset.Start, err = ReadUvarint(&data)
    if err != nil { return value, err }
    value.MetadataOffset.End, err = ReadUvarint(&data)
    if err != nil { return value, err }

    if len(data) != 0 {
        return value, errors.New("Unexpected trailing footer data")
    }

    return value, nil
}