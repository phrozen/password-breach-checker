package main

import (
	"crypto/sha1"
	_ "embed"
	"flag"
	"fmt"

	"github.com/phrozen/password-breach-checker/pkg/database"
	"github.com/phrozen/password-breach-checker/pkg/format"
)

func main() {
	input := flag.String("f", "", "Input filename of passwords ordered by hash (.bin)")
	password := flag.String("p", "", "Password to check for breaches")
	flag.Parse()
	// Create a new Database from the binary file
	db, err := database.New(*input)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Printf("DB size: %s - length: %d passwords\n", format.Bytes(uint64(db.Length()*24)), db.Length())
	// Hash the password with SHA1 and check it against the database
	hash := sha1.Sum([]byte(*password))
	fmt.Printf("Password: '%s' found in %d breaches\n", *password, db.Search(hash[:]))
}
