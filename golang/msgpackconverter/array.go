package msgpackconverter

import "errors"

func (conv *MsgPackConverter) encodeList(value []interface{}) ([]byte, error) {
	var msgpackData []byte
	listLength := len(value)

	if listLength <= 15 {
		msgpackData = append(msgpackData, byte(fixArrayMinType|listLength))
	} else if listLength <= 0xffff {
		msgpackData = append(msgpackData, array16Type)
		msgpackData = append(msgpackData, byte(listLength>>8), byte(listLength&0xff))
	} else {
		msgpackData = append(msgpackData, array32Type)
		msgpackData = append(msgpackData, byte(listLength>>24), byte(listLength>>16), byte(listLength>>8), byte(listLength&0xff))
	}

	for _, item := range value {
		encodedItem, err := conv.EncodeJSONToMsgPack(item)
		if err != nil {
			return nil, err
		}
		msgpackData = append(msgpackData, encodedItem...)
	}

	return msgpackData, nil
}

func (conv *MsgPackConverter) decodeList(data []byte, firstByte byte) ([]interface{}, error) {
	var size int
	if firstByte == array16Type {
		if len(data) < 2 {
			return nil, errors.New("not enough data to decode array16Type length")
		}
		size = int(data[0])<<8 + int(data[1])
		data = data[2:]
	} else if firstByte == array32Type {
		if len(data) < 4 {
			return nil, errors.New("not enough data to decode array32Type length")
		}
		size = int(data[0])<<24 + int(data[1])<<16 + int(data[2])<<8 + int(data[3])
		data = data[4:]
	} else {
		return nil, errors.New("invalid first byte for decoding list")
	}

	var decodedList []interface{}
	for i := 0; i < size; i++ {
		item, err := conv.DecodeMsgPackToJSON(data)
		if err != nil {
			return nil, err
		}
		item2, _ := conv.EncodeJSONToMsgPack(item)
		data = data[len(item2):]
		decodedList = append(decodedList, item)
	}

	return decodedList, nil

}

func (conv *MsgPackConverter) decodeFixArray(data []byte, firstByte byte) ([]interface{}, error) {
	listLength := int(firstByte & 0x0f)
	decodedList := make([]interface{}, 0, listLength)

	for i := 0; i < listLength; i++ {
		item, err := conv.DecodeMsgPackToJSON(data)
		if err != nil {
			return nil, err
		}
		item2, _ := conv.EncodeJSONToMsgPack(item)
		data = data[len(item2):]
		decodedList = append(decodedList, item)
	}

	return decodedList, nil
}
