package CoreHexUint

import (
	"bytes"
	"encoding/binary"
)

func IntToBytes(n int) (dataByte []byte, err error) {
	data := int64(n)
	buf := bytes.NewBuffer([]byte{})
	err = binary.Write(buf, binary.BigEndian, data)
	if err != nil {
		return
	}
	dataByte = buf.Bytes()
	return
}

func BytesToInt(bys []byte) (dataInt int, err error) {
	buf := bytes.NewBuffer(bys)
	var data int64
	err = binary.Read(buf, binary.BigEndian, &data)
	if err != nil {
		return
	}
	dataInt = int(data)
	return
}
