package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main(){
	var ipAddress string
	var userName string
	var help bool
	flag.StringVar(&ipAddress, "i", "", "IP Address of target to shutdown (shorthand)")
	flag.StringVar(&ipAddress, "ip-address", "", "IP Address of target to shutdown")
	flag.StringVar(&userName, "u", "", "User name who will shutdown target (shorthand)")
	flag.StringVar(&userName, "user-name", "", "User name who will shutdown target")
	flag.BoolVar(&help, "h", false, "Show usage instructions (shorthand)")
	flag.BoolVar(&help, "help", false, "Show usage instructions")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Pass -i to specify IP Address of target to shutdown (required)")
		fmt.Fprintln(os.Stderr, "  -i\t\tstring\n  --ip-address")
		fmt.Fprintln(os.Stderr, "Pass -u to specify User name who will shutdown target (required)")
		fmt.Fprintln(os.Stderr, "  -u\t\tstring\n  --user-name")
		fmt.Fprintln(os.Stderr, "Pass -h to show usage instructions")
		fmt.Fprintln(os.Stderr, "  -h\t\tstring\n  --help")
		fmt.Fprintln(os.Stderr, "----------------------------------------------------------------------------------------")
	}
	flag.Parse()

	// Print the Help docuemntation to the terminal if user passes help flag
	if help {
		flag.Usage()
		return
	}

	// check required flags
	if ipAddress == "" || userName == "" {
		log.Fatal("from sol(): required parameters cannot be blank")
		return
	}

	// Execute the command
	cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", userName, ipAddress), "sudo", "shutdown", "now")
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "closed by remote host"){
		log.Fatalf("from sol(): function exec.Command() failed: %v", string(output))
		return
	}
	log.Printf("from sol(): %v", output)
}