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

type RRType = uint16

const (
	A     = 1
	NS    = 2
	MD    = 3
	MF    = 4
	CNAME = 5
	SOA   = 6
	MB    = 7
	MG    = 8
	MR    = 9
	NULL  = 10
	WKS   = 11
	PTR   = 12
	HINFO = 13
	MINFO = 14
	MX    = 15
	TXT   = 16
)

type RRClass = uint16

const (
	IN = 1
	CS = 2
	CH = 3
	HS = 4
)

type Question struct {
	qname  string
	qtype  QType
	qclass QClass
}

type QType = uint16
type QClass = uint16
