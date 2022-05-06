package pkg

import "encoding/binary"

func htons(data uint16) []byte {
	buf16 := make([]byte, 2)
	binary.BigEndian.PutUint16(buf16, data)
	return buf16
}

func htonl(data uint32) []byte {
	buf32 := make([]byte, 4)
	binary.BigEndian.PutUint32(buf32, data)
	return buf32
}

func ntohs(data []byte) uint16 {
	if len(data) < 2 {
		return 0
	}
	return binary.BigEndian.Uint16(data)
}

func ntohl(data []byte) uint32 {
	if len(data) < 4 {
		return 0
	}
	return binary.BigEndian.Uint32(data)
}
