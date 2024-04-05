package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var ipAddr string
	var help bool
	// Define CLI flags in shrot and long form
	flag.StringVar(&ipAddr, "i", "", "IP adress of host (shorthand)")
	flag.StringVar(&ipAddr, "ip-address", "", "IP adress of host")
	flag.BoolVar(&help, "h", false, "Show usage instructions (shorthand)")
	flag.BoolVar(&help, "help", false, "Show usage instructions")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Pass -i to specify the IP address of the host:")
		fmt.Fprintln(os.Stderr, "  -i\t\tstring\n  --ip-address")
		fmt.Fprintln(os.Stderr, "  \tIP adress of host")
		fmt.Fprintln(os.Stderr, "\n  -h\n  --help")
		fmt.Fprintln(os.Stderr, "  \tShow usage instructions")
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
	}
	flag.Parse()

	out, _ := exec.Command("ping", ipAddr, "-c 5", "-i 3").Output()
	if strings.Contains(string(out), "Destination Host Unreachable") || strings.Contains(string(out), "No route to host") || strings.Contains(string(out), "Request timeout for icmp_seq") {
		fmt.Println("{'result':false}")
	} else {
		fmt.Println("{'result':true}")
	}
}
