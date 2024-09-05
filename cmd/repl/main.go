package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/eunomie/ducker/dockerfile/repl"
)

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Whalecome %s!\n", u.Username)
	fmt.Println("This is the Docker REPL")
	fmt.Println("Feel free to type in commands")
	repl.Start(os.Stdin, os.Stdout)
}
