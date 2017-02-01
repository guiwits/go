package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// ps aux | grep -v grep | grep "ProductClient.jar"
func main() {
	out, err := exec.Command("pgrep", "-f", "ProductClient.jar").Output()
	if err != nil {
		log.Fatal(err)
		fmt.Printf("CRITICAL. No PID found for PDL")
		os.Exit(2) // fatal
	}
	fmt.Printf("OK. The PID for PDL is %s", out)
	os.Exit(0)
}
