package tree_sha256_hash

import "crypto/sha256"
import "github.com/ivan386/go-hash-tree"

func New() *hash_tree.Tree {
	return hash_tree.New(sha256.New(), 1024, []byte{0x00}, []byte{0x01}, true)
}