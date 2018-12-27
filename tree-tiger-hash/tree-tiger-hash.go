package tree_tiger_hash

import "github.com/cxmcc/tiger"
import "github.com/ivan386/go-hash-tree"

func New() *hash_tree.Tree {
	return hash_tree.New(tiger.New(), 1024, []byte{0x00}, []byte{0x01}, true)
}