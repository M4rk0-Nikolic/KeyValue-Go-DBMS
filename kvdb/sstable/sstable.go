package sstable

type SStable struct{
	Data []WALEntry
	Filter *bloomfilter.BloomFilter

}