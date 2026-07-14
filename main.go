package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "version" {
		fmt.Println("todo-bench seed")
		return
	}
	fmt.Fprintln(os.Stderr, "usage: todo <add|list|done>")
	os.Exit(2)
}
