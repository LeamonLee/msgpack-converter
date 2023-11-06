package msgpackconverter

func (conv *MsgPackConverter) encodeDict(value map[string]interface{}) ([]byte, error) {
	var msgpackData []byte

	dictLength := len(value)
	// fmt.Printf("__encode_dict dictLength:%d\n", dictLength)
	switch {
	case dictLength <= 15:
		// fmt.Println("__encode_dict fixmap")
		msgpackData = append(msgpackData, byte(fixMapMinType|dictLength))
	case dictLength <= 0xffff:
		// fmt.Println("__encode_dict map16Type")
		msgpackData = append(msgpackData, map16Type)
		msgpackData = append(msgpackData, byte(dictLength>>8), byte(dictLength))
	default:
		// fmt.Println("__encode_dict map32Type")
		msgpackData = append(msgpackData, map32Type)
		msgpackData = append(msgpackData, byte(dictLength>>24), byte(dictLength>>16), byte(dictLength>>8), byte(dictLength))
	}

	for key, val := range value {
		encodedKey, err := conv.EncodeJSONToMsgPack(key)
		if err != nil {
			return nil, err
		}
		msgpackData = append(msgpackData, encodedKey...)

		encodedValue, err := conv.EncodeJSONToMsgPack(val)
		if err != nil {
			return nil, err
		}
		msgpackData = append(msgpackData, encodedValue...)
	}

	return msgpackData, nil
}

func (conv *MsgPackConverter) decodeFixMap(msgpackData []byte, firstByte byte) (map[string]interface{}, error) {
	dictLength := int(firstByte & 0x0f)
	decodedDict := make(map[string]interface{})
	var err error

	// fmt.Println("decodeFixMap dictLength: ", dictLength)
	for i := 0; i < dictLength; i++ {
		var key interface{}
		key, err = conv.DecodeMsgPackToJSON(msgpackData)
		if err != nil {
			return nil, err
		}
		keyBytes, _ := conv.EncodeJSONToMsgPack(key)
		msgpackData = msgpackData[len(keyBytes):]

		var value interface{}
		value, err = conv.DecodeMsgPackToJSON(msgpackData)
		if err != nil {
			return nil, err
		}
		decodedDict[key.(string)] = value
		valueBytes, _ := conv.EncodeJSONToMsgPack(value)
		msgpackData = msgpackData[len(valueBytes):]
	}

	return decodedDict, nil
}

func (conv *MsgPackConverter) decodeDict(msgpackData []byte, firstByte byte) (map[string]interface{}, error) {
	sizes := []int{2, 4}
	size := sizes[firstByte-map16Type]

	var dictLength int
	if size == 2 {
		dictLength = int(msgpackData[0])<<8 | int(msgpackData[1])
	} else {
		dictLength = int(msgpackData[0])<<24 | int(msgpackData[1])<<16 | int(msgpackData[2])<<8 | int(msgpackData[3])
	}

	msgpackData = msgpackData[size:]
	decodedDict := make(map[string]interface{})
	var err error

	for i := 0; i < dictLength; i++ {
		var key interface{}
		key, err = conv.DecodeMsgPackToJSON(msgpackData)
		if err != nil {
			return nil, err
		}
		keyBytes, _ := conv.EncodeJSONToMsgPack(key)
		msgpackData = msgpackData[len(keyBytes):]

		var value interface{}
		value, err = conv.DecodeMsgPackToJSON(msgpackData)
		if err != nil {
			return nil, err
		}
		decodedDict[key.(string)] = value
		valueBytes, _ := conv.EncodeJSONToMsgPack(value)
		msgpackData = msgpackData[len(valueBytes):]
	}

	return decodedDict, nil
}
