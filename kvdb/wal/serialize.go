package wal

import (
	"encoding/binary"
)

const (
	CrcSize = 4
	SizeSize = 2
	TypeSize = 1
	Total = CrcSize + SizeSize + TypeSize

	TimestampSize = 8
	TombstoneSize = 1
	KeysizeSize = 8
	NumberofopsSize = 2
	ValuesizeSize = 8
	TcSize = 8
)

func (op Operation)Serialize() []byte{
	length := len(op.Key) + len(op.Value) + TombstoneSize + KeysizeSize + ValuesizeSize
	serialized := make([]byte, length)
	offset := 0

	if op.Tombstone {
        serialized[offset] = 1
    }
    offset += TombstoneSize

	binary.LittleEndian.PutUint64(serialized[offset], op.KeySize)
	offset += KeysizeSize
	
	binary.LittleEndian.PutUint64(serialized[offset], op.ValueSize)
	offset += ValuesizeSize

	offset += copy(serialized[offset:], op.Key)
	offset += copy(serialized[offset:], op.Value)

	return serialized
}

func DeserializeOperation(data []byte) (Operation, int, error){
	var operation Operation
	offset := 0

	if len(data) < TombstoneSize + KeysizeSize + ValuesizeSize {
		return operation, 0, errors.New("Incomplete operation header")
	}

	operation.Tombstone = data[offset] != 0
	offset += TombstoneSize

	operation.KeySize = binary.LittleEndian.Uint64(data[offset])
	offset += KeysizeSize

	operation.ValueSize = binary.LittleEndian.Uint64(data[offset])
	offset += ValuesizeSize

	if operation.KeySize > uint64(len(data) - offset) {
		return operation, 0, errors.New("Incomplete key")
	}

	operation.Key = append([]byte{}, data[offset:offset + int(operation.KeySize)]...)
	offset += int(operation.KeySize)

	if operation.ValueSize > uint64(len(data) - offset) {
		return operation, 0, errors.New("Incomplete value")
	}

	operation.Value = append([]byte{}, data[offset:offset+int(operation.ValueSize)]...)
	offset += int(operation.ValueSize)

	return operation, offset, nil
}


func (entry Entry)Serialize() []byte{
	serialized := make([]byte, TimestampSize + NumberofopsSize)
	offset := 0

	binary.LittleEndian.PutUint64(serialized[offset], entry.Timestamp)
	offset += TimestampSize

	binary.LittleEndian.PutUint16(serialized[offset], entry.NumberOfOperations)
	offset += NumberofopsSize

	for _, op := range entry.Operations{
		serialized = append(serialized, op.Serialize()...)
	}

	tc := make([]byte, TcSize)
    binary.LittleEndian.PutUint64(tc, entry.TC)
    serialized = append(serialized, tc...)

	return serialized
}

func DeserializeEntry(data []byte) (Entry, int, error) {
	var entry Entry
	offset := 0

	if len(data) < TimestampSize + NumberofopsSize {
		return entry, 0, errors.New("Incomplete entry header")
	}

	entry.Timestamp = binary.LittleEndian.Uint64(data[offset])
	offset += TimestampSize

	entry.NumberOfOperations = binary.LittleEndian.Uint16(data[offset])
	offset += NumberofopsSize

	for i := 0; i < int(entry.NumberOfOperations); i++ {
		op, n, err := DeserializeOperation(data[offset:])
		if err != nil {
			return entry, 0, err
		}

		entry.Operations = append(entry.Operations, op)
		offset += n
	}

	if len(data) - offset < TcSize {
		return entry, 0, errors.New("Missing trailing canary")
	}

	entry.TC = binary.LittleEndian.Uint64(data[offset])
	offset += TcSize

	return entry, offset, nil
}

func SerializeFragment(fragment Fragment) []byte{
	result := make([]byte, fragment.Size + SizeSize + CrcSize)
	
	binary.LittleEndian.PutUint32(result, fragment.Hash)
	binary.LittleEndian.PutUint16(result[CrcSize:], fragment.Size)

	result[CrcSize+SizeSize] = byte(fragment.Type)

	copy(result[SizeSize + CrcSize:], fragment.Payload)

	return result
}

func DeserializeFragment(data []byte) (Fragment, int, error) {
    var fragment Fragment
    offset := 0

    headerSize := CrcSize + SizeSize + TypeSize

    if len(data) < headerSize {
        return fragment, 0, errors.New("Incomplete fragment header")
    }

    fragment.Hash = binary.LittleEndian.Uint32(data[offset:])
    offset += CrcSize

    fragment.Size = binary.LittleEndian.Uint16(data[offset:])
    offset += SizeSize

    fragment.Type = FragmentType(data[offset])
    offset += TypeSize

    if int(fragment.Size) > len(data) - offset {
        return fragment, 0, errors.New("Incomplete fragment payload")
    }

    fragment.Payload = make([]byte, int(fragment.Size))

	copy(fragment.Payload, data[offset:offset+int(fragment.Size)])
    offset += int(fragment.Size)

    return fragment, offset, nil
}