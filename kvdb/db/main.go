package main

import (
	"fmt"
)

func main() {
	manager := bm.NewBlockManager()
	wal := wal.NewSegWAL()

	db := db.DB{
		BlockManager: manager
		WAL: wal
	}
	fmt.Println("Hello, World!")
}