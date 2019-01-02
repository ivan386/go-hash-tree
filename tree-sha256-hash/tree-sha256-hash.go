package tree_sha256_hash

import sha256 "crypto/sha256"
import hash_tree "github.com/ivan386/go-hash-tree"
import hash_list "github.com/ivan386/go-hash-list"

func New() *hash_tree.Tree {
	return hash_tree.New(sha256.New(), 1024, []byte{0x00}, []byte{0x01})
}

func NewHashList(level uint) *hash_list.List {
	return hash_list.New(hash_tree.New(sha256.New(), 1024, []byte{0x00}, []byte{0x01}), 1024<<level)
}
