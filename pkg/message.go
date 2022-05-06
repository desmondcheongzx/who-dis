package pkg

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

// func decodeQuestion(data []byte) (*Question, n, error) {
// 	// Get the qname.
// 	qname, n, err := decodeDN(data)
// 	if err != nil {
// 		return nil, -1, err
// 	}
// 	// Get the qtype and qclass.
// 	qtype := QType(data[n:n+2])
// }

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
