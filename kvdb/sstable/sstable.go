package sstable

type SStable struct{
	Data []WALEntry
	Filter *bloomfilter.BloomFilter
	Index 
}

type Index struct{
	Key []byte
	Offset int
}