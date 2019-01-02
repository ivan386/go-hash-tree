package tree_any_hash

import hash "hash"
import hash_tree "github.com/ivan386/go-hash-tree"
import hash_list "github.com/ivan386/go-hash-list"

func New(hasher hash.Hash) *hash_tree.Tree {
	return hash_tree.New(hasher, 1024, []byte{0x00}, []byte{0x01})
}

func NewHashList(hasher hash.Hash, level uint) *hash_list.List {
	return hash_list.New(hash_tree.New(hasher, 1024, []byte{0x00}, []byte{0x01}), 1024<<level)
}
