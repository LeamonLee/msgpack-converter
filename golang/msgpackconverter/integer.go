package msgpackconverter

import (
	"fmt"
)

func putUint16(b []byte, v uint16) {
	b[0] = byte(v >> 8)
	b[1] = byte(v)
}

func putUint32(b []byte, v uint32) {
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)
}

func putUint64(b []byte, v uint64) {
	b[0] = byte(v >> 56)
	b[1] = byte(v >> 48)
	b[2] = byte(v >> 40)
	b[3] = byte(v >> 32)
	b[4] = byte(v >> 24)
	b[5] = byte(v >> 16)
	b[6] = byte(v >> 8)
	b[7] = byte(v)
}

func putInt16(b []byte, v int16) {
	putUint16(b, uint16(v))
}

func putInt32(b []byte, v int32) {
	putUint32(b, uint32(v))
}

func putInt64(b []byte, v int64) {
	putUint64(b, uint64(v))
}

func getUint16(b []byte) uint16 {
	return uint16(b[0])<<8 | uint16(b[1])
}

func getUint32(b []byte) uint32 {
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

func getUint64(b []byte) uint64 {
	return uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])
}

func getInt16(b []byte) int16 {
	return int16(getUint16(b))
}

func getInt32(b []byte) int32 {
	return int32(getUint32(b))
}

func getInt64(b []byte) int64 {
	return int64(getUint64(b))
}

func (conv *MsgPackConverter) encodeFixInteger(value int) ([]byte, error) {
	if value < -32 || value > 127 {
		return nil, fmt.Errorf("value out of range for fixint encoding")
	}

	msgpackData := make([]byte, 1)
	if value >= 0 && value <= 127 { // positive fixint
		msgpackData[0] = byte(value)
	} else if value >= -32 && value < 0 { // negative fixint
		msgpackData[0] = byte(value & 0xFF)
	}

	return msgpackData, nil
}

func (conv *MsgPackConverter) encodeInteger(value interface{}) ([]byte, error) {
	// fixIntValue := value.(int)
	// if fixIntValue >= -32 && fixIntValue <= 127 {
	// 	return conv.encodeFixInteger(fixIntValue)
	// }
	var buf []byte
	switch v := value.(type) {
	case uint8:
		// fmt.Println("Encode type uint8")
		buf = []byte{uint8Type, v}
	case uint16:
		// fmt.Println("Encode type uint16")
		buf = append([]byte{uint16Type}, make([]byte, 2)...)
		putUint16(buf[1:], v)
	case uint32:
		// fmt.Println("Encode type uint32")
		buf = append([]byte{uint32Type}, make([]byte, 4)...)
		putUint32(buf[1:], v)
	case uint64:
		// fmt.Println("Encode type uint64")
		buf = append([]byte{uint64Type}, make([]byte, 8)...)
		putUint64(buf[1:], v)
	case int8:
		// fmt.Println("Encode type int8")
		buf = []byte{int8Type, byte(v)}
	case int16:
		// fmt.Println("Encode type int16")
		buf = append([]byte{int16Type}, make([]byte, 2)...)
		putInt16(buf[1:], v)
	case int32:
		// fmt.Println("Encode type int32")
		buf = append([]byte{int32Type}, make([]byte, 4)...)
		putInt32(buf[1:], v)
	case int64:
		// fmt.Println("Encode type int64")
		buf = append([]byte{int64Type}, make([]byte, 8)...)
		putInt64(buf[1:], v)
	default:
		return nil, fmt.Errorf("encodeInteger unsupported type %T", v)
	}
	return buf, nil
}

func (conv *MsgPackConverter) decodeFixInteger(data byte) (int, error) {
	// if len(data) != 1 {
	// 	return 0, fmt.Errorf("invalid data length for fixint decoding")
	// }

	// firstByte := data[0]

	if data&fixMapMinType == 0 { // positive fixint
		return int(data), nil
	} else if data >= 0xE0 { // negative fixint
		// Convert to a signed 8-bit integer to ensure correct negative value
		return int(int8(data)), nil
	}

	return 0, fmt.Errorf("invalid data for fixint decoding")
}

func (conv *MsgPackConverter) decodeInteger(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	switch data[0] {
	case uint8Type:
		if len(data) < 2 {
			return nil, fmt.Errorf("not enough data for uint8")
		}
		return uint8(data[1]), nil
	case uint16Type:
		if len(data) < 3 {
			return nil, fmt.Errorf("not enough data for uint16")
		}
		return getUint16(data[1:3]), nil
	case uint32Type:
		if len(data) < 5 {
			return nil, fmt.Errorf("not enough data for uint32")
		}
		return getUint32(data[1:5]), nil
	case uint64Type:
		if len(data) < 9 {
			return nil, fmt.Errorf("not enough data for uint64")
		}
		return getUint64(data[1:9]), nil
	case int8Type:
		if len(data) < 2 {
			return nil, fmt.Errorf("not enough data for int8")
		}
		return int8(data[1]), nil
	case int16Type:
		if len(data) < 3 {
			return nil, fmt.Errorf("not enough data for int16")
		}
		return getInt16(data[1:3]), nil
	case int32Type:
		if len(data) < 5 {
			return nil, fmt.Errorf("not enough data for int32")
		}
		return getInt32(data[1:5]), nil
	case int64Type:
		if len(data) < 9 {
			return nil, fmt.Errorf("not enough data for int64")
		}
		return getInt64(data[1:9]), nil
	default:
		return nil, fmt.Errorf("unsupported type prefix: %x", data[0])
	}
}
