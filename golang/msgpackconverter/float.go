package msgpackconverter

import (
	"errors"
	"math"
)

func (conv *MsgPackConverter) encodeFloat(value interface{}) ([]byte, error) {
	var buf []byte

	switch v := value.(type) {
	case float32:
		// fmt.Println("encodeFloat float32")
		buf = append(buf, float32Type)
		bits := math.Float32bits(v)
		buf = append(buf,
			byte(bits>>24),
			byte(bits>>16),
			byte(bits>>8),
			byte(bits),
		)
	case float64:
		// fmt.Println("encodeFloat float64")
		buf = append(buf, float64Type)
		bits := math.Float64bits(v)
		buf = append(buf,
			byte(bits>>56),
			byte(bits>>48),
			byte(bits>>40),
			byte(bits>>32),
			byte(bits>>24),
			byte(bits>>16),
			byte(bits>>8),
			byte(bits),
		)
	}

	return buf, nil
}

func (conv *MsgPackConverter) decodeFloat32(data []byte) (interface{}, error) {
	if len(data) < 5 {
		return nil, errors.New("not enough bytes to decode float32")
	}
	bits := uint32(data[1])<<24 | uint32(data[2])<<16 | uint32(data[3])<<8 | uint32(data[4])
	return math.Float32frombits(bits), nil
}

func (conv *MsgPackConverter) decodeFloat64(data []byte) (interface{}, error) {
	if len(data) < 9 {
		return nil, errors.New("not enough bytes to decode float64")
	}
	bits := uint64(data[1])<<56 | uint64(data[2])<<48 | uint64(data[3])<<40 | uint64(data[4])<<32 |
		uint64(data[5])<<24 | uint64(data[6])<<16 | uint64(data[7])<<8 | uint64(data[8])
	return math.Float64frombits(bits), nil
}
