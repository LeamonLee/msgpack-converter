package msgpackconverter

import (
	"errors"
)

// Encodes bytes into MessagePack binary format.
func (conv *MsgPackConverter) encodeBytes(value []byte) ([]byte, error) {
	length := len(value)
	switch {
	case length < (1 << 8):
		// fmt.Println("Encode type bin8")
		return append([]byte{bin8Type, byte(length)}, value...), nil
	case length < (1 << 16):
		// fmt.Println("Encode type bin16")
		return append([]byte{bin16Type, byte(length >> 8), byte(length)}, value...), nil
	case length < (1 << 32):
		// fmt.Println("Encode type bin32")
		return append([]byte{bin32Type, byte(length >> 24), byte(length >> 16), byte(length >> 8), byte(length)}, value...), nil
	default:
		return nil, errors.New("value too large to encode in MessagePack format")
	}
}

// Decodes a bin 8, 16, or 32 byte slice from MessagePack format.
func (conv *MsgPackConverter) decodeBytes(msgpackData []byte) ([]byte, error) {
	if len(msgpackData) < 2 {
		return nil, errors.New("insufficient data for decoding")
	}

	var binLength int
	lengthByte := msgpackData[0]

	switch lengthByte {
	case bin8Type:
		binLength = int(msgpackData[1])
		if len(msgpackData) < binLength+2 {
			return nil, errors.New("data length mismatch for bin8 decoding")
		}
		return msgpackData[2 : 2+binLength], nil

	case bin16Type:
		if len(msgpackData) < 3 {
			return nil, errors.New("insufficient data for bin16 decoding")
		}
		binLength = int(msgpackData[1])<<8 | int(msgpackData[2])
		if len(msgpackData) < binLength+3 {
			return nil, errors.New("data length mismatch for bin16 decoding")
		}
		return msgpackData[3 : 3+binLength], nil

	case bin32Type:
		if len(msgpackData) < 5 {
			return nil, errors.New("insufficient data for bin32 decoding")
		}
		binLength = int(msgpackData[1])<<24 | int(msgpackData[2])<<16 | int(msgpackData[3])<<8 | int(msgpackData[4])
		if len(msgpackData) < binLength+5 {
			return nil, errors.New("data length mismatch for bin32 decoding")
		}
		return msgpackData[5 : 5+binLength], nil

	default:
		return nil, errors.New("unknown data type for decoding")
	}
}
