package sstable

import (
	"crypto/sha256"
)

type Hash [32]byte
type Block []byte

type Hashable interface {
	hash() Hash
}

type Node struct {
	Left Hashable
	Right Hashable
}

func hash(data []byte) Hash{
	return sha256.Sum256(data)
}

func (b Block) hash() Hash {
	return hash([]byte(b)[:])
}

func (n Node) hash() Hash {
	var l, r Hash 
	l = n.Left.hash()
	r = n.Right.hash()
	return hash(append(l[:], r[:]...))
}

type EmptyBlock struct {}

func (_ EmptyBlock) hash() Hash {
	return [32]byte{}
}

func BuildTree(parts []Hashable) []Hashable {
	var nodes []Hashable
	var i int
	for i = 0; i < len(parts); i += 2 {
		if i + 1 < len(parts) {
			nodes = append(nodes, Node{Left: parts[i], Right: parts[i+1]})
		} else {
			nodes = append(nodes, Node{Left: parts[i], Right: EmptyBlock{}})
		}
	}
	if len(nodes) == 1 {
		return nodes
	} else {
		return BuildTree(nodes)
	}
}

type MerkleTree struct {
	Nodes []Hash
}

func (m MerkleTree) Serialize() []byte {
    result := make([]byte, 0)

    result = binary.AppendUvarint(result, uint64(len(m.Nodes)))

    for _, h := range m.Nodes {
        result = append(result, h[:]...)
    }

    return result
}

func DeserializeMerkleTree(data []byte) (MerkleTree, error) {
    nodeCount, n := binary.Uvarint(data)
    if n <= 0 {
        return MerkleTree{}, errors.New("Invalid node count")
    }
    data = data[n:]

    if len(data) < int(nodeCount)*32 {
        return MerkleTree{}, errors.New("Invalid merkle tree data")
    }

    nodes := make([]Hash, nodeCount)

    for i := range nodes {
        copy(nodes[i][:], data[i*32:(i+1)*32])
    }

    return MerkleTree{
        Nodes:     nodes
    }, nil
}