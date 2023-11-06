package msgpackconverter

import (
	"errors"
)

func (conv *MsgPackConverter) encodeString(value string) ([]byte, error) {
	var msgpackData []byte
	strLength := len(value)

	if strLength <= 31 {
		msgpackData = append(msgpackData, byte(fixStrMinType|strLength))
	} else if strLength <= (1<<8)-1 {
		msgpackData = append(msgpackData, str8Type, byte(strLength))
	} else if strLength <= (1<<16)-1 {
		msgpackData = append(msgpackData, str16Type)
		msgpackData = append(msgpackData, byte(strLength>>8), byte(strLength))
	} else {
		msgpackData = append(msgpackData, str32Type)
		msgpackData = append(msgpackData, byte(strLength>>24), byte(strLength>>16), byte(strLength>>8), byte(strLength))
	}

	msgpackData = append(msgpackData, []byte(value)...)

	return msgpackData, nil
}

// func (conv *MsgPackConverter) decodeString(data []byte) (string, error) {
// 	var length int
// 	switch data[0] {
// 	case str8Type:
// 		length = int(data[1])
// 		return string(data[2 : 2+length]), nil
// 	case str16Type:
// 		length = int(binary.BigEndian.Uint16(data[1:3]))
// 		return string(data[3 : 3+length]), nil
// 	case str32Type:
// 		length = int(binary.BigEndian.Uint32(data[1:5]))
// 		return string(data[5 : 5+length]), nil
// 	default:
// 		return "", errors.New("unsupported string type")
// 	}
// }

func (conv *MsgPackConverter) decodeString(data []byte) (string, error) {
	var length int
	switch data[0] {
	case str8Type:
		length = int(data[1])
	case str16Type:
		length = (int(data[1]) << 8) | int(data[2])
	case str32Type:
		length = (int(data[1]) << 24) | (int(data[2]) << 16) | (int(data[3]) << 8) | int(data[4])
	default:
		return "", errors.New("unsupported string type")
	}

	return string(data[1+bytesNeededForLength(length) : 1+bytesNeededForLength(length)+length]), nil
}

// Helper function to determine the number of bytes needed to encode the length
func bytesNeededForLength(length int) int {
	switch {
	case length <= (1<<8)-1:
		return 1 // str8Type
	case length <= (1<<16)-1:
		return 2 // str16Type
	case length <= (1<<32)-1:
		return 4 // str32Type
	default:
		return 0
	}
}
