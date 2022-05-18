package pkg

import (
	"fmt"
	"log"
	"net"
)

var ROOT_IP string = "199.9.14.201"

type DNSClient struct {
	cache map[string]net.IP
}

func NewDNSClient() *DNSClient {
	return &DNSClient{
		cache: make(map[string]net.IP),
	}
}

func recursiveQuery(dn string, ns string) error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:53", ns))
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
	fmt.Println()
	fmt.Println("### Response from", ns, "###")
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
	answers := make([]*ResourceRecord, 0)
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
			answers = append(answers, rr)
			fmt.Printf("%v\t\t%v\tIN\tA\t%v\n", rr.name, rr.ttl, net.IP(rr.rdata))
		}
	}
	fmt.Println()
	fmt.Println(";; AUTHORITY SECTION:")
	authDNs := make([]string, 0)
	for i := 0; i < int(replyHdr.nscount); i++ {
		rr := &ResourceRecord{}
		n, err := rr.deserialize(buf, startidx, bytesRead)
		if err != nil {
			return err
		}
		startidx += n
		if rr.rrType == RR_NS {
			authDN, _, err := decodeDomainName(buf, startidx-len(rr.rdata), bytesRead)
			if err != nil {
				log.Println(err)
			}
			authDNs = append(authDNs, authDN)
			fmt.Printf("%v\t\t%v\tIN\tNS\t%v\n", rr.name, rr.ttl, authDN)
		}
	}
	fmt.Println()
	fmt.Println(";; ADDITIONAL SECTION:")
	servers := make([]*ResourceRecord, 0)
	for i := 0; i < int(replyHdr.arcount); i++ {
		rr := &ResourceRecord{}
		n, err := rr.deserialize(buf, startidx, bytesRead)
		if err != nil {
			return err
		}
		startidx += n
		servers = append(servers, rr)
		fmt.Printf("%v\t\t%v\tIN\tA\t%v\n", rr.name, rr.ttl, net.IP(rr.rdata))
	}
	// Check if we're done recursing
	for _, rr := range answers {
		if rr.name == dn {
			return nil
		}
	}
	// Query an authoritative server
	for _, rr := range servers {
		for _, authdn := range authDNs {
			if authdn == rr.name {
				return recursiveQuery(dn, net.IP(rr.rdata).String())
			}
		}
	}
	return nil
}

func (client *DNSClient) Query(dn string, recursive bool) error {
	if recursive {
		return recursiveQuery(dn, ROOT_IP) // Query a root server
	}
	return recursiveQuery(dn, "8.8.8.8")
}
