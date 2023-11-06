import struct
import math


class MsgPackConverter:
    def __init__(self):
        pass

    @staticmethod
    def __bytes_to_int(byte_array):
        result = 0
        for b in byte_array:
            result = result * 256 + int(b)
        return result
    
    @staticmethod
    def __encode_bytes(value):
        def int_to_bytes(n, length):
            return [(n >> (i * 8)) & 0xff for i in range(length)]
        
        length = len(value)
        if length < (1 << 8):
            return [0xc4] + [length] + list(value)
        elif length < (1 << 16):
            return [0xc5] + int_to_bytes(length, 2) + list(value)
        else:
            return [0xc6] + int_to_bytes(length, 4) + list(value)

    @staticmethod
    def __encode_bool(value):
        msgpack_data = bytearray()
        msgpack_data.append(0xc3 if value else 0xc2)

        return msgpack_data

    @staticmethod
    def __encode_integer(value):
        """Encode an integer value."""
        msgpack_data = bytearray()

        if 0 <= value <= 127:
            msgpack_data.append(value)
        elif 128 <= value <= 255:  # Unsigned 8-bit int
            msgpack_data.append(0xcc)
            msgpack_data.append(value)
        elif 256 <= value <= 65535:  # Unsigned 16-bit int
            msgpack_data.append(0xcd)
            msgpack_data.extend(value.to_bytes(2, 'big'))
        elif 65536 <= value <= 4294967295:  # Unsigned 32-bit int
            msgpack_data.append(0xce)
            msgpack_data.extend(value.to_bytes(4, 'big'))
        elif value > 4294967295:  # uint 64
            msgpack_data.append(0xcf)
            msgpack_data.extend(value.to_bytes(8, 'big'))
        elif value < 0:
            if -128 <= value <= -1:
                msgpack_data.append(0xd0)
                msgpack_data.append(value & 0xff)
            elif -32768 <= value:
                msgpack_data.append(0xd1)
                msgpack_data.extend(struct.pack(">h", value))
                # msgpack_data.extend(value.to_bytes(2, 'big'))
            elif -2147483648 <= value:
                msgpack_data.append(0xd2)
                msgpack_data.extend(struct.pack(">i", value))
                # msgpack_data.extend(value.to_bytes(4, 'big'))
            else:
                msgpack_data.append(0xd3)
                msgpack_data.extend(struct.pack(">q", value))

        return msgpack_data

    @staticmethod
    def __encode_float(value):
        """Encode a float value."""
        def float32(value):
            return struct.unpack('f', struct.pack('f', value))[0]
        
        msgpack_data = bytearray()

        if math.isnan(value):
            # Encode NaN
            msgpack_data.append(0xca)
            msgpack_data.extend(struct.pack("!f", value))
        elif math.isinf(value):
            # Encode Infinity and -Infinity
            if value > 0:  # Positive Infinity
                msgpack_data.append(0xca)
                msgpack_data.extend(struct.pack("!f", float('inf')))
            else:  # Negative Infinity
                msgpack_data.append(0xca)
                msgpack_data.extend(struct.pack("!f", float('-inf')))
        else:
            if value == float32(value):  # Check if it's a float32
                msgpack_data.append(0xca)
                msgpack_data.extend(struct.pack("!f", value))
            else:  # Otherwise consider it as float64
                msgpack_data.append(0xcb)
                msgpack_data.extend(struct.pack("!d", value))

        return msgpack_data

    @staticmethod
    def __encode_string(value):
        """Encode a string value."""
        msgpack_data = bytearray()

        str_length = len(value.encode())
        if str_length <= 31:
            msgpack_data.append(0xa0 | str_length)
        elif str_length <= (1 << 8) - 1:
            msgpack_data.append(0xd9)
            msgpack_data.append(str_length)
        elif str_length <= (1 << 16) - 1:
            msgpack_data.append(0xda)
            msgpack_data.extend(struct.pack(">H", str_length))
        else:
            msgpack_data.append(0xdb)
            msgpack_data.extend(struct.pack(">I", str_length))

        msgpack_data.extend(value.encode())

        return msgpack_data

    def __encode_list(self, value):
        """Encode a list."""
        msgpack_data = bytearray()

        list_length = len(value)
        if list_length <= 15:
            msgpack_data.append(0x90 | list_length)
        elif list_length <= 0xffff:
            msgpack_data.append(0xdc)
            msgpack_data.extend(list_length.to_bytes(2, 'big'))
        else:
            msgpack_data.append(0xdd)
            msgpack_data.extend(list_length.to_bytes(4, 'big'))
        for item in value:
            msgpack_data.extend(self.encode_json_to_msgpack(item))

        return msgpack_data

    def __encode_dict(self, value):
        """Encode a dictionary."""
        msgpack_data = bytearray()

        dict_length = len(value)
        # print(f"__encode_dict dict_length:{dict_length}")
        if dict_length <= 15:
            msgpack_data.append(0x80 | dict_length)
        elif dict_length <= 0xffff:
            # print(f"__encode_dict 0xde")
            msgpack_data.append(0xde)
            msgpack_data.extend(dict_length.to_bytes(2, 'big'))
        else:
            # print(f"__encode_dict 0xdf")
            msgpack_data.append(0xdf)
            msgpack_data.extend(dict_length.to_bytes(4, 'big'))
        for key, value in value.items():
            msgpack_data.extend(self.encode_json_to_msgpack(key))
            msgpack_data.extend(self.encode_json_to_msgpack(value))

        return msgpack_data

    # Encoding Methods
    def encode_json_to_msgpack(self, json_data):
        """Encode JSON data into MessagePack format."""
        try:
            if json_data is None:
                return [0xc0]
            elif isinstance(json_data, bytes):
                return self.__encode_bytes(json_data)
            elif isinstance(json_data, bool):
                return self.__encode_bool(json_data)
            elif isinstance(json_data, int):
                return self.__encode_integer(json_data)
            elif isinstance(json_data, float):
                return self.__encode_float(json_data)
            elif isinstance(json_data, str):
                return self.__encode_string(json_data)
            elif isinstance(json_data, list):
                return self.__encode_list(json_data)
            elif isinstance(json_data, dict):
                return self.__encode_dict(json_data)
            else:
                raise ValueError(f"Unsupported data type: {type(json_data)}")
        except Exception as e:
            print(f"An error occurred while encoding: {e}")
    

    # Decoding Methods
    def decode_msgpack_to_json(self, msgpack_data):
        """Decode MessagePack data into JSON format."""
        try:
            msgpack_data = bytearray(msgpack_data)
            first_byte = msgpack_data.pop(0)

            if first_byte == 0xc0:
                return None

            elif first_byte == 0xc2:
                return False

            elif first_byte == 0xc3:
                return True

            elif first_byte == 0xc4:
                return self.__decode_bytes_8(msgpack_data)

            elif first_byte == 0xc5:
                return self.__decode_bytes_16(msgpack_data)

            elif first_byte == 0xc6:
                return self.__decode_bytes_32(msgpack_data)

            elif 0x00 <= first_byte <= 0x7f:    # positive fixint
                return first_byte
            elif 0xe0 <= first_byte <= 0xff:    # negative fixint
                return first_byte - 256

            elif 0x80 <= first_byte <= 0x8f:    # fixmap
                return self.__decode_fixmap(msgpack_data, first_byte)

            elif 0x90 <= first_byte <= 0x9f:    # fixarray
                return self.__decode_fixarr(msgpack_data, first_byte)

            elif 0xa0 <= first_byte <= 0xbf:    # fixstr
                str_length = first_byte & 0x1f
                return msgpack_data[:str_length].decode()

            elif first_byte == 0xca:  # Float 32
                return struct.unpack("!f", msgpack_data[:4])[0]

            elif first_byte == 0xcb:  # Float 64
                return struct.unpack("!d", msgpack_data[:8])[0]

            elif 0xcc <= first_byte <= 0xcf:  # Handling unsigned integers
                return self.__decode_unsigned_integer(msgpack_data, first_byte)

            elif 0xd0 <= first_byte <= 0xd3:
                return self.__decode_signed_integer(msgpack_data, first_byte)

            elif first_byte == 0xd9 or (0xda <= first_byte <= 0xdb):
                return self.__decode_string(msgpack_data, first_byte)

            elif 0xdc <= first_byte <= 0xdd:
                return self.__decode_list(msgpack_data, first_byte)

            elif 0xde <= first_byte <= 0xdf:
                return self.__decode_dict(msgpack_data, first_byte)

            else:
                raise ValueError("Unsupported type")
        except Exception as e:
            print(f"An error occurred while decoding: {e}")

    @staticmethod
    def __decode_bytes_8(msgpack_data):
        bin_length = msgpack_data.pop(0)
        return bytes(msgpack_data[:bin_length])

    @staticmethod
    def __decode_bytes_16(msgpack_data):
        bin_length = msgpack_data.pop(0) << 8 | msgpack_data.pop(0)
        return bytes(msgpack_data[:bin_length])

    @staticmethod
    def __decode_bytes_32(msgpack_data):
        bin_length = (msgpack_data.pop(0) << 24 | msgpack_data.pop(0) << 16 |
                        msgpack_data.pop(0) << 8 | msgpack_data.pop(0))
        return bytes(msgpack_data[:bin_length])

    def __decode_fixmap(self, msgpack_data, first_byte):
        dict_length = first_byte & 0x0f
        decoded_dict = {}
        for _ in range(dict_length):
            key = self.decode_msgpack_to_json(msgpack_data)
            msgpack_data = msgpack_data[len(self.encode_json_to_msgpack(key)):]
            value = self.decode_msgpack_to_json(msgpack_data)
            msgpack_data = msgpack_data[len(self.encode_json_to_msgpack(value)):]
            decoded_dict[key] = value
        return decoded_dict

    def __decode_fixarr(self, msgpack_data, first_byte):
        list_length = first_byte & 0x0f
        decoded_list = []
        for _ in range(list_length):
            item = self.decode_msgpack_to_json(msgpack_data)
            msgpack_data = msgpack_data[len(self.encode_json_to_msgpack(item)):]
            decoded_list.append(item)
        return decoded_list

    @staticmethod
    def __decode_unsigned_integer(msgpack_data, first_byte):
        """Decode an unsigned integer value."""
        sizes = [1, 2, 4, 8]
        size = sizes[first_byte - 0xcc]
        num = int.from_bytes(msgpack_data[:size], byteorder='big', signed=False)
        del msgpack_data[:size]
        return num

    @staticmethod
    def __decode_signed_integer(msgpack_data, first_byte):
        """Decode an signed integer value."""
        sizes = [1, 2, 4, 8]
        size = sizes[first_byte - 0xd0]
        num = int.from_bytes(msgpack_data[:size], byteorder='big', signed=True)
        del msgpack_data[:size]
        return num

    @staticmethod
    def __decode_string(msgpack_data, first_byte):
        """Decode a string value."""
        sizes = [1, 2, 4]
        size = sizes[(first_byte & 0x0f) - 9]
        str_length = int.from_bytes(msgpack_data[:size], byteorder='big')
        del msgpack_data[:size]
        string = msgpack_data[:str_length].decode()
        del msgpack_data[:str_length]
        return string


    def __decode_list(self, msgpack_data, first_byte):
        """Decode a list."""
        sizes = [2,4]
        size = sizes[first_byte - 0xdc]
        list_length = int.from_bytes(msgpack_data[:size], 'big')
        msgpack_data = msgpack_data[size:]
        decoded_list = []
        for _ in range(list_length):
            item = self.decode_msgpack_to_json(msgpack_data)
            msgpack_data = msgpack_data[len(self.encode_json_to_msgpack(item)):]
            decoded_list.append(item)
        return decoded_list

    def __decode_dict(self, msgpack_data, first_byte):
        """Decode a dictionary."""
        sizes = [2,4]
        size = sizes[first_byte - 0xde]
        dict_length = int.from_bytes(msgpack_data[:size], 'big')
        msgpack_data = msgpack_data[size:]
        decoded_dict = {}
        for _ in range(dict_length):
            key = self.decode_msgpack_to_json(msgpack_data)
            msgpack_data = msgpack_data[len(self.encode_json_to_msgpack(key)):]
            value = self.decode_msgpack_to_json(msgpack_data)
            msgpack_data = msgpack_data[len(self.encode_json_to_msgpack(value)):]
            decoded_dict[key] = value
        return decoded_dict

    def test_conversion(self):
        """Method to test the encoding and decoding."""
        test_cases = [
            {"nan_val": math.nan, "inf_val": math.inf, "neg_inf_val": -math.inf},
            {"list_val": [1, 2, math.nan, math.inf, -math.inf]},
            {"dict_val": {"a": math.nan, "b": math.inf, "c": -math.inf}},
            {"mixed_val": [math.nan, {"a": 1, "b": math.inf}, -math.inf, "test"]}
        ]

        for i, test_case in enumerate(test_cases, 1):
            try:
                print(f"Test Case {i}:")
                msgpack_data = self.encode_json_to_msgpack(test_case)
                print("Encoded MessagePack data:", msgpack_data)

                decoded_json_data = self.decode_msgpack_to_json(msgpack_data)
                print("Decoded JSON data:", decoded_json_data)
                print("-" * 40)
            except Exception as e:
                print(f"An error occurred in the test case: {e}")
                print("-" * 40)

if __name__ == "__main__":
    msgPacker = MsgPackConverter()
    msgPacker.test_conversion()