package pkg

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var ROOT_IP string = "199.9.14.201"

type DNSClient struct {
	cache *Cache
}

func NewDNSClient() *DNSClient {
	return &DNSClient{
		cache: NewCache(DB_PATH),
	}
}

func (client *DNSClient) recursiveQuery(dn string, ns string, useCaching bool) error {
	// See if we already have the result.
	if useCaching {
		toks := strings.Split(dn, ".")
		for i := 0; i < len(toks); i++ {
			dnfrag := strings.Join(toks[i:], ".")
			rr, cacheHit := client.cache.Get(dnfrag)
			if cacheHit && rr.timestamp+rr.ttl >= uint32(time.Now().Unix()) {
				fmt.Println("### Cached response ###")
				fmt.Printf("%s\t%s\n", dnfrag, rr.addr.String())
				if i == 0 {
					return nil
				} else {
					return client.recursiveQuery(dn, rr.addr.String(), false)
				}
			}
		}
	}

	// Grab a connection to the nameserver.
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

	// Construct the DNS query, send query.
	message := newQuery([]Question{question}, nil, nil, nil)
	payload, err := message.serialize()
	if err != nil {
		return err
	}
	if _, err = conn.Write(payload); err != nil {
		return err
	}

	// Get response from the nameserver.
	buf := make([]byte, 512)
	bytesRead, err := conn.Read(buf)
	if err != nil {
		return err
	}

	// Deserialize nameserver response.
	replyHdr := &Header{}
	startidx := 12 // End of header
	replyHdr.deserialize(buf[:startidx])

	// Print out question.
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

	// Print out answers.
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
		answers = append(answers, rr)
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

	// Print out authorites.
	fmt.Println()
	fmt.Println(";; AUTHORITY SECTION:")
	authDNs := make([]string, 0)
	authRecords := make([]*ResourceRecord, 0)
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
			authRecords = append(authRecords, rr)
			fmt.Printf("%v\t\t%v\tIN\tNS\t%v\n", rr.name, rr.ttl, authDN)
		}
	}

	// Print out additional data.
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

	// Cache answers.
	for _, answer := range answers {
		if answer.rrType != RR_CNAME {
			cr := NewCachedRecord(answer)
			client.cache.Store(answer.name, cr)
		}
	}

	// Cache authDNs.
	added := make(map[string]bool)
	for i, authRecord := range authRecords {
		if _, found := added[authRecord.name]; !found {
			added[authRecord.name] = true
			for _, server := range servers {
				if server.name == authDNs[i] {
					cr := NewCachedRecord(server)
					client.cache.Store(authRecord.name, cr)
					break
				}
			}
		}
	}

	// Check if we're done recursing.
	for _, rr := range answers {
		if rr.name == dn {
			return nil
		}
	}

	// Query an authoritative server.
	for _, rr := range servers {
		for _, authdn := range authDNs {
			if authdn == rr.name {
				return client.recursiveQuery(dn, net.IP(rr.rdata).String(), useCaching)
			}
		}
	}
	return nil
}

func (client *DNSClient) Query(dn string, recursive, useCaching bool) error {
	// Make sure no trailing dot.
	if len(dn) > 0 && dn[len(dn)-1] == '.' {
		dn = dn[:len(dn)-1]
	}
	// Query.
	if recursive {
		return client.recursiveQuery(dn, ROOT_IP, useCaching) // Query a root server
	}
	return client.recursiveQuery(dn, "8.8.8.8", useCaching)
}
