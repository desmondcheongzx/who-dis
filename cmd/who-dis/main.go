package main

import (
	"flag"
	"fmt"

	dns "github.com/desmondcheongzx/who-dis/pkg"
)

func main() {
	var nocacheFlag bool
	flag.BoolVar(&nocacheFlag, "nocache", false, "Turn off caching.")
	var recursiveFlag bool
	flag.BoolVar(&recursiveFlag, "trace", false, "Turn on recursive queries.")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("usage: ./who-dis <domain name>")
		return
	}
	dn := args[0]
	client := dns.NewDNSClient()
	fmt.Printf("nochache = %v; trace = %v\n", nocacheFlag, recursiveFlag)
	err := client.Query(dn, recursiveFlag, !nocacheFlag)
	if err != nil {
		fmt.Println(err)
	}
}
