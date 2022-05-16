package main

import (
	"fmt"
	"os"

	dns "github.com/desmondcheongzx/who-dis/pkg"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./who-dis <domain name>")
		return
	}
	dn := os.Args[1]
	client := dns.NewDNSClient()
	ip, err := client.Query(dn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ip.String())
}
