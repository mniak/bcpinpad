package encoding

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

func N(length, value int) ([]byte, error) {
	return Number(length, value)
}
func Number(length, value int) ([]byte, error) {
	if length < 1 {
		return nil, errors.New("could not format as number. length must be at least 1")
	}
	if value < 0 {
		return nil, errors.New("could not format as number. value cannot be negative")
	}
	max := int(math.Pow10(length)) - 1
	if value > max {
		return nil, errors.New("could not format as number. value is greater than the maximum value")
	}
	f := "%0" + strconv.Itoa(length) + "d"
	str := fmt.Sprintf(f, value)
	return []byte(str), nil
}
