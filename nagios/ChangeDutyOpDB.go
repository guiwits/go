package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"os"
	"time"
)

//
// Global map to map Last Name with username.
//
var userToDutyOp = map[string]string{
	"Last" : "username",
	"Last" : "username",
	"Last" : "username",
	"Last" : "username",
	"Last" : "username",
}

//
// DB Select something like:
// select firstname from table where weekstart='2016-03-12';
//
func getDutyOpFromDB() (string, string) {
	var fname string
	var lname string
	dbquery := "select whatever from whatever ?"
	year, month, day := time.Now().Date()
	opDate := fmt.Sprintf("%d-%.2d-%d", year, month, (day - 1))
	db, err := sql.Open("mysql", "username:password@tcp(hostname:<port>)/dbname")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(dbquery, opDate)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&fname, &lname)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	rows.Close()
	return fname, lname
}

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

func restartNagios() {
  cmd := exec.Command ("/etc/init.d/nagios", "restart")
  err := cmd.Start()
  if err != nil {
    log.Fatal(err)
    err = cmd.Wait()
    log.Printf ("Command finished with error: %v", err)
}

//
// Main func()
//
func main() {
	fname, lname := getDutyOpFromDB()
	dutyOp, present := userToDutyOp[lname]
	year, month, day := time.Now().Date()
	cfile := fmt.Sprintf("contacts.cfg.%d.%.2d.%d", year, month, day)
	ferr := copyFile("contacts.cfg", cfile)

	if ferr != nil {
		log.Fatal(ferr)
	}

	if present == false {
		fmt.Printf("Unable to find user: %s %s. Exiting.\n", fname, lname)
		os.Exit(1)
	}

	lines, err := readContacts("contacts.cfg", dutyOp)

	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	for i, line := range lines {
		fmt.Println(i, line)
	}

	if err := writeContacts(lines, "contacts.cfg"); err != nil {
		log.Fatalf("writeLines %s", err)
	}

  restartNagios()

	fmt.Printf("Duty Operator changed to %s %s\n", fname, lname)
}
