package tree_pairs

import "hash"

type TreePairs struct {
	hasher      hash.Hash
	pair_prefix []byte

	levels     [][]byte
	index      uint
	state      uint
	last_hash  []byte
	hash_count uint
}

func New(hasher hash.Hash, pair_prefix []byte) *TreePairs {
	tree := new(TreePairs)
	tree.hasher = hasher
	tree.pair_prefix = pair_prefix
	return tree
}

func (tree *TreePairs) BlockSize() int {
	return tree.hasher.Size()
}

func (tree *TreePairs) Size() int {
	return tree.hasher.Size()
}

func (tree *TreePairs) Reset() {
	tree.hasher.Reset()
	tree.levels = nil
	tree.index = 0
	tree.state = 0
	tree.last_hash = nil
	tree.hash_count = 0
}

func (tree *TreePairs) Sum(in []byte) []byte {

	for tree.state > 0 {
		if tree.state&1 == 1 {
			tree.last_hash = tree.pairHash(tree.levels[tree.index], tree.last_hash)
			tree.levels[tree.index] = nil
		}
		tree.state >>= 1
		tree.index += 1
	}

	return append(in, tree.last_hash...)
}

func (tree *TreePairs) Write(hashes []byte) (int, error) {
	writen := len(hashes)
	hash_size := tree.BlockSize()

	for len(hashes) >= hash_size {
		tree.AppendHash(hashes[:hash_size])
		hashes = hashes[hash_size:]
	}

	return writen, nil
}

func (tree *TreePairs) AppendHash(hashes ...[]byte) {
	for _, hash_value := range hashes {
		tree.state = tree.hash_count
		tree.index = 0

		for tree.state&1 == 1 {
			hash_value = tree.pairHash(tree.levels[tree.index], hash_value)
			tree.levels[tree.index] = nil
			tree.state >>= 1
			tree.index += 1
		}

		if len(tree.levels) == int(tree.index) {
			tree.levels = append(tree.levels, hash_value)
		} else {
			tree.levels[tree.index] = hash_value
		}

		tree.last_hash = hash_value
		tree.hash_count += 1
	}
}

func (tree *TreePairs) pairHash(left []byte, right []byte) []byte {
	tree.hasher.Reset()
	tree.hasher.Write(tree.pair_prefix)
	tree.hasher.Write(left)
	tree.hasher.Write(right)
	return tree.hasher.Sum(nil)
}
