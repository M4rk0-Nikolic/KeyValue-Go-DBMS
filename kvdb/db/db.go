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


func (d *DB) Get(key any) ([]byte, error){
	var keyBytes []byte

	switch v := key.(type){
	case string:
		keyBytes = []byte(v)
	case []byte:
		keyBytes = v
	default:
		return nil, errors.New("Incorrect type passed. Get only supports string or []byte types.")
	}
	
	return keyBytes, nil
}

func (d *DB) PUT(key any, value []byte) error{
	var keyBytes []byte

	switch v := key.(type){
	case string:
		keyBytes = []byte(v)
	case []byte:
		keyBytes = v
	default:
		return errors.New("Incorrect type passed. PUT only supports string or []byte types.")
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