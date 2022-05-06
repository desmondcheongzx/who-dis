package pkg

import (
	"errors"
	"strings"
)

// +---------------------+
// |        Header       |
// +---------------------+
// |       Question      | the question for the name server
// +---------------------+
// |        Answer       | RRs answering the question
// +---------------------+
// |      Authority      | RRs pointing toward an authority
// +---------------------+
// |      Additional     | RRs holding additional information
// +---------------------+

// The header contains the following fields:

//                                     1  1  1  1  1  1
//       0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                      ID                       |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    QDCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    ANCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    NSCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//     |                    ARCOUNT                    |
//     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type Message struct {
	header      Header
	questions   []Question
	answers     []ResourceRecord
	authorities []ResourceRecord
	additionals []ResourceRecord
}

type Header struct {
	id     uint16
	qr     bool  // query or response
	opcode uint8 // 4 bits
	aa     bool  // authoritative answer
	tc     bool  // TrunCation
	rd     bool  // recursion desired
	ra     bool  // recursion available
	// Note: Z is reservered for future use and should be 0
	rcode   uint8  //response code
	qdcount uint16 // question count
	ancount uint16 // answer count
	nscount uint16 // autorities count
	arcount uint16 // additionals count
}

type ResourceRecord struct {
	name    string
	rrType  RRType
	rrClass RRClass
	ttl     uint32
	rdlen   uint16
	rdata   []byte
}

type Question struct {
	qname  string
	qtype  QType
	qclass QClass
}

// Function to parse variable-length domain names from bytes.
func decodeDN(data []byte) (string, error) {
	// Initialize.
	var sb strings.Builder
	var n int
	// Iterate over data.
	for _, b := range data {
		if n == 0 {
			n = int(b)
			sb.WriteString(".")
		} else {
			n = n - 1
			sb.WriteByte(b)
		}
	}
	// If n isn't 0, data was malformed; else return.
	if n != 0 || len(sb.String()) == 0 {
		return "", errors.New("domain name data malformed")
	}
	return sb.String()[1:], nil
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
	return data, nil
}

func (hdr *Header) serialize() []byte {
	buf := make([]byte, 0)
	buf = append(buf, htons(hdr.id)...)
	var flags1 uint8 = 0
	if hdr.qr {
		flags1 |= 1 << 7
	}
	flags1 |= uint8(hdr.opcode << 3)
	if hdr.aa {
		flags1 |= 1 << 2
	}
	if hdr.tc {
		flags1 |= 1 << 1
	}
	if hdr.rd {
		flags1 |= 1
	}
	buf = append(buf, []byte{flags1}...)
	var flags2 uint8 = 8
	if hdr.ra {
		flags1 |= 1
	}
	flags2 |= hdr.rcode
	buf = append(buf, []byte{flags2}...)
	buf = append(buf, htons(hdr.qdcount)...)
	buf = append(buf, htons(hdr.ancount)...)
	buf = append(buf, htons(hdr.nscount)...)
	buf = append(buf, htons(hdr.arcount)...)
	return buf
}

func (hdr *Header) deserialize(buf []byte) {
}
