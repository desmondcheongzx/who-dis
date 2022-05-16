package pkg

import (
	"fmt"
	"log"
	"net"
)

type DNSClient struct {
	cache map[string]net.IP
}

func NewDNSClient() *DNSClient {
	return &DNSClient{
		cache: make(map[string]net.IP),
	}
}

func (client *DNSClient) Query(dn string) error {
	addr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	question := Question{
		qname:  dn,
		qtype:  Q_A,
		qclass: Q_IN,
	}
	message := newQuery([]Question{question}, nil, nil, nil)
	payload, err := message.serialize()
	if err != nil {
		return err
	}
	if _, err = conn.Write(payload); err != nil {
		return err
	}
	buf := make([]byte, 512)
	bytesRead, err := conn.Read(buf)
	if err != nil {
		return err
	}

	replyHdr := &Header{}
	startidx := 12 // End of header
	replyHdr.deserialize(buf[:startidx])
	fmt.Println(";; QUESTION SECTION:")
	for i := 0; i < int(replyHdr.qdcount); i++ {
		q := &Question{}
		n, err := q.deserialize(buf, startidx, bytesRead)
		if err != nil {
			return err
		}
		fmt.Printf("%v\t\t\tIN\tA\n", q.qname)
		startidx += n
	}
	fmt.Println()
	fmt.Println(";; ANSWER SECTION:")
	for i := 0; i < int(replyHdr.ancount); i++ {
		rr := &ResourceRecord{}
		n, err := rr.deserialize(buf, startidx, bytesRead)
		if err != nil {
			return err
		}
		startidx += n
		if rr.rrType == RR_CNAME {
			alias, _, err := decodeDomainName(buf, startidx-len(rr.rdata), bytesRead)
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("%v\t\t%v\tIN\tCNAME\t%v\n", rr.name, rr.ttl, alias)
		} else {
			fmt.Printf("%v\t\t%v\tIN\tA\t%v\n", rr.name, rr.ttl, net.IP(rr.rdata))
		}
	}
	return nil
}
