package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	passwordFlag := flag.String("password", "", "Plain text password to hash (omit to read from stdin)")
	cost := flag.Int("cost", bcrypt.DefaultCost, "bcrypt cost (10-16 recommended)")
	flag.Parse()

	password := strings.TrimSpace(*passwordFlag)
	if password == "" {
		fmt.Print("Enter password: ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			fmt.Fprintln(os.Stderr, "no password provided")
			os.Exit(1)
		}
		password = strings.TrimSpace(scanner.Text())
	}

	if password == "" {
		fmt.Fprintln(os.Stderr, "password cannot be empty")
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), *cost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to hash password: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Password hash: %s\n", string(hash))
}
