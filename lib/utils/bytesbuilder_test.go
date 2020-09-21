package utils

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().Unix())
}
func TestBytesBuilder_NewHasNoBytes(t *testing.T) {
	bytes := NewBytesBuilder().Bytes()
	assert.Empty(t, bytes)
}

func TestBytesBuilder_AddByte(t *testing.T) {
	b1 := byte(rand.Intn(256))
	b2 := byte(rand.Intn(256))
	b3 := byte(rand.Intn(256))
	bytes := NewBytesBuilder().
		AddByte(b1).
		AddByte(b2, b3).
		Bytes()
	assert.EqualValues(t, bytes, []byte{b1, b2, b3})
}

func TestBytesBuilder_AddBytes(t *testing.T) {
	b1 := byte(rand.Intn(256))
	b2 := byte(rand.Intn(256))
	b3 := byte(rand.Intn(256))
	b4 := byte(rand.Intn(256))
	bytes := NewBytesBuilder().
		AddBytes([]byte{b1, b2}).
		AddBytes([]byte{b3, b4}).
		Bytes()
	assert.EqualValues(t, bytes, []byte{b1, b2, b3, b4})
}

func randchar() byte {
	return byte(rand.Intn(0x7e-0x20+1) + 0x20)
}
func TestBytesBuilder_AddString(t *testing.T) {
	b1 := randchar()
	b2 := randchar()
	b3 := randchar()
	b4 := randchar()
	b5 := randchar()
	b6 := randchar()
	bytes := NewBytesBuilder().
		AddString(string([]byte{b1, b2})).
		AddString(string([]byte{b3, b4}), string([]byte{b5, b6})).
		Bytes()
	assert.EqualValues(t, bytes, []byte{b1, b2, b3, b4, b5, b6})
}

func TestBytesBuilder_AddString_CharsLowerThan0x20AreIgnored(t *testing.T) {
	b1 := byte(rand.Intn(0x20))
	b2 := randchar()
	b3 := byte(rand.Intn(0x20))
	b4 := randchar()
	b5 := byte(rand.Intn(0x20))
	b6 := randchar()
	bytes := NewBytesBuilder().
		AddString(string([]byte{b1, b2})).
		AddString(string([]byte{b3, b4}), string([]byte{b5, b6})).
		Bytes()
	assert.EqualValues(t, bytes, []byte{b2, b4, b6})
}

func TestBytesBuilder_AddString_CharsGreaterThan0x7eAreIgnored(t *testing.T) {
	b1 := randchar()
	b2 := byte(rand.Intn(256-0x7e-1) + 0x7e + 1)
	b3 := randchar()
	b4 := byte(rand.Intn(256-0x7e-1) + 0x7e + 1)
	b5 := randchar()
	b6 := byte(rand.Intn(256-0x7e-1) + 0x7e + 1)
	bytes := NewBytesBuilder().
		AddString(string([]byte{b1, b2})).
		AddString(string([]byte{b3, b4}), string([]byte{b5, b6})).
		Bytes()
	assert.EqualValues(t, bytes, []byte{b1, b3, b5})
}
