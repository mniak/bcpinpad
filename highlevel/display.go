package highlevel

import (
	"github.com/mniak/bcpinpad/encoding"
	"github.com/mniak/bcpinpad/utils"
)

type DisplayParams struct {
	Message string
}

func (p DisplayParams) Bytes() ([]byte, error) {
	msg, err := encoding.Alpha(32, p.Message)
	if err != nil {
		return nil, err
	}
	return utils.NewBytesBuilder().
		AddString("DSP").
		AddBytes(msg).
		Bytes(), nil
}

type DisplayResult struct {
}

func parseDisplayResult(b []byte) (result DisplayResult, err error) {
	return
}

func (pp Pinpad) Display(p DisplayParams) (DisplayResult, error) {
	bytes, err := p.Bytes()
	if err != nil {
		return DisplayResult{}, err
	}
	resultBytes, err := pp.pp.Send(bytes)
	if err != nil {
		return DisplayResult{}, err
	}
	result, err := parseDisplayResult(resultBytes)
	return result, err
}
