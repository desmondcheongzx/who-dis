package pkg

import (
	"errors"
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

func newQuery(questions []Question, answers []ResourceRecord, authorities []ResourceRecord, additionals []ResourceRecord) *Message {
	header := Header{
		id:      genRandomID(),
		qr:      false, // Query
		opcode:  0,     // Standard query
		rd:      true,  // Recursion desired
		qdcount: uint16(len(questions)),
		ancount: uint16(len(answers)),
		nscount: uint16(len(authorities)),
		arcount: uint16(len(additionals)),
	}
	return &Message{
		header:      header,
		questions:   questions,
		answers:     answers,
		authorities: authorities,
		additionals: additionals,
	}
}

func (msg *Message) serialize() ([]byte, error) {
	buf := make([]byte, 0)
	buf = append(buf, msg.header.serialize()...)
	for _, q := range msg.questions {
		qbuf, err := q.serialize()
		if err != nil {
			return nil, err
		}
		buf = append(buf, qbuf...)
	}
	for _, anws := range msg.answers {
		abuf, err := anws.serialize()
		if err != nil {
			return nil, err
		}
		buf = append(buf, abuf...)
	}
	for _, auths := range msg.authorities {
		abuf, err := auths.serialize()
		if err != nil {
			return nil, err
		}
		buf = append(buf, abuf...)
	}
	for _, adds := range msg.additionals {
		abuf, err := adds.serialize()
		if err != nil {
			return nil, err
		}
		buf = append(buf, abuf...)
	}
	return buf, nil
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
	hdr.id = ntohs(buf[0:2])
	flags1 := uint8(buf[2])
	hdr.qr = flags1&(1<<7) != 0
	hdr.opcode = (flags1 >> 3) & 0b1111
	hdr.aa = flags1&(1) != 0
	flags2 := uint8(buf[3])
	hdr.tc = flags2&(1<<7) != 0
	hdr.rd = flags2&(1<<6) != 0
	hdr.ra = flags2&(1<<5) != 0
	hdr.rcode = flags1 & 0b1111
	hdr.qdcount = ntohs(buf[4:6])
	hdr.ancount = ntohs(buf[6:8])
	hdr.nscount = ntohs(buf[8:10])
	hdr.arcount = ntohs(buf[10:12])
}

type ResourceRecord struct {
	name    string
	rrType  RRType
	rrClass RRClass
	ttl     uint32
	rdlen   uint16
	rdata   []byte
}

func (rr *ResourceRecord) serialize() ([]byte, error) {
	buf := make([]byte, 0)
	nData, err := encodeDomainName(rr.name)
	if err != nil {
		return nil, err
	}
	buf = append(buf, nData...)
	buf = append(buf, htons(rr.rrType)...)
	buf = append(buf, htons(rr.rrType)...)
	buf = append(buf, htonl(rr.ttl)...)
	buf = append(buf, htons(rr.rdlen)...)
	buf = append(buf, rr.rdata...)
	return buf, nil
}

func (rr *ResourceRecord) deserialize(data []byte, idx int, maxlen int) (int, error) {
	// Get the qname.
	name, n, err := decodeDomainName(data, idx, maxlen)
	if err != nil {
		return -1, err
	}
	// If not enough data, error.
	if len(data) < idx+n+10 {
		return -1, errors.New("data malformed; too short")
	}
	// Get the qtype and qclass.
	rrType := ntohs(data[idx+n : idx+n+2])
	rrClass := ntohs(data[idx+n+2 : idx+n+4])
	ttl := ntohl(data[idx+n+4 : idx+n+8])
	rdlen := ntohs(data[idx+n+8 : idx+n+10])
	if len(data) < idx+n+10+int(rdlen) {
		return -1, errors.New("data malformed; too short")
	}
	rdata := append(make([]byte, 0), data[idx+n+10:idx+n+10+int(rdlen)]...)
	// Assign and return.
	rr.name = name
	rr.rrType = rrType
	rr.rrClass = rrClass
	rr.ttl = ttl
	rr.rdlen = rdlen
	rr.rdata = rdata
	return n + 10 + int(rdlen), nil
}

type Question struct {
	qname  string
	qtype  QType
	qclass QClass
}

func (q *Question) serialize() ([]byte, error) {
	buf := make([]byte, 0)
	qnData, err := encodeDomainName(q.qname)
	if err != nil {
		return nil, err
	}
	buf = append(buf, qnData...)
	buf = append(buf, htons(q.qtype)...)
	buf = append(buf, htons(q.qclass)...)
	return buf, nil
}

func (q *Question) deserialize(data []byte, idx int, maxlen int) (int, error) {
	// Get the qname.
	qname, n, err := decodeDomainName(data, idx, maxlen)
	if err != nil {
		return -1, err
	}
	// If not enough data, error.
	if len(data) < n+4 {
		return -1, errors.New("data malformed; too short")
	}
	// Get the qtype and qclass.
	qtype := ntohs(data[idx+n : idx+n+2])
	qclass := ntohs(data[idx+n+2 : idx+n+4])
	// Assign and return.
	q.qname = qname
	q.qtype = qtype
	q.qclass = qclass
	return n + 4, nil
}
