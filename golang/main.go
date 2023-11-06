package main

import (
	"fmt"

	"github.com/LeamonLee/msgpack-converter/msgpackconverter"
)

func main() {
	// conv := &MsgPackConverter{}

	// encodedData, _ := conv.EncodeJSONToMsgPack(12345)
	// fmt.Printf("Encoded data: %v\n", encodedData)

	// decodedData, err := conv.DecodeMsgPackToJSON(encodedData)
	// if err != nil {
	// 	fmt.Println("An error occurred:", err)
	// } else {
	// 	fmt.Printf("Decoded data: %v\n", decodedData)
	// }

	// // Test float encoding and decoding
	// floatEncoded, err := conv.EncodeJSONToMsgPack(3.1445384234902380)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(floatEncoded)
	// floatDecoded, err := conv.DecodeMsgPackToJSON(floatEncoded)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(floatDecoded)

	// // Test string encoding and decoding
	// str := strings.Repeat("Hello,", 10)
	// stringEncoded, err := conv.EncodeJSONToMsgPack(str)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// stringDecoded, err := conv.DecodeMsgPackToJSON(stringEncoded)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(stringDecoded)

	// Example usage:
	conv := &msgpackconverter.MsgPackConverter{}
	jsonArray := []interface{}{"example322", 123, uint16(1000), uint32(100000), float32(3.14), 3.141592653589763, true, false, nil, []interface{}{"nested", "array", 100, 99}}
	encoded, err := conv.EncodeJSONToMsgPack(jsonArray)
	if err != nil {
		fmt.Println("Encode error:", err)
		return
	}
	fmt.Printf("Encoded MessagePack: %v\n", encoded)

	decoded, err := conv.DecodeMsgPackToJSON(encoded)
	if err != nil {
		fmt.Println("Decode error:", err)
		return
	}
	fmt.Printf("Decoded JSON: %v\n", decoded)

	// Encode Json example
	jsonData := map[string]interface{}{"k1": "value1", "k2": 100, "k3": int16(10000), "k4": int32(100000), "k5": float32(3.14), "k6": float64(3.141592653589)}
	encodedData, err := conv.EncodeJSONToMsgPack(jsonData)
	if err != nil {
		fmt.Println("An error occurred while encoding:", err)
		return
	}
	fmt.Println("Encoded Data:", encodedData)

	// Decode Json example
	decodedData, err := conv.DecodeMsgPackToJSON(encodedData)
	if err != nil {
		fmt.Println("An error occurred while decoding:", err)
		return
	}
	fmt.Println("Decoded Data:", decodedData)

	// // Create a map with string keys and empty interface type values
	// mixedTypeMap := make(map[string]interface{})

	// // Generate more than 100 keys with different types of values
	// for i := 0; i < 100; i++ {
	// 	key := fmt.Sprintf("key%d", i)
	// 	// Assign a value to the map entry based on the modulus of i
	// 	switch i % 13 { // Updated modulus to 13 to include new cases
	// 	case 0:
	// 		mixedTypeMap[key] = nil
	// 	case 1:
	// 		mixedTypeMap[key] = true
	// 	case 2:
	// 		mixedTypeMap[key] = "a string value"
	// 	case 3:
	// 		mixedTypeMap[key] = uint8(rand.Intn(256))
	// 	case 4:
	// 		mixedTypeMap[key] = uint16(rand.Intn(65536))
	// 	case 5:
	// 		mixedTypeMap[key] = uint32(rand.Uint32())
	// 	case 6:
	// 		mixedTypeMap[key] = int8(rand.Intn(256) - 128)
	// 	case 7:
	// 		mixedTypeMap[key] = int16(rand.Intn(65536) - 32768)
	// 	case 8:
	// 		mixedTypeMap[key] = int32(rand.Int31())
	// 	case 9:
	// 		mixedTypeMap[key] = float32(rand.Float32())
	// 	case 10:
	// 		mixedTypeMap[key] = rand.Float64()
	// 	case 11: // Nested array
	// 		mixedTypeMap[key] = []interface{}{rand.Intn(100), "nested", false}
	// 	case 12: // Nested map
	// 		nestedMap := make(map[string]interface{})
	// 		nestedMap["nestedKey1"] = rand.Intn(100)
	// 		nestedMap["nestedKey2"] = "nested value"
	// 		nestedMap["nestedKey3"] = []interface{}{"array", "inside a nested map", 42}
	// 		mixedTypeMap[key] = nestedMap
	// 	}
	// }

	// // Encode Json example
	// encodedData, err = conv.EncodeJSONToMsgPack(mixedTypeMap)
	// if err != nil {
	// 	fmt.Println("An error occurred while encoding:", err)
	// 	return
	// }
	// fmt.Println("Encoded Data:", encodedData)

	// // Decode Json example
	// decodedData, err = conv.DecodeMsgPackToJSON(encodedData)
	// if err != nil {
	// 	fmt.Println("An error occurred while decoding:", err)
	// 	return
	// }
	// fmt.Println("Decoded Data:", decodedData)
}
