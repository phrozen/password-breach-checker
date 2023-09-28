package format

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	tests := map[uint64]string{
		1:                        "1B",
		2000:                     "1.95KB",
		3000_000:                 "2.86MB",
		4000_000_000:             "3.73GB",
		5000_000_000_000:         "4.55TB",
		6000_000_000_000_000:     "5.33PB",
		7000_000_000_000_000_000: "6.07EB",
	}

	for bytes, want := range tests {
		got := Bytes(bytes)
		assert.Equal(t, want, got, "they should be equal")
	}
}

func BenchmarkBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Bytes(math.MaxUint64)
	}
}
