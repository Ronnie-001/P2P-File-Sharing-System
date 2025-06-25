package ui // User interaction with the command line interface.

import (
	"flag"
)

func SetIdentity() string {
	name := flag.String("n", "(?)anon", "Identity of connected user.")
	flag.Parse()

	return *name
}

func GrabUsers() ([]string, error) {return []string{}, nil}
