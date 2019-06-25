package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	hash_list "github.com/ivan386/go-hash-list"
	hash_tree "github.com/ivan386/go-hash-tree"
	hash_pairs "github.com/ivan386/go-hash-tree/tree-pairs"
	"math/rand"
	"strings"
	"time"
)

// Data block size
const block_size = 1024

// Sha256 hash size in bytes
const hash_size = int(256 / 8)

func get_tree_hasher() *hash_tree.Tree {
	// Use sha256 hash for blocks and pairs hashing
	// Data block size block_size bytes
	// Data block prefix 1 byte = 0x00
	// Hash pair prefix 1 byte = 0x01
	return hash_tree.New(sha256.New(), block_size, []byte{0x00}, []byte{0x01})
}

func get_pairs_hasher() *hash_pairs.TreePairs {
	// Use sha256 hash for pairs hashing
	// Hash pair prefix 1 byte = 0x01
	// Size of input hash
	haser := sha256.New()
	return hash_pairs.New(haser, []byte{0x01}, haser.Size())
}

func get_root(data []byte) []byte {
	tree := get_tree_hasher()
	tree.Write(data)
	return tree.Sum(nil)
}

func check_root(root []byte, data []byte) bool {
	// Get root hash of data
	data_root := get_root(data)

	// Compare with given root
	return bytes.Equal(root, data_root)
}

func get_list(data []byte, block_size int) []byte {
	tree := get_tree_hasher()

	list := hash_list.New(tree, block_size)

	list.Write(data)

	// Get hash list
	return list.Sum(nil)
}

func get_list_root(list []byte) []byte {
	pairs := get_pairs_hasher()

	// Add hash list to tree
	pairs.Write(list)

	// Get hash
	return pairs.Sum(nil)
}

func check_all_parts(list []byte, data []byte, block_size int, hash_size int) bool {
	result := true

	for i := int(0); i < len(list)/hash_size; i++ {
		part_root := list[i*hash_size : i*hash_size+hash_size]
		data_part := data[i*block_size:]
		if len(data_part) > block_size {
			data_part = data_part[:block_size]
		}
		if !check_root(part_root, data_part) {
			result = false
			fmt.Printf("range bad: %v %v-%v %x != %x\n", i, i*block_size, i*block_size+len(data_part), part_root, get_root(data_part))
		} else {
			fmt.Printf("range good: %v %v-%v %x\n", i, i*block_size, i*block_size+len(data_part), part_root)
		}
	}

	return result
}

func main() {
	data := []byte(strings.Repeat("123456789ABCD", block_size<<2+1))
	data_root := get_root(data)
	fmt.Printf("data root: %x\n", data_root)
	fmt.Printf("check root equal data root: %v\n\n", check_root(data_root, data))

	list1 := get_list(data, block_size)

	fmt.Printf("\n\nlist1: %x\n\n", list1)

	list1_root := get_list_root(list1)
	fmt.Printf("list root: %x\n", list1_root)
	fmt.Printf("list root equal data root: %v\n\n", bytes.Equal(data_root, list1_root))

	fmt.Printf("check all parts by list1: %v\n", check_all_parts(list1, data, block_size, hash_size))

	list2 := get_list(data, block_size<<1)

	fmt.Printf("\n\nlist2: %x\n\n", list2)

	list2_root := get_list_root(list2)
	fmt.Printf("list2 root: %x\n", list2_root)
	fmt.Printf("list2 root equal data root: %v\n\n", bytes.Equal(data_root, list2_root))

	fmt.Printf("check all parts by list2: %v\n", check_all_parts(list2, data, block_size<<1, hash_size))

	fmt.Println("\n\nchange random byte of data\n\n")
	rand.Seed(time.Now().UnixNano())
	byte_index := rand.Intn(len(data))
	data[byte_index] = ^data[byte_index]
	fmt.Printf("check all parts by list1: %v\n", check_all_parts(list1, data, block_size, hash_size))
}
