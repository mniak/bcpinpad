package encoding

import (
	"fmt"
)

func A(length int, value string) ([]byte, error) {
	return Alpha(length, value)
}
func Alpha(length int, value string) ([]byte, error) {
	bytes := []byte(value)
	if length != len(bytes) {
		return nil, fmt.Errorf("invalid value. it was expected length %d but got %d", length, len(bytes))
	}

	return bytes, nil
}
