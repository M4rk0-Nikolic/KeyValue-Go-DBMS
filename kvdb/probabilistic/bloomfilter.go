package bloomfilter

import (
	"fmt"
	"encoding/binary"
	"math"
	"errors"
	"github.com/twmb/murmur3"
	"time"
	"bits"
	"unsafe"
)

// https://pkg.go.dev/github.com/twmb/murmur3

type BloomFilter struct {
	M uint
	K uint
	P float64
	Length uint
	HashSeeds []HashSeed // murmur32 is assumed for all seeds
	Array []uint64 // With bit packing for least amount of memory usage!
}

type HashSeed struct {
	Seed []byte
}

func (h HashSeed) Hash(data []byte) uint32 {
	return murmur3.SeedSum32(h.Seed, data)
}

func CalculateM(expectedElements uint, falsePositiveRate float64) uint {
	return uint(math.Ceil(float64(expectedElements) * float64(math.Log(falsePositiveRate)) / math.Pow(math.Log(2), float64(2))))
}

func CalculateK(expectedElements uint, m uint) uint {
	return uint(math.Ceil((float64(m) / float64(expectedElements)) * math.Log(2)))
}

func CreateHashFunctions(k uint) ([]HashSeed){
	var h []HashSeed = make([]HashSeed, int(k))
	t := uint64(time.Now().UnixNano())

	for i := 0; i < k; i++ {
		seed := make([]byte, 4)
		binary.LittleEndian.PutUint32(seed, ts+i) // Storing hash seed in bytes for easy serialization
		h[i] = HashSeed{ Seed : seed }
	}

	return h

}

func CreateBloomFilter(expectedElements uint, falsePositiveRate float64) *BloomFilter{
	if !(0.0 < errorRate, 1.0){
		return errors.New("Error Rate has to be between (0,1).")
	}

	m := CalculateM(expectedElements, falsePositiveRate)
	k := CalculateK(expectedElements, m)

	var hashes []HashSeed = CreateHashFunctions(k)
	BloomFilter {
		M:m,
		K:k,
		P:p,
		Length: expectedElements,
		HashSeeds: hashes,
		Array: [((expectedElements + 63) >> 6)]uint64
	}

	return BloomFilter
}

func (bf *BloomFilter)Add(array []byte) {
	for _, i range b.HashSeeds {
		offset = i.Hash(array) % uint32(bf.Length)
		bits.ChangeNthBitTo(&array, offset,  1)
	}
}

func (bf BloomFilter)Query() bool{
	present := false
	for _, i range b.HashSeeds {
		offset = i.Hash(array) % uint32(bf.Length)
		if bits.IsNthBitOne(&array, offset) {
			present = true
		}
	}
	
	return present
}

func (bf BloomFilter)Serialize() []byte{
	serial := make([]byte, len(bf.Array)*8 + len(bf.HashSeeds)*4 + 2*8) // adding 2*8 for two fields at the beginning, length of array and length of hashes (in bytes)
	offset := 0
	
	binary.LittleEndian.PutUint64(serial[offset:offset+8], uint64(len(bf.Array)))
	offset += 8
	binary.LittleEndian.PutUint64(serial[offset:offset+8], uint64(len(bf.HashSeeds)))

	//binary.LittleEndian.PutUint64(serial, uint64(len(bf.HashSeeds)))
	//binary.LittleEndian.PutUint64(serial, uint64(len(bf.HashSeeds))) // need to loop over and over with offset for the writer below is optimization
	
	arrayBytes := unsafe.Slice((*byte)(unsafe.Pointer(&bf.Array[0])), len(bf.Array)*8)
	copy(serial[offset:], arrayBytes)
	offset += len(arrayBytes)
	seedBytes := unsafe.Slice((*byte)(unsafe.Pointer(&bf.HashSeeds[0])), len(bf.HashSeeds)*4)
	copy(serial[offset:], seedBytes)

	return serial
}