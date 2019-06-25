package hash_tree

import "hash"
import tree_pairs "github.com/ivan386/go-hash-tree/tree-pairs"

type Tree struct {
	hasher      hash.Hash
	pairs       hash.Hash
	block_size  int
	data_prefix []byte

	block_len int
	last_hash []byte
}

func NewDefault(hasher hash.Hash) *Tree {
	return New(hasher, 1024, []byte{0x00}, []byte{0x01})
}

func New(hasher hash.Hash, block_size int, data_prefix []byte, pair_prefix []byte) *Tree {
	tree := new(Tree)
	tree.hasher = hasher
	tree.pairs = tree_pairs.New(hasher, pair_prefix, hasher.Size())

	tree.block_size = block_size
	tree.data_prefix = data_prefix
	return tree
}

func (tree *Tree) BlockSize() int {
	return tree.block_size
}

func (tree *Tree) Size() int {
	return tree.hasher.Size()
}

func (tree *Tree) Reset() {
	tree.hasher.Reset()
	tree.pairs.Reset()
	tree.block_len = 0
	tree.last_hash = nil
}

func (tree *Tree) Write(data []byte) (int, error) {
	writen := len(data)

	if tree.block_len > 0 {
		part := tree.block_size - tree.block_len
		if part <= len(data) {
			tree.hasher.Write(data[:part])
			tree.pairs.Write(tree.hasher.Sum(nil))
			tree.block_len = 0
			data = data[part:]
		}
	}

	for len(data) >= tree.block_size {
		tree.hasher.Reset()
		tree.hasher.Write(tree.data_prefix)
		tree.hasher.Write(data[:tree.block_size])
		tree.pairs.Write(tree.hasher.Sum(nil))
		data = data[tree.block_size:]
	}

	if len(data) > 0 {
		if tree.block_len == 0 {
			tree.hasher.Reset()
			tree.hasher.Write(tree.data_prefix)
		}
		tree.hasher.Write(data)
		tree.block_len += len(data)
	}

	return writen, nil
}

func (tree *Tree) Sum(in []byte) []byte {
	if tree.block_len > 0 {
		tree.pairs.Write(tree.hasher.Sum(nil))
		tree.block_len = 0
	}

	if len(tree.last_hash) == 0 {
		tree.last_hash = tree.pairs.Sum(nil)
		if len(tree.last_hash) == 0 {
			tree.hasher.Reset()
			tree.hasher.Write(tree.data_prefix)
			tree.last_hash = tree.hasher.Sum(nil)
		}
	}

	return append(in, tree.last_hash...)
}
