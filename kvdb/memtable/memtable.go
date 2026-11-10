package memtable

type Memtable interface {
    Put(record Record) error
    Get(key string) (Record, bool)
    Delete(key string) error
    Iterator() Iterator
}