import unittest
from msgPackConverter import MsgPackConverter
import struct

class TestMsgPackConverter(unittest.TestCase):
    def setUp(self):
        self.converter = MsgPackConverter()

    def test_positive_fixint(self):
        encoded = self.converter.encode_json_to_msgpack(127)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded), 127)

    def test_negative_fixint(self):
        encoded = self.converter.encode_json_to_msgpack(-1)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded), -1)
        
    def test_fixmap(self):
        encoded = self.converter.encode_json_to_msgpack({"a": 1})
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded), {"a": 1})

    def test_fixarray(self):
        encoded = self.converter.encode_json_to_msgpack([1, 2, 3])
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded), [1, 2, 3])

    def test_fixstr(self):
        encoded = self.converter.encode_json_to_msgpack("hello")
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded), "hello")
        
    def test_nil(self):
        encoded = self.converter.encode_json_to_msgpack(None)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded), None)
        
    def test_false(self):
        encoded = self.converter.encode_json_to_msgpack(False)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded), False)

    def test_true(self):
        encoded = self.converter.encode_json_to_msgpack(True)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded), True)

    def test_bin8(self):
        original_data = bytes([65, 66, 67])
        encoded_data = self.converter.encode_json_to_msgpack(original_data)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded_data), original_data)
        
    def test_bin16(self):
        original_data = bytes([i % 256 for i in range(300)])
        encoded_data = self.converter.encode_json_to_msgpack(original_data)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded_data), original_data)
        
    def test_bin32(self):
        original_data = bytes([i % 256 for i in range(70000)])
        encoded_data = self.converter.encode_json_to_msgpack(original_data)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded_data), original_data)
        
    def test_float32(self):
        # Assuming MsgPackConverter can handle float32 and float64
        # original = struct.unpack('f', struct.pack('f', 3.4028234663852886e+38))[0]  # Max float32 value
        original = 3.4028234663852886e+38
        encoded_data = self.converter.encode_json_to_msgpack(original)
        decoded_data = self.converter.decode_msgpack_to_json(encoded_data)
        self.assertAlmostEqual(decoded_data, original, places=5)

    def test_float64(self):
        # original = struct.unpack('d', struct.pack('d', 1.7976931348623157e+308))[0]  # Max float64 value
        original = 1.7976931348623157e+308
        encoded_data = self.converter.encode_json_to_msgpack(original)
        decoded_data = self.converter.decode_msgpack_to_json(encoded_data)
        self.assertEqual(decoded_data, original)

    def test_str8(self):
        original_data = "a" * 220
        encoded_data = self.converter.encode_json_to_msgpack(original_data)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded_data), original_data)
        
    def test_str16(self):
        original_data = "a" * 700
        encoded_data = self.converter.encode_json_to_msgpack(original_data)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded_data), original_data)
        
    def test_str32(self):
        original_data = "a" * 70000
        encoded_data = self.converter.encode_json_to_msgpack(original_data)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded_data), original_data)
        
    def test_array16(self):
        original_data = [i for i in range(500)]
        encoded_data = self.converter.encode_json_to_msgpack(original_data)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded_data), original_data)
        
    def test_array32(self):
        original_data = [i for i in range(70000)]
        encoded_data = self.converter.encode_json_to_msgpack(original_data)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded_data), original_data)
        
    def test_map16(self):
        original_data = {str(i): i for i in range(300)}
        encoded_data = self.converter.encode_json_to_msgpack(original_data)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded_data), original_data)
        
    def test_map32(self):
        original_data = {str(i): i for i in range(70000)}
        encoded_data = self.converter.encode_json_to_msgpack(original_data)
        self.assertEqual(self.converter.decode_msgpack_to_json(encoded_data), original_data)

    def test_complex_data(self):
        json_data = {
            "byte_val": b"12435grgtrwew9766",
            "bool_val": True, 
            "null_val": None, 
            "int_val": 2**50, 
            "float_val": 12.344443423157,
            "str_val": "world"*100,
            "list_val": [1], 
            "dict_val": {"a": 1, "b": 2}
        }
        
        # Encoding the JSON data to MessagePack format
        encoded_data = self.converter.encode_json_to_msgpack(json_data)
        
        # Decoding the MessagePack data back to JSON format
        decoded_data = self.converter.decode_msgpack_to_json(encoded_data)
        
        # Asserting if the decoded data matches the original JSON data
        self.assertEqual(decoded_data, json_data)

if __name__ == '__main__':
    unittest.main()
