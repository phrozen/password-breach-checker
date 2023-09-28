package database

import (
	"bytes"
	"crypto/rand"
	random "math/rand"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DatabaseTestSuite struct {
	suite.Suite
	db *DB
}

func (ts *DatabaseTestSuite) SetupSuite() {
	db, err := New("../../data/pwned-passwords-sha1-ordered-by-hash-v8.bin")
	ts.Equal(err, nil)
	ts.db = db
}

func (ts *DatabaseTestSuite) TearDownSuite() {
	err := ts.db.Close()
	ts.Equal(err, nil)
}

func (ts *DatabaseTestSuite) TestOpenFileError() {
	_, err := New("non-existent-file.bin")
	ts.NotNil(err)
}

func (ts *DatabaseTestSuite) TestSearchFound() {
	for i := 0; i < 1000; i++ {
		idx := random.Intn(ts.db.Length())
		have := ts.db.countAt(idx)
		hash := ts.db.hashAt(idx)
		want := ts.db.Search(hash)
		ts.Equal(want, have, "counts should match")
	}
}

func (ts *DatabaseTestSuite) TestSearchBounds() {
	bounds := []int{0, ts.db.Length() - 1}
	for _, i := range bounds {
		hash := ts.db.hashAt(i)
		want := ts.db.countAt(i)
		have := ts.db.Search(hash)
		ts.Equal(want, have, "counts should match")
	}
}

func (ts *DatabaseTestSuite) TestSearchNotFound() {
	b := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}
	want := uint32(0)
	for i := range b {
		hash := bytes.Repeat(b[i:i+1], 20)
		have := ts.db.Search(hash)
		ts.Equal(want, have, "count should be 0")
	}
}

func (ts *DatabaseTestSuite) TestSearchNotFoundRandom() {
	hash := make([]byte, 20)
	want := uint32(0)
	for i := 0; i < 1000; i++ {
		rand.Read(hash)
		have := ts.db.Search(hash)
		ts.Equal(want, have, "counts should be 0")
	}
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func BenchmarkSearch(b *testing.B) {
	db, err := New("../../data/pwned-passwords-sha1-ordered-by-hash-v8.bin")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	hash := make([]byte, 20)
	rand.Read(hash)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Search(hash)
	}
}

func BenchmarkSearchRandom(b *testing.B) {
	db, err := New("../../data/pwned-passwords-sha1-ordered-by-hash-v8.bin")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	hashes := make([]byte, b.N+20)
	rand.Read(hashes)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Search(hashes[i : i+20])
	}
}

func BenchmarkSearchCompare(b *testing.B) {
	db, err := New("../../data/pwned-passwords-sha1-ordered-by-hash-v8.bin")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		have := db.countAt(i)
		hash := db.hashAt(i)
		want := db.Search(hash)
		if have != want {
			b.Fatal("counts do not match")
		}
	}
}
