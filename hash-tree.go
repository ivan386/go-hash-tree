package hash_tree

import "hash"

type Tree struct {
	hasher      hash.Hash
	block_size  int
	data_prefix []byte
	pair_prefix []byte
	skip_empty  bool

	block_len   int
	levels      [][]byte
	index       uint
	state       uint
	last_hash   []byte
	block_index uint
}

func New(hasher hash.Hash) *Tree {
	return NewCustom(hasher, 1024, []byte{0x00}, []byte{0x01}, true)
}

func NewCustom(hasher hash.Hash, block_size int, data_prefix []byte, pair_prefix []byte, skip_empty bool) *Tree {
	tree := new(Tree)
	tree.hasher = hasher
	tree.block_size = block_size
	tree.data_prefix = data_prefix
	tree.pair_prefix = pair_prefix
	tree.skip_empty = skip_empty
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
	tree.block_len = 0
	tree.levels = nil
	tree.index = 0
	tree.state = 0
	tree.last_hash = nil
	tree.block_index = 0
}

func (tree *Tree) Write(data []byte) (int, error) {
	writen := len(data)

	if tree.block_len > 0 {
		part := tree.block_size - tree.block_len
		if part <= len(data) {
			tree.hasher.Write(data[:part])
			tree.AppendHash(tree.hasher.Sum(nil))
			tree.block_len = 0
			data = data[part:]
		}
	}

	for len(data) >= tree.block_size {
		tree.hasher.Reset()
		tree.hasher.Write(tree.data_prefix)
		tree.hasher.Write(data[:tree.block_size])
		tree.AppendHash(tree.hasher.Sum(nil))
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
		tree.AppendHash(tree.hasher.Sum(nil))
	}

	if tree.block_index == 0 {
		if tree.block_len == 0 {
			tree.hasher.Reset()
			tree.hasher.Write(tree.data_prefix)
		}
		tree.last_hash = tree.hasher.Sum(nil)
		tree.block_index = 1
	} else {
		for tree.state > 0 {
			if tree.state&1 == 1 {
				tree.last_hash = tree.pairHash(tree.levels[tree.index], tree.last_hash)
			} else if !tree.skip_empty {
				tree.last_hash = tree.pairHash(tree.last_hash, tree.last_hash)
			}
			tree.state >>= 1
			tree.index += 1
		}
	}

	return append(in, tree.last_hash...)
}

func (tree *Tree) AppendHashList(hashes []byte) {
	hash_size := tree.Size()
	for len(hashes) >= hash_size {
		tree.AppendHash(hashes[:hash_size])
		hashes = hashes[hash_size:]
	}
}

func (tree *Tree) AppendHash(hashes ...[]byte) {
	for _, hash_value := range hashes {
		if len(hash_value) != tree.Size() {
			panic("wrong hash size")
		}
		tree.state = tree.block_index
		tree.index = 0

		for tree.state&1 == 1 {
			hash_value = tree.pairHash(tree.levels[tree.index], hash_value)
			tree.state >>= 1
			tree.index += 1
		}

		if len(tree.levels) == int(tree.index) {
			tree.levels = append(tree.levels, hash_value)
		} else {
			tree.levels[tree.index] = hash_value
		}

		tree.last_hash = hash_value
		tree.block_index += 1
	}
}

func (tree *Tree) pairHash(left []byte, right []byte) []byte {
	tree.hasher.Reset()
	tree.hasher.Write(tree.pair_prefix)
	tree.hasher.Write(left)
	tree.hasher.Write(right)
	return tree.hasher.Sum(nil)
}
