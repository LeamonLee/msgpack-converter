package msgpackconverter

import (
	"errors"
)

type MsgPackConverter struct{}

func (conv *MsgPackConverter) EncodeJSONToMsgPack(data interface{}) ([]byte, error) {
	switch v := data.(type) {
	case nil:
		return []byte{nilType}, nil
	case bool:
		if v {
			return []byte{trueType}, nil
		} else {
			return []byte{falseType}, nil
		}
	case int:
		// fmt.Println("Encode int generic type")
		return conv.encodeFixInteger(v)
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		// fmt.Println("Encode int/unit type")
		return conv.encodeInteger(v)
	case []byte:
		// fmt.Println("Encode []byte type")
		return conv.encodeBytes(v)
	case float32, float64:
		// fmt.Println("Encode []float type")
		return conv.encodeFloat(v)
	case string:
		// fmt.Println("Encode string type")
		return conv.encodeString(v)
	case []interface{}:
		// fmt.Println("Encode array type")
		return conv.encodeList(v)
	case map[string]interface{}:
		return conv.encodeDict(v)
	default:
		// fmt.Println("data type: ", v, "data: ", data)
		return nil, errors.New("unsupported type")
	}
}

func (conv *MsgPackConverter) DecodeMsgPackToJSON(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	switch data[0] {
	case nilType:
		return nil, nil
	case falseType:
		return false, nil
	case trueType:
		return true, nil
	case bin8Type, bin16Type, bin32Type:
		return conv.decodeBytes(data)
	case uint8Type, uint16Type, uint32Type, uint64Type, int8Type, int16Type, int32Type, int64Type:
		// fmt.Println("Decode uint/int type")
		return conv.decodeInteger(data)
	case float32Type:
		// fmt.Println("Decode float32 type")
		return conv.decodeFloat32(data)
	case float64Type:
		// fmt.Println("Decode float64 type")
		return conv.decodeFloat64(data)
	case str8Type, str16Type, str32Type:
		// fmt.Println("Decode string type")
		return conv.decodeString(data)
	case array16Type, array32Type:
		// fmt.Println("Decode array16Type, array32Type type")
		return conv.decodeList(data[1:], data[0])
	case map16Type, map32Type:
		// fmt.Println("Decode dict type")
		return conv.decodeDict(data[1:], data[0])
	default:
		if (data[0] >= fixPositiveIntMinType && data[0] <= fixPositiveIntMaxType) || (data[0] >= fixNegativeIntMinType && data[0] <= fixNegativeIntMaxType) { // fixint
			return conv.decodeFixInteger(data[0])
		} else if fixStrMinType <= data[0] && data[0] <= fixStrMaxType { // fixstr
			strLength := int(data[0] & 0x1f)
			return string(data[1 : strLength+1]), nil
		} else if fixArrayMinType <= data[0] && data[0] <= fixArrayMaxType { // fixarr
			return conv.decodeFixArray(data[1:], data[0])
		} else if fixMapMinType <= data[0] && data[0] <= fixMapMaxType { // fixmap
			// fmt.Println("Decode fixmap type")
			return conv.decodeFixMap(data[1:], data[0])
		}
		return nil, errors.New("Decode unsupported type")
	}
}
