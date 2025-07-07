package ui // User interaction with the command line interface.

import (
	"flag"
	"fmt"
	"strconv"
)

func GetIdentity() string {
	name := flag.String("n", "(?)anon", "Identity of connected user.")
	flag.Parse()

	return *name
}

func DisplayUsers(listOfUsers []string) {
	
	fmt.Println("-------LIST OF CONNECTED USERS-----------")
	for i, v := range listOfUsers {
		fmt.Println(strconv.Itoa(i) + ". " + v)
	}
}
