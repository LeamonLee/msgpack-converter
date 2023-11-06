package msgpackconverter

const (
	fixPositiveIntMinType = 0x00
	fixPositiveIntMaxType = 0x7f

	fixMapMinType = 0x80
	fixMapMaxType = 0x8f

	fixArrayMinType = 0x90
	fixArrayMaxType = 0x9f

	fixStrMinType = 0xa0
	fixStrMaxType = 0xbf

	nilType = 0xc0

	falseType = 0xc2
	trueType  = 0xc3

	bin8Type  = 0xc4
	bin16Type = 0xc5
	bin32Type = 0xc6

	float32Type = 0xca
	float64Type = 0xcb

	uint8Type  = 0xcc
	uint16Type = 0xcd
	uint32Type = 0xce
	uint64Type = 0xcf
	int8Type   = 0xd0
	int16Type  = 0xd1
	int32Type  = 0xd2
	int64Type  = 0xd3

	str8Type  = 0xd9
	str16Type = 0xda
	str32Type = 0xdb

	array16Type = 0xdc
	array32Type = 0xdd

	map16Type             = 0xde
	map32Type             = 0xdf
	fixNegativeIntMinType = 0xe0
	fixNegativeIntMaxType = 0xff
)
