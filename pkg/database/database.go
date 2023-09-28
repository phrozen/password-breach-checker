package database

import (
	"bytes"
	"encoding/binary"

	"github.com/prometheus/prometheus/tsdb/fileutil"
)

// DB represents a database of password hash/count pairs
// in a sorted memory mapped binary file from:
// https://haveibeenpwned.com/Passwords
// Previously converted to binary format with the cmd/process
// utility, to reduce size and guarantee fixed sizes.
type DB struct {
	data   []byte
	file   *fileutil.MmapFile
	length int
}

// New creates a new memory mapped database based on input binary file.
func New(filename string) (*DB, error) {
	mmap, err := fileutil.OpenMmapFile(filename)
	if err != nil {
		return nil, err
	}
	db := &DB{
		data:   mmap.Bytes(),
		file:   mmap,
		length: len(mmap.Bytes()) / 24,
	}
	return db, nil
}

// Length is the number of hash/count pairs in the file.
// Each pair is 24 bytes in length, 20 bytes (160 bits)
// for the SHA-1 hash and 4 bytes (32 bits) for the count.
func (db *DB) Length() int {
	return db.length
}

// Close the underlying memory mapped file
func (db *DB) Close() error {
	return db.file.Close()
}

// Search performs binary search for the given hash on the
// database and returns the count if found, 0 means the hash
// was not found in the database.
// There is no bounds check, hash should be a []byte (len = 20)
// slice representing a SHA-1 sum of the password to search.
func (db *DB) Search(hash []byte) uint32 {
	// Binary Search: https://go.dev/src/sort/search.go
	low, mid, high := 0, 0, db.length
	for low <= high {
		mid = (low + high) >> 1
		switch bytes.Compare(hash, db.hashAt(mid)) {
		case 1:
			low = mid + 1
		case -1:
			high = mid - 1
		case 0:
			return db.countAt(mid)
		}
	}
	return 0
}

// hashAt returns the hash value at position i
// No mutex used as database is READ ONLY.
// No bounds check, internal only
func (db *DB) hashAt(i int) []byte {
	// Each hash/count pair is 24 bytes
	idx := i * 24
	// each hash is 20 bytes (SHA1 160 bits)
	return db.data[idx : idx+20]
}

// countAt returns the count value at position i
// No bounds check, internal only
func (db *DB) countAt(i int) uint32 {
	// Each hash/count pair is 24 bytes
	idx := i*24 + 20
	// count is 4 bytes with an offset of 20
	return binary.BigEndian.Uint32(db.data[idx : idx+4])
}
