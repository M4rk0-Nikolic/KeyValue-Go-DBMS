package skiplist

import (
	"math/rand"
)

const (
	Probability float64 = 0.5
	MaxLevelHeight int = 16
)

type Node struct {
	RightNode []*Node
	Key string
	Value []byte
}

type SkipList struct {
	Head *Node
}

func CreateSkipList() *SkipList{
	return &SkipList{
		Head: &Node{
			RightNode: nil,
			Key: "", // Key "" signifies -inf while RightNode: nil signifies +inf, so here it is the head and the tail
			Value: nil
		}
	}
}

// Here the second returned value is a bool to signify whether a Key was found or not, if it wasn't return the key before it
func (skiplist SkipList)Search(targetKey string) (*Node, bool){
	var current_node *Node = skiplist.Head

	level := len(skiplist.Head.RightNode) - 1

	for ; ; {
		if current_node.RightNode[level].Key == targetKey {
			return current_node.RightNode[level], true
		} else if current_node.RightNode[level].Key > targetKey{
			if level == 1 {
				return current_node, false
			}
			level--
			continue
		} else if current_node.RightNode[level].Key < targetKey{
			current_node = current_node.RightNode[level]
		}
	}

}

func (skiplist *SkipList)Append(Key string, Value []byte) bool{
	levels := 0

	for ; rand.Float64() <= Probability ; levels++ {
		if levels == MaxLevelHeight{
			break
		}
	}

	node, wasfound := skiplist.Search(Key)

	if wasfound {
		return false
	}

	left_neighbours := make([]*Node, len(skiplist.Head.RightNode))
	current_node := skiplist.Head
	for level := len(skiplist.Head.RightNode) - 1; level >= 0; level-- {
		for current_node.RightNode[level] != nil && current_node.RightNode[level].Key < Key {
			current_node = current_node.RightNode[level]
		}
		update[level] = current_node
	}

	newnode := Node{
		RightNode: make([]*Node, levels)
		Key: Key,
		Value: Value
	}

	for i := 0; i < levels; i++ {
		newnode.RightNode[i] = update[i].RightNode[i]
		update[i].RightNode[i] = newnode
	}

	return true
}
