package pkg

import (
	"encoding/binary"
	"errors"
	"strings"
)

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

// Function to parse variable-length domain names from bytes.
func decodeDN(data []byte) (string, int, error) {
	// Initialize.
	var sb strings.Builder
	var n int
	// Iterate over data.
	for _, b := range data {
		if b == 0 {
			break
		} else if n == 0 {
			n = int(b)
			sb.WriteString(".")
		} else {
			n = n - 1
			sb.WriteByte(b)
		}
	}
	// If n isn't 0, data was malformed; else return.
	dn := sb.String()
	if n != 0 || len(dn) == 0 {
		return "", -1, errors.New("domain name data malformed")
	}
	return dn[1:], len(dn), nil
}

// Function to encode a domain name as bytes.
func encodeDN(dn string) ([]byte, error) {
	// Split toks, append in specified manner.
	data, toks := make([]byte, 0), strings.Split(dn, ".")
	for _, t := range toks {
		if len(t) == 0 {
			return nil, errors.New("domain name malformed")
		}
		data = append(data, byte(len(t)))
		data = append(data, []byte(t)...)
	}
	data = append(data, byte(0))
	return data, nil
}
