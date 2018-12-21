package hash_tree

import "hash"
import "github.com/cxmcc/tiger"

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

func New() *Tree {
	return NewTree(tiger.New())
}

func NewTree(hasher hash.Hash) *Tree {
	return NewCustomTree(hasher, 1024, []byte{0x00}, []byte{0x01}, true)
}

func NewCustomTree(hasher hash.Hash, block_size int, data_prefix []byte, pair_prefix []byte, skip_empty bool) *Tree {
	t := new(Tree)
	t.hasher = hasher
	t.block_size = block_size
	t.data_prefix = data_prefix
	t.pair_prefix = pair_prefix
	t.skip_empty = skip_empty
	return t
}

func (t *Tree) BlockSize() int {
	return t.block_size
}

func (t *Tree) Size() int {
	return t.hasher.Size()
}

func (t *Tree) Reset() {
	t.hasher.Reset()
	t.block_len = 0
	t.levels = nil
	t.index = 0
	t.state = 0
	t.last_hash = nil
	t.block_index = 0
}

func (t *Tree) Write(data []byte) (int, error) {
	writen := len(data)

	if t.block_len > 0 {
		part := t.block_size - t.block_len
		if part <= len(data) {
			t.hasher.Write(data[:part])
			t.AppendHash(t.hasher.Sum(nil))
			t.block_len = 0
			data = data[part:]
		}
	}

	for len(data) >= t.block_size {
		t.hasher.Reset()
		t.hasher.Write(t.data_prefix)
		t.hasher.Write(data[:t.block_size])
		t.AppendHash(t.hasher.Sum(nil))
		data = data[t.block_size:]
	}

	if len(data) > 0 {
		if t.block_len == 0 {
			t.hasher.Reset()
			t.hasher.Write(t.data_prefix)
		}
		t.hasher.Write(data)
		t.block_len += len(data)
	}

	return writen, nil
}

func (t *Tree) Sum(in []byte) []byte {
	if t.block_len > 0 {
		t.AppendHash(t.hasher.Sum(nil))
	}
	
	if t.block_index == 0 {
		if t.block_len == 0 {
			t.hasher.Reset()
			t.hasher.Write(t.data_prefix)
		}
		t.last_hash = t.hasher.Sum(nil)
		t.block_index = 1
	} else {
		for t.state > 0 {
			if t.state&1 == 1 {
				t.last_hash = t.pairHash(t.levels[t.index], t.last_hash)
			} else if !t.skip_empty {
				t.last_hash = t.pairHash(t.last_hash, t.last_hash)
			}
			t.state >>= 1
			t.index += 1
		}
	}
	
	return append(in, t.last_hash...)
}

func (t *Tree) AppendHash(hashes ...[]byte) {
	for _, hash_value := range hashes {
		t.state = t.block_index
		t.index = 0
		
		for t.state&1 == 1 {
			hash_value = t.pairHash(t.levels[t.index], hash_value)
			t.state >>= 1
			t.index += 1
		}

		if len(t.levels) == int(t.index) {
			t.levels = append(t.levels, hash_value)
		} else {
			t.levels[t.index] = hash_value
		}

		t.last_hash = hash_value
		t.block_index += 1
	}
}

func (t *Tree) pairHash(left []byte, right []byte) []byte {
	t.hasher.Reset()
	t.hasher.Write(t.pair_prefix)
	t.hasher.Write(left)
	t.hasher.Write(right)
	return t.hasher.Sum(nil)
}
