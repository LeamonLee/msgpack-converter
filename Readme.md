# MessagePack-JSON Converter

This repository contains two versions of a MessagePack-JSON converter: one implemented in Python and the other in Go. These converters allow for encoding JSON to MessagePack and decoding MessagePack to JSON back without using any third-party libraries.

## Installation

Clone the repository using the following command:

```bash
git clone https://github.com/LeamonLee/msgpack-converter.git
cd msgpack-converter
```

## Go Version

The Go version provides two functions:

- `EncodeJSONToMsgPack(jsonData []byte) ([]byte, error)`: Encodes JSON data to MessagePack format.
- `DecodeMsgPackToJSON(msgpackData []byte) ([]byte, error)`: Decodes MessagePack data to JSON format.

### Running Go Unit Tests

To run the unit tests for the Go version, follow these steps:

1. Navigate to the Go implementation directory (if applicable):
    ```bash
    cd golang/msgpackconverter
    ```
2. Run the tests using the `go test` command:
    ```bash
    go test
    ```

### Example

In Go, use the MsgPackConverter struct from the msgpackconverter package to convert between JSON and MessagePack.

```go
package main

import (
	"fmt"
	"strings"
	"github.com/LeamonLee/msgpack-converter/msgpackconverter"
)

func main() {
	converter := &msgpackconverter.MsgPackConverter{}
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

	encodedData, err := converter.EncodeJSONToMsgPack(json_data)
	if err != nil {
		fmt.Printf("Encoding failed: %v", err)
		return
	}

	decodedData, err := converter.DecodeMsgPackToJSON(encodedData)
	if err != nil {
		fmt.Printf("Decoding failed: %v", err)
		return
	}

	// Output the decoded data
	fmt.Println(decodedData)
}

```

## Python Version

The Python version provides two functions:

- `encode_json_to_msgpack(json_data)`: Encodes JSON data to MessagePack format.
- `decode_msgpack_to_json(msgpack_data)`: Decodes MessagePack data to JSON format.

### Running Python Unit Tests

To run the unit tests for the Python version, follow these steps:

1. Navigate to the Python implementation directory (if applicable):
    ```bash
    cd python
    ```
2. Run the tests using the unittest module that comes with the standard library of Python:
    ```bash
    python -m unittest test_msgpackconverter.py
    ```

Replace `python` with `python3` if you are running on a system where Python 3 is not the default version.

### Example

In Python, you can use the `MsgPackConverter` class to convert between JSON and MessagePack.

```python
from msgPackConverter import MsgPackConverter  # Replace with your actual import

converter = MsgPackConverter()
json_data = {
    "byte_val": b"12435grgtrwew9766",
    "bool_val": True, 
    "null_val": None, 
    "int_val": 2**50, 
    "float_val": 12.344443423157,
    "str_val": "world" * 100,
    "list_val": [1], 
    "dict_val": {"a": 1, "b": 2}
}

# Encoding the JSON data to MessagePack format
encoded_data = converter.encode_json_to_msgpack(json_data)

# Decoding the MessagePack data back to JSON format
decoded_data = converter.decode_msgpack_to_json(encoded_data)
```
