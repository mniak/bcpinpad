package highlevel

import (
	"bytes"
	"fmt"

	"github.com/mniak/bcpinpad/encoding"
	"github.com/mniak/bcpinpad/utils"
)

const (
	General  GetInfoType = 0
	Amex     GetInfoType = 1
	Redecard GetInfoType = 2
	Cielo    GetInfoType = 3
)

type GetInfoType int
type GetInfoParams struct {
	Type GetInfoType
}

func (p GetInfoParams) Bytes() ([]byte, error) {
	t, err := encoding.N(2, int(p.Type))
	if err != nil {
		return nil, err
	}
	return utils.NewBytesBuilder().
		AddString("GIN").
		AddBytes(t).
		Bytes(), nil
}

type GetInfoResult struct {
	Name string
}

func parseGetInfo00Result(b []byte) (result GetInfoResult, err error) {
	buffer := bytes.NewBuffer(b)
	name := make([]byte, 20)
	n, err := buffer.Read(name)
	if err != nil {
		return
	}
	if n < 20 {
		return result, fmt.Errorf("expecting 20 bytes but got only %d", n)
	}
	result.Name = string(name)
	return
}

func (pp Pinpad) GetInfo(p GetInfoParams) (GetInfoResult, error) {
	bytes, err := p.Bytes()
	if err != nil {
		return GetInfoResult{}, err
	}
	resultBytes, err := pp.pp.Send(bytes)
	if err != nil {
		return GetInfoResult{}, err
	}
	result, err := parseGetInfo00Result(resultBytes)
	return result, err
}
