package CoreHexUint

//16进制转[]byte
func Uint16ToBytes(n uint16) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
	}
}

//[]byte转16进制
func BytesToUint16(array []byte) uint16 {
	var data uint16 = 0
	for i := 0; i < len(array); i++ {
		data = data + uint16(uint(array[i])<<uint(8*i))
	}

	return data
}

//32进制转[]byte
func Uint32ToBytes(n uint32) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
	}
}

//64进制转[]byte
func Uint64ToBytes(n uint64) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
		byte(n >> 32),
		byte(n >> 40),
		byte(n >> 48),
		byte(n >> 56),
	}
}
