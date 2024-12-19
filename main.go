package main

import (
	"flag"
	"golang.org/x/term"
	"os"
)
import (
	"encoding/json"
	"fmt"
	"github.com/TylerBrock/colorjson"
	"strings"
)
import "github.com/mpage/onepassword"

/*
* non-parameters
*   master pwd (not to get into command history in plaintext)
* required parameters:
*   dbpath
*   filter by title
*   filter by email
*   filter (generic)
* */
func main() {
	dbpath := flag.String("dbpath", "./OnePassword.sqlite", "path to OnePassword.sqlite database")
	filter := flag.String("filter", "", "filter results by")
	flag.Parse()
	fmt.Println("Enter master password: ")
	password, err := term.ReadPassword(0)
	if err != nil {
		fmt.Println("error reading entered password")
		os.Exit(-1)
	}
	pwd := string(password)
	s, err := onepassword.NewVault(pwd, onepassword.VaultConfig{DBPath: *dbpath, Profile: "default"})
	if err != nil {
		fmt.Println("incorrect master password")
		os.Exit(-2)
	}
	rows, err := s.LookupItems(func(i *onepassword.Item) bool {
		if strings.Contains(strings.ToLower(i.Title), strings.ToLower(*filter)) {
			return true
		}
		return false
	})
	if err != nil {
		fmt.Println("did not find login entry")
		os.Exit(-3)
	}
	for i := range len(rows) {
		r := rows[i]
		// det := string(r.Details[:])
		fmt.Printf("Title: %v, url: %v, tags: %v, cat: %v", r.Title, r.Url, r.Tags, r.Category)
		var obj map[string]interface{}
		json.Unmarshal(r.Details, &obj)

		// Make a custom formatter with indent set
		f := colorjson.NewFormatter()
		f.Indent = 4

		// Marshall the Colorized JSON
		s, _ := f.Marshal(obj)
		fmt.Println(string(s))
	}
}
