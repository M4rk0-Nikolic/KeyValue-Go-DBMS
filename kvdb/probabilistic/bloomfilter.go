package bloomfilter

import (
	"fmt"
	"encoding/binary"
	"math"
	"errors"
	"github.com/twmb/murmur3"
	"hash"
	"list/List"
)

// https://pkg.go.dev/github.com/twmb/murmur3

type BloomFilter struct {
	M uint
	K uint
	P float64
	HashFunctions []hash.Hash32
	Array []uint64 // With bit packing for least amount of memory usage!
}


func CalculateM(expectedElements uint, falsePositiveRate float64) uint {
	return uint(math.Ceil(float64(expectedElements) * float64(math.Log(falsePositiveRate)) / math.Pow(math.Log(2), float64(2))))
}

func CalculateK(expectedElements uint, m uint) uint {
	return uint(math.Ceil((float64(m) / float64(expectedElements)) * math.Log(2)))
}

func CreateHashFunctions(k uint) ([]hash.Hash32){
	var h []hash.Hash32

}


func CreateBloomFilter(expectedElements uint, falsePositiveRate float64) error{
	if !(0.0 < errorRate, 1.0){
		return errors.New("Error Rate has to be between (0,1).")
	}

	m := CalculateM(expectedElements, falsePositiveRate)
	k := CalculateK(expectedElements, m)








	return nil
}
