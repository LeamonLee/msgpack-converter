package msgpackconverter

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"
	// "github.com/LeamonLee/msgpack-converter/msgpackconverter"
)

func TestMsgPackConverter(t *testing.T) {
	// conv := &msgpackconverter.MsgPackConverter{}
	conv := &MsgPackConverter{}

	t.Run("positive fixint", func(t *testing.T) {
		input := 127
		encoded, err := conv.EncodeJSONToMsgPack(input)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack failed: %v", err)
		}

		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if decoded != input {
			t.Errorf("Expected 127, got %v", decoded)
		}
	})

	t.Run("negative fixint", func(t *testing.T) {
		input := -1
		encoded, err := conv.EncodeJSONToMsgPack(input)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack failed: %v", err)
		}

		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if decoded != input {
			t.Errorf("Expected -1, got %v", decoded)
		}
	})

	t.Run("fixmap", func(t *testing.T) {
		input := map[string]interface{}{"a": 1}
		encoded, err := conv.EncodeJSONToMsgPack(input)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack failed: %v", err)
		}
		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("DecodeMsgPackToJSON failed: %v", err)
		}
		// Compare JSON strings to avoid issues with map ordering
		inputJSON, _ := json.Marshal(input)
		decodedJSON, _ := json.Marshal(decoded)
		if string(inputJSON) != string(decodedJSON) {
			t.Errorf("Expected %v, got %v", string(inputJSON), string(decodedJSON))
		}
	})

	t.Run("fixarray", func(t *testing.T) {
		input := []interface{}{1, 2, 3}
		encoded, err := conv.EncodeJSONToMsgPack(input)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack failed: %v", err)
		}
		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("DecodeMsgPackToJSON failed: %v", err)
		}
		decodedArray, ok := decoded.([]interface{})
		if !ok {
			t.Fatalf("Decoded is not an array: %v", decoded)
		}
		if len(decodedArray) != len(input) {
			t.Fatalf("Arrays are not the same length: got %v want %v", len(decodedArray), len(input))
		}
		for i, v := range input {
			if decodedArray[i] != v {
				t.Errorf("At index %d, got %v want %v", i, decodedArray[i], v)
			}
		}
	})

	t.Run("fixstr", func(t *testing.T) {
		input := "hello"
		encoded, err := conv.EncodeJSONToMsgPack(input)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack failed: %v", err)
		}
		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("DecodeMsgPackToJSON failed: %v", err)
		}
		if decoded != input {
			t.Errorf("Expected %v, got %v", input, decoded)
		}
	})

	t.Run("nil", func(t *testing.T) {
		encoded, err := conv.EncodeJSONToMsgPack(nil)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack failed: %v", err)
		}
		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("DecodeMsgPackToJSON failed: %v", err)
		}
		if decoded != nil {
			t.Errorf("Expected nil, got %v", decoded)
		}
	})

	t.Run("true", func(t *testing.T) {
		encoded, err := conv.EncodeJSONToMsgPack(true)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack failed: %v", err)
		}
		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("DecodeMsgPackToJSON failed: %v", err)
		}
		if decoded != true {
			t.Errorf("Expected true, got %v", decoded)
		}
	})

	t.Run("false", func(t *testing.T) {
		encoded, err := conv.EncodeJSONToMsgPack(false)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack failed: %v", err)
		}
		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("DecodeMsgPackToJSON failed: %v", err)
		}
		if decoded != false {
			t.Errorf("Expected false, got %v", decoded)
		}
	})

	t.Run("bin8", func(t *testing.T) {
		originalData := []byte{65, 66, 67} // ABC
		encoded, err := conv.EncodeJSONToMsgPack(originalData)
		if err != nil {
			t.Fatalf("Encoding failed: %v", err)
		}

		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("Decoding failed: %v", err)
		}

		if !bytes.Equal(decoded.([]byte), originalData) {
			t.Errorf("Decoded data %v does not match expected %v", decoded, originalData)
		}
	})

	t.Run("bin16", func(t *testing.T) {
		originalData := make([]byte, 300)
		for i := range originalData {
			originalData[i] = byte(i % 256)
		}
		encoded, err := conv.EncodeJSONToMsgPack(originalData)
		if err != nil {
			t.Fatalf("Encoding failed: %v", err)
		}

		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("Decoding failed: %v", err)
		}

		if !bytes.Equal(decoded.([]byte), originalData) {
			t.Errorf("Decoded data %v does not match expected %v", decoded, originalData)
		}
	})

	t.Run("bin32", func(t *testing.T) {
		originalData := make([]byte, 70000)
		for i := range originalData {
			originalData[i] = byte(i % 256)
		}
		encoded, err := conv.EncodeJSONToMsgPack(originalData)
		if err != nil {
			t.Fatalf("Encoding failed: %v", err)
		}

		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("Decoding failed: %v", err)
		}

		if !bytes.Equal(decoded.([]byte), originalData) {
			t.Errorf("Decoded data does not match expected data")
		}
	})

	t.Run("float32", func(t *testing.T) {
		original := float32(math.MaxFloat32)
		encoded, err := conv.EncodeJSONToMsgPack(original)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Errorf("Unexpected error during decoding: %v", err)
		}
		if decoded != original {
			t.Errorf("Expected %v, got %v", original, decoded)
		}
	})

	t.Run("float64", func(t *testing.T) {
		original := math.MaxFloat64
		encoded, err := conv.EncodeJSONToMsgPack(original)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Errorf("Unexpected error during decoding: %v", err)
		}
		if decoded != original {
			t.Errorf("Expected %v, got %v", original, decoded)
		}
	})

	t.Run("str8", func(t *testing.T) {
		// str8 can represent a string that is up to 255 bytes long.
		input := string(make([]byte, 220))
		encoded, err := conv.EncodeJSONToMsgPack(input)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack failed: %v", err)
		}
		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("DecodeMsgPackToJSON failed: %v", err)
		}
		if decoded != input {
			t.Errorf("Expected %v, got %v", input, decoded)
		}
	})

	t.Run("str16", func(t *testing.T) {
		// str16 can represent a string that is up to (2^16)-1 bytes long.
		input := string(make([]byte, 700))
		encoded, err := conv.EncodeJSONToMsgPack(input)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack failed: %v", err)
		}
		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("DecodeMsgPackToJSON failed: %v", err)
		}
		if decoded != input {
			t.Errorf("Expected string of length %d, got length %d", len(input), len(decoded.(string)))
		}
	})

	t.Run("str32", func(t *testing.T) {
		// str32 can represent a string that is up to (2^32)-1 bytes long.
		input := string(make([]byte, 70000))
		encoded, err := conv.EncodeJSONToMsgPack(input)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack failed: %v", err)
		}
		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("DecodeMsgPackToJSON failed: %v", err)
		}
		if decoded != input {
			t.Errorf("Expected string of length %d, got length %d", len(input), len(decoded.(string)))
		}
	})

	t.Run("array16", func(t *testing.T) {
		// array16 can store 2^16 - 1 elements.
		input := make([]interface{}, 500)
		for i := 0; i < 500; i++ {
			input[i] = uint16(i)
		}
		encoded, err := conv.EncodeJSONToMsgPack(input)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack failed: %v", err)
		}
		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("DecodeMsgPackToJSON failed: %v", err)
		}
		decodedArray, ok := decoded.([]interface{})
		if !ok {
			t.Fatalf("Decoded is not an array: %v", decoded)
		}
		for i, v := range input {
			if decodedArray[i] != v {
				t.Errorf("At index %d, got %v want %v", i, decodedArray[i], v)
			}
		}
	})

	t.Run("array32", func(t *testing.T) {
		// array16 can store 2^32 - 1 elements.
		input := make([]interface{}, 70000)
		for i := 0; i < 70000; i++ {
			input[i] = uint32(i)
		}
		encoded, err := conv.EncodeJSONToMsgPack(input)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack failed: %v", err)
		}
		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("DecodeMsgPackToJSON failed: %v", err)
		}
		decodedArray, ok := decoded.([]interface{})
		if !ok {
			t.Fatalf("Decoded is not an array: %v", decoded)
		}
		for i, v := range input {
			if decodedArray[i] != v {
				t.Errorf("At index %d, got %v want %v", i, decodedArray[i], v)
			}
		}
	})

	t.Run("map16", func(t *testing.T) {
		originalData := make(map[string]interface{})
		for i := 0; i < 300; i++ {
			originalData[strconv.Itoa(i)] = uint16(i)
		}
		encodedData, err := conv.EncodeJSONToMsgPack(originalData)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack error: %v", err)
		}

		decodedData, err := conv.DecodeMsgPackToJSON(encodedData)
		if err != nil {
			t.Fatalf("DecodeMsgPackToJSON error: %v", err)
		}

		if !reflect.DeepEqual(originalData, decodedData) {
			t.Errorf("map16: Expected and decoded data do not match")
		}
	})

	// t.Run("map32", func(t *testing.T) {
	// 	largeMap := make(map[string]interface{})
	// 	for i := 0; i < 70000; i++ {
	// 		largeMap[strconv.Itoa(i)] = uint32(i)
	// 	}
	// 	// Encode the large map into MessagePack format.
	// 	encoded, err := conv.EncodeJSONToMsgPack(largeMap)
	// 	if err != nil {
	// 		t.Errorf("Unexpected error during encoding: %v", err)
	// 	}

	// 	// Decode the MessagePack format back into a map.
	// 	decoded, err := conv.DecodeMsgPackToJSON(encoded)
	// 	if err != nil {
	// 		t.Errorf("Unexpected error during decoding: %v", err)
	// 	}

	// 	// Convert the decoded interface{} back to a map for comparison.
	// 	decodedMap, ok := decoded.(map[string]interface{})
	// 	if !ok {
	// 		t.Fatal("Decoded value is not a map")
	// 	}

	// 	// Verify that the decoded map has the same size as the original.
	// 	if len(decodedMap) != len(largeMap) {
	// 		t.Errorf("Expected map of size %d, got map of size %d", len(largeMap), len(decodedMap))
	// 	}

	// 	// Check that all keys and values are equal.
	// 	for k, v := range largeMap {
	// 		if decodedVal, ok := decodedMap[k]; !ok || int(decodedVal.(float64)) != v {
	// 			t.Errorf("Mismatched value for key '%s': expected %d, got %v", k, v, decodedVal)
	// 		}
	// 	}
	// })

	t.Run("map32", func(t *testing.T) {
		originalData := make(map[string]interface{})
		for i := 0; i < 70000; i++ {
			originalData[strconv.Itoa(i)] = uint32(i)
		}
		encodedData, err := conv.EncodeJSONToMsgPack(originalData)
		if err != nil {
			t.Fatalf("EncodeJSONToMsgPack error: %v", err)
		}

		decodedData, err := conv.DecodeMsgPackToJSON(encodedData)
		if err != nil {
			t.Fatalf("DecodeMsgPackToJSON error: %v", err)
		}

		if !reflect.DeepEqual(originalData, decodedData) {
			t.Errorf("map32: Expected and decoded data do not match")
		}
	})

	t.Run("complex data", func(t *testing.T) {
		json_data := map[string]interface{}{
			"byte_val":  []byte("12435grgtrwew9766"),
			"bool_val":  true,
			"null_val":  nil,
			"int_val":   int64(1) << 50,
			"float_val": 12.344443423157,
			"str_val":   strings.Repeat("world", 100),
			"list_val":  []interface{}{1},
			"dict_val":  map[string]interface{}{"a": 1, "b": 2},
		}

		encoded, err := conv.EncodeJSONToMsgPack(json_data)
		if err != nil {
			t.Fatalf("Encoding failed: %v", err)
		}

		decoded, err := conv.DecodeMsgPackToJSON(encoded)
		if err != nil {
			t.Fatalf("Decoding failed: %v", err)
		}

		if !reflect.DeepEqual(decoded, json_data) {
			t.Errorf("Decoded data %#v does not match expected %#v", decoded, json_data)
		}
	})

}
