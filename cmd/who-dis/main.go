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
	err := client.Query(dn, true, true)
	if err != nil {
		fmt.Println(err)
	}
}
