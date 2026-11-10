package db

import (
	"errors"
	"fmt"
	"log"
)


type DB struct {
	BlockManager *bm.BlockManager
	WAL *wal.SegWAL
}


func (d *DB) Get(key string) ([]byte, error){

	return []byte(), nil
}

func (d *DB) PUT(key string, value any) error{
	var valueBytes []byte

	switch v := valueBytes.(type){
	case string:
		valueBytes = []byte(v)
	case []byte:
		valueBytes = v
	default:
		return errors.New("Incorrect type passed to value parameter. PUT only supports string or []byte value types.")
	}
	
	return nil
}

func (d *DB) Delete(key any) error{
	var keyBytes []byte

	switch v := key.(type){
	case string:
		keyBytes = []byte(v)
	case []byte:
		keyBytes = v
	default:
		return errors.New("Incorrect type passed. Delete only supports string or []byte types.")
	}
	
	return nil
}