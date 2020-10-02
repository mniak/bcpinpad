package utils

// BytesBuilder
type BytesBuilder struct {
	bytes []byte
}

// NewBytesBuilder creates a Builder
func NewBytesBuilder() *BytesBuilder {
	return &BytesBuilder{}
}

// AddByte adds one or many bytes to the Builder
func (bb *BytesBuilder) AddByte(bytes ...byte) *BytesBuilder {
	return bb.AddBytes(bytes)
}

// AddBytes adds many bytes to the Builder
func (bb *BytesBuilder) AddBytes(bytes []byte) *BytesBuilder {
	bb.bytes = append(bb.bytes, bytes...)
	return bb
}

// AddString adds one or many strings to the Builder
func (bb *BytesBuilder) AddString(strings ...string) *BytesBuilder {
	for _, s := range strings {
		bytes := ascii(s)
		bb = bb.AddBytes(bytes)
	}
	return bb
}

func ascii(s string) []byte {
	result := make([]byte, 0, len(s))
	for _, r := range s {
		if r >= 0x20 && r <= 0x7e {
			result = append(result, byte(r))
		}
	}
	return result
}

// Bytes returns the bytes of the builder
func (bb *BytesBuilder) Bytes() []byte {
	return bb.bytes
}
