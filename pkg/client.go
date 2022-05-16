package pkg

import (
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

func (client *DNSClient) Query(dn string) (net.IP, error) {
	addr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	if _, err = conn.Write(payload); err != nil {
		return nil, err
	}
	buf := make([]byte, 65535)
	bytesRead, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	replyHdr := &Header{}
	startidx := 12 // End of header
	replyHdr.deserialize(buf[:startidx])
	for i := 0; i < int(replyHdr.qdcount); i++ {
		q := &Question{}
		n, err := q.deserialize(buf, startidx, bytesRead)
		if err != nil {
			return nil, err
		}
		startidx += n
	}
	for i := 0; i < int(replyHdr.ancount); i++ {
		rr := &ResourceRecord{}
		n, err := rr.deserialize(buf, startidx, bytesRead)
		if err != nil {
			return nil, err
		}
		startidx += n
		if rr.name == dn {
			return net.IP(rr.rdata), nil
		}
	}
	return nil, nil
}
