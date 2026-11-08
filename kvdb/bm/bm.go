package bm

import (
	"os"
	"errors"
	"time"
	"container/list"
	"runtime" // To be used to detrmine whether running on Linux
	"syscall" // fdatasync and fsync, difference being fsync syncs meta data as well, fdatasync syncs data and only critical meta data
)

type BlockID struct {
	Path string
	Position uint64
}

type Block struct {
	ID BlockID
	Data []byte
}


// Map for O(1) point search (Blocks)
// Doubly Linked List for O(1) insertions (LRUEvictionOrder)
type BlockManager struct {
	BlockSizeBytes uint64
	BlockCacheSizeKB uint64
	Blocks map[BlockID]*list.Element
	LRUEvictionOrder list.List
}

func NewBlockManager() *BlockManager {
	var defBlockSize uint64
	var defBlockCacheSize uint64
	
	// Read config yaml file here, either have a seperate .go file for reading .yaml file and somehow share it as a global variable OR read when needed; also which lib for reading yaml?
	// Also need to perfom check that BlockSizeBytes is multiple of BlockCacheSizeKB * 1024
	// and that cache is non-zero

    return &BlockManager{
        BlockSizeBytes: 4096,
        BlockCacheSizeKB: 1024,

        Blocks: make(map[BlockID]*list.Element),
    }
}

func (bm *BlockManager) Add(block Block) {
	id := block.ID
	
	// Not in cache
	if element, present := bm.Blocks[id]; !present {
		// Full Block cache
		if bm.LRUEvictionOrder.Len() == int(bm.BlockCacheSizeKB * 1024 / bm.BlockSizeBytes){ // Assumes that BlockSizeBytes is a multiple of bm.BlockCacheSizeKB * 1024
			last := bm.LRUEvictionOrder.Back()
			old := last.Value.(Block)
			delete(bm.Blocks, old.ID)
			bm.LRUEvictionOrder.Remove(bm.LRUEvictionOrder.Back())
		}
		bm.Blocks[id] = bm.LRUEvictionOrder.PushFront(block)
	} else {
		// In cache
		element.Value = block
		bm.LRUEvictionOrder.MoveToFront(element)
	}

}

func (bm *BlockManager)IsCached(blockid BlockID) bool{
	if _, isCached := bm.Blocks[block.ID]; isCached {
		return true
	return false
	}
}

func (bm *BlockManager) ReadBlock(path string, blockNumber uint64) ([]byte, error){
	blockid := BlockID{path, blockNumber}
	if present := bm.IsCached(blockid); present{
		element := bm.Blocks[blockid]
		block := element.Value.(Block) // type assertion? Need to convert to Block here in order to access Data field, elements of the array are otherwise type any
		
		bm.Add(Block { ID: blockid, Data: block.Data })
		return block.Data, nil
	}

	file, error := os.OpenFile(path, os.O_RDONLY, 0644) //0644 RW for owner (us), R only for other apps
	
	if error != nil {
		log.Println("Failed to open file in order to read block.")
		return []byte{}, errors.New("Failed to open file.")
	}
	defer file.Close()

	offset := int64(bm.BlockSizeBytes*blockNumber)
	block := make([]byte, int64(bm.BlockSizeBytes))

	_, err := file.ReadAt(block, offset)
	if err != nil {
		log.Println("Failed to read block.")
        return []byte{}, errors.New("Failed to read file.")
    }

	bm.Add(Block { ID: blockid, Data: block })

	return block, nil
}


// In order to create file call WriteBlock(path string, 0, data []byte)
func (bm *BlockManager) WriteBlock(path string, blockNumber uint64, data []byte) error{
	file, erorr := os.OpenFile(path, os.O_CREATE | os.O_WRONLY, 0644)

	blockid := BlockID{path, blockNumber}

	if error != nil {
		log.Println("Failed to open file in order to write block.")
		return errors.New("Failed to write to file.")
	}

	defer file.Close()

	offset := int64(bm.BlockSizeBytes*blockNumber)
	block := make([]byte, bm.BlockSizeBytes)

	file.WriteAt(block, int64(offset))
	syscall.Fdatasync(int(file.Fd()))

	if err != nil {
		log.Println("Failed to write block.")
        return errors.New("Failed to write file.")
    }

	if _, present := bm.Blocks[blockid]; present {
		bm.Add(Block{ ID: blockid, Data: data}) // If in cache update, but do not insert if not in cache
	}

	return nil

}