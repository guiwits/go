package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

//
// Create a copy of the contracts file with the time and date appended to it.
//
func copyFile(source, destination string) (err error) {
	in, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()
	out, err := os.Create(destination)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		log.Fatal(err)
	}
	err = out.Sync()
	return
}

//
// Read in the configuration file and find the admins contact
//
func readContacts(path string, dutyOp string) ([]string, error) {
	newDutyOp := fmt.Sprintf("        members                 nagiosadmin, xxxx, xxxx, xxxx, %s", dutyOp)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "        alias                   Nagios Administrators" {
			fmt.Println("Found Nagios Administrators.")
			lines = append(lines, scanner.Text())
			scanner.Scan()
			lines = append(lines, newDutyOp)
		} else {
			lines = append(lines, scanner.Text())
		}
	}
	file.Close()
	return lines, scanner.Err()
}

//
// Write out each line back to contacts.cfg with new dutyop.
//
func writeContacts(lines []string, path string) error {

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

//
// Function to restart the nagios daemon
//
func restartNagios() {
	cmd := exec.Command("/etc/init.d/nagios", "restart")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}

//
// Main func()
//
func main() {
	var user string
	nagfile := "/opt/nagios/current/etc/objects/contacts.cfg"

	if len(os.Args) > 1 {
		user = os.Args[1]
	} else {
		fmt.Println("%s requires the new duty op user as a command line argument.\n", os.Args[0])
		os.Exit(-1)
	}

	year, month, day := time.Now().Date()
	cfile := fmt.Sprintf("/opt/nagios/current/etc/objects/contacts.cfg.%d.%.2d.%d", year, month, day)
	ferr := copyFile(nagfile, cfile)

	if ferr != nil {
		log.Fatal(ferr)
	} else {
		fmt.Printf("Created a backup of the contact.cfg file to file: %s\n", cfile)
	}

	lines, err := readContacts(nagfile, user)

	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	for i, line := range lines {
		fmt.Println(i, line)
	}

	if err := writeContacts(lines, nagfile); err != nil {
		log.Fatalf("writeLines %s", err)
	}

	restartNagios()
	fmt.Printf("Duty Operator changed to %s\n", user)
}
