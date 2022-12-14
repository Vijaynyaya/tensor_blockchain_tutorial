package blockchain

import (
	"bytes"
	"encoding/binary"
	"log"
)

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	// big endian: leftmost bit is the most significant bit 💪
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// panics after logging the error
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
